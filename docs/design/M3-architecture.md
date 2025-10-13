# GoLisp アーキテクチャ設計書（M3）

## 目的

このドキュメントは、GoLisp M3（関数とクロージャ）の技術選定とアーキテクチャ設計を記録します。

## 1. M3の目標

**関数定義と呼び出しができる**

```lisp
> (lambda (x) (* x x))
#<FUNCTION>

> ((lambda (x) (* x x)) 5)
25

> (defun square (x) (* x x))
SQUARE

> (square 5)
25

> (defun make-adder (n)
    (lambda (x) (+ x n)))
MAKE-ADDER

> (funcall (make-adder 10) 5)
15
```

## 2. 必要な機能

### 2.1 特殊形式

M3で実装する特殊形式：

- `lambda` - 無名関数（クロージャ）
- `defun` - 関数定義
- `quote` / `'` - クォート（評価を抑制）
- `if` - 条件分岐

### 2.2 組み込み関数

- `funcall` - 関数呼び出し
- `apply` - 引数リストを展開して関数呼び出し

## 3. データ構造の拡張

### 3.1 Lambda（ユーザー定義関数）

```go
// internal/types/types.go に追加

// Lambda はユーザー定義関数を表現
type Lambda struct {
    Params []string      // 仮引数のリスト
    Body   Expr          // 関数本体（S式）
    Env    *Environment  // クロージャ（定義時の環境を保持）
}

func (l *Lambda) String() string {
    return "#<FUNCTION>"
}
```

**重要**: `Lambda`は定義時の環境（`Env`）を保持することで、**クロージャ**を実現します。

### 3.2 なぜEnvironmentへの参照が必要か

```lisp
(defun make-adder (n)
  (lambda (x) (+ x n)))  ; nを参照したい

(define add10 (make-adder 10))
(add10 5)  ; => 15
```

`lambda`が作られた時点で、`n=10`という環境を保持しておく必要があります。

## 4. 特殊形式の実装

### 4.1 特殊形式とは

通常の関数と異なり、**引数を評価せずに受け取る**形式です。

例：
```lisp
(if (> 5 3) 'yes 'no)
```

- 通常の関数: すべての引数を評価してから関数を呼び出す
- 特殊形式: 引数を評価せず、特殊形式の中で必要に応じて評価

### 4.2 実装方針

```go
// internal/eval/special.go
package eval

import (
    "fmt"
    "github.com/koplec/gospl/internal/types"
)

// evalSpecialForm は特殊形式を評価
func evalSpecialForm(name string, args types.Expr, env *Environment) (types.Expr, error) {
    switch name {
    case "quote":
        return evalQuote(args)
    case "if":
        return evalIf(args, env)
    case "lambda":
        return evalLambda(args, env)
    case "defun":
        return evalDefun(args, env)
    default:
        return nil, fmt.Errorf("unknown special form: %s", name)
    }
}

// isSpecialForm は特殊形式かどうかを判定
func isSpecialForm(name string) bool {
    switch name {
    case "quote", "if", "lambda", "defun":
        return true
    default:
        return false
    }
}
```

### 4.3 quote の実装

```go
// evalQuote は引数を評価せずにそのまま返す
func evalQuote(args types.Expr) (types.Expr, error) {
    // quoteは引数を1つだけ取る
    cons, ok := args.(*types.Cons)
    if !ok {
        return nil, fmt.Errorf("quote requires exactly 1 argument")
    }

    // 引数が1つだけか確認
    if _, ok := cons.Cdr.(*types.Nil); !ok {
        return nil, fmt.Errorf("quote requires exactly 1 argument")
    }

    // 評価せずにそのまま返す
    return cons.Car, nil
}
```

**動作例**:
```lisp
> 'x
X

> '(1 2 3)
(1 2 3)

> (quote (+ 1 2))
(+ 1 2)  ; 評価されない
```

### 4.4 if の実装

```go
// evalIf は条件分岐
func evalIf(args types.Expr, env *Environment) (types.Expr, error) {
    // (if condition then-expr else-expr)
    cons, ok := args.(*types.Cons)
    if !ok {
        return nil, fmt.Errorf("if requires 2 or 3 arguments")
    }

    // 条件を評価
    condition, err := Eval(cons.Car, env)
    if err != nil {
        return nil, err
    }

    // then-exprとelse-exprを取得
    rest, ok := cons.Cdr.(*types.Cons)
    if !ok {
        return nil, fmt.Errorf("if requires 2 or 3 arguments")
    }

    thenExpr := rest.Car

    // else-exprはオプション
    var elseExpr types.Expr = &types.Nil{}
    if cons2, ok := rest.Cdr.(*types.Cons); ok {
        elseExpr = cons2.Car
    }

    // 条件がnilまたはfalse以外なら真
    if isTrue(condition) {
        return Eval(thenExpr, env)
    } else {
        return Eval(elseExpr, env)
    }
}

// isTrue はLispの真偽値判定（nilとfalse以外はすべて真）
func isTrue(expr types.Expr) bool {
    if _, ok := expr.(*types.Nil); ok {
        return false
    }
    if b, ok := expr.(types.Boolean); ok {
        return b.Value
    }
    return true
}
```

**動作例**:
```lisp
> (if (> 5 3) 'yes 'no)
YES

> (if nil 'yes 'no)
NO

> (if t 42)
42

> (if nil 42)
NIL
```

### 4.5 lambda の実装

```go
// evalLambda は無名関数を作成
func evalLambda(args types.Expr, env *Environment) (types.Expr, error) {
    // (lambda (params...) body)
    cons, ok := args.(*types.Cons)
    if !ok {
        return nil, fmt.Errorf("lambda requires at least 2 arguments")
    }

    // 仮引数リストを解析
    params, err := parseParams(cons.Car)
    if err != nil {
        return nil, err
    }

    // 関数本体を取得（複数の式がある場合はprogn相当）
    rest, ok := cons.Cdr.(*types.Cons)
    if !ok {
        return nil, fmt.Errorf("lambda requires a body")
    }

    // M3では本体は1つの式のみサポート（prognは後で実装）
    body := rest.Car

    // クロージャを作成（現在の環境を保持）
    return &types.Lambda{
        Params: params,
        Body:   body,
        Env:    env,  // 重要: 定義時の環境を保持
    }, nil
}

// parseParams は仮引数リストをパース
func parseParams(expr types.Expr) ([]string, error) {
    // 空リストの場合
    if _, ok := expr.(*types.Nil); ok {
        return []string{}, nil
    }

    var params []string
    current := expr

    for {
        if _, ok := current.(*types.Nil); ok {
            break
        }

        cons, ok := current.(*types.Cons)
        if !ok {
            return nil, fmt.Errorf("invalid parameter list")
        }

        // パラメータはシンボルでなければならない
        sym, ok := cons.Car.(types.Symbol)
        if !ok {
            return nil, fmt.Errorf("parameter must be a symbol, got %T", cons.Car)
        }

        params = append(params, sym.Name)
        current = cons.Cdr
    }

    return params, nil
}
```

**動作例**:
```lisp
> (lambda (x) (* x x))
#<FUNCTION>

> ((lambda (x) (* x x)) 5)
25
```

### 4.6 defun の実装

```go
// evalDefun は関数を定義してグローバル環境に登録
func evalDefun(args types.Expr, env *Environment) (types.Expr, error) {
    // (defun name (params...) body)
    cons, ok := args.(*types.Cons)
    if !ok {
        return nil, fmt.Errorf("defun requires at least 3 arguments")
    }

    // 関数名を取得
    nameSym, ok := cons.Car.(types.Symbol)
    if !ok {
        return nil, fmt.Errorf("function name must be a symbol, got %T", cons.Car)
    }

    // 残りの引数は (params...) body と同じ
    rest := cons.Cdr

    // lambdaと同じ処理
    lambda, err := evalLambda(rest, env)
    if err != nil {
        return nil, err
    }

    // グローバル環境に登録
    env.Set(nameSym.Name, lambda)

    // 関数名のシンボルを返す
    return nameSym, nil
}
```

**動作例**:
```lisp
> (defun square (x) (* x x))
SQUARE

> (square 5)
25

> square
#<FUNCTION>
```

## 5. 評価器の更新

### 5.1 evalListの変更

```go
// internal/eval/eval.go

func evalList(list *types.Cons, env *Environment) (types.Expr, error) {
    if list == nil {
        return nil, fmt.Errorf("cannot evaluate nil pointer (internal error)")
    }

    // 先頭要素を取得（評価はまだしない）
    first := list.Car

    // シンボルの場合、特殊形式かチェック
    if sym, ok := first.(types.Symbol); ok {
        if isSpecialForm(sym.Name) {
            // 特殊形式: 引数を評価せずに渡す
            return evalSpecialForm(sym.Name, list.Cdr, env)
        }
    }

    // 通常の関数呼び出し: 先頭を評価
    fn, err := Eval(first, env)
    if err != nil {
        return nil, err
    }

    // 引数を評価
    args, err := evalArgs(list.Cdr, env)
    if err != nil {
        return nil, err
    }

    // 関数適用
    return apply(fn, args)
}
```

### 5.2 applyの拡張

```go
// internal/eval/eval.go

func apply(fn types.Expr, args []types.Expr) (types.Expr, error) {
    switch f := fn.(type) {
    case BuiltinFunc:
        // 組み込み関数
        return f.Call(args)

    case *types.Lambda:
        // ユーザー定義関数
        return applyLambda(f, args)

    default:
        return nil, fmt.Errorf("not a function: %v", fn)
    }
}

func applyLambda(lambda *types.Lambda, args []types.Expr) (types.Expr, error) {
    // 引数の数をチェック
    if len(args) != len(lambda.Params) {
        return nil, fmt.Errorf("wrong number of arguments: expected %d, got %d",
            len(lambda.Params), len(args))
    }

    // 新しい環境を作成（クロージャの環境を親とする）
    newEnv := NewEnvironment(lambda.Env)

    // 仮引数に実引数を束縛
    for i, param := range lambda.Params {
        newEnv.Set(param, args[i])
    }

    // 関数本体を新しい環境で評価
    return Eval(lambda.Body, newEnv)
}
```

## 6. funcallとapplyの実装

### 6.1 funcall

```go
// internal/eval/builtins.go に追加

func builtinFuncall(args []types.Expr) (types.Expr, error) {
    if len(args) < 1 {
        return nil, fmt.Errorf("funcall requires at least 1 argument")
    }

    fn := args[0]
    fnArgs := args[1:]

    return apply(fn, fnArgs)
}
```

**使用例**:
```lisp
> (funcall (lambda (x) (* x x)) 5)
25

> (defun add (x y) (+ x y))
ADD

> (funcall add 3 4)
7
```

### 6.2 apply（組み込み関数版）

```go
// internal/eval/builtins.go に追加

func builtinApply(args []types.Expr) (types.Expr, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("apply requires exactly 2 arguments")
    }

    fn := args[0]

    // 第2引数はリストでなければならない
    argList, err := listToSlice(args[1])
    if err != nil {
        return nil, err
    }

    return apply(fn, argList)
}

// listToSlice はLispリストをGoのスライスに変換
func listToSlice(expr types.Expr) ([]types.Expr, error) {
    if _, ok := expr.(*types.Nil); ok {
        return []types.Expr{}, nil
    }

    var result []types.Expr
    current := expr

    for {
        if _, ok := current.(*types.Nil); ok {
            break
        }

        cons, ok := current.(*types.Cons)
        if !ok {
            return nil, fmt.Errorf("not a proper list")
        }

        result = append(result, cons.Car)
        current = cons.Cdr
    }

    return result, nil
}
```

**使用例**:
```lisp
> (apply + '(1 2 3))
6

> (apply * '(2 3 4))
24
```

## 7. Environmentの循環参照対策

### 7.1 問題

`types.Lambda`が`*Environment`を持ち、`Environment`が`types.Expr`（Lambdaを含む）を持つため、循環参照になります。

```
Lambda → Environment → Lambda → ...
```

### 7.2 解決策

Goはガベージコレクションを持つため、循環参照自体は問題ありません。ただし、`Environment`を`types`パッケージから参照する必要があります。

**オプション1: Environmentをtypesパッケージに移動**

M3では`eval`パッケージのままで進め、後でリファクタリングすることも可能です。

**オプション2: インターフェースで抽象化**

```go
// internal/types/types.go

type Env interface {
    Get(name string) (Expr, error)
    Set(name string, value Expr)
}

type Lambda struct {
    Params []string
    Body   Expr
    Env    Env  // インターフェース
}
```

M3では**オプション1**（evalパッケージのまま）で進めます。

## 8. モジュール構成

### 8.1 ディレクトリ構造

```
gospl/
├── cmd/
│   └── gospl/
│       └── main.go
├── internal/
│   ├── types/
│   │   ├── types.go        ← Lambda追加
│   │   └── printer.go
│   ├── reader/
│   │   ├── lexer.go
│   │   ├── lexer_test.go
│   │   ├── parser.go
│   │   └── parser_test.go
│   ├── eval/
│   │   ├── eval.go         ← evalList, apply更新
│   │   ├── eval_test.go
│   │   ├── env.go
│   │   ├── env_test.go
│   │   ├── builtins.go     ← funcall, apply追加
│   │   ├── builtins_test.go
│   │   ├── special.go      ← 新規: 特殊形式
│   │   └── special_test.go ← 新規: テスト
│   └── repl/
│       └── repl.go
└── docs/
    ├── SPECIFICATION.md
    └── design/
        ├── M1-architecture.md
        ├── M2-architecture.md
        └── M3-architecture.md
```

## 9. M3完成の定義

### 9.1 実装する機能

- [ ] `internal/types/types.go` - Lambda構造体追加
- [ ] `internal/eval/special.go` - 特殊形式の実装
  - [ ] quote
  - [ ] if
  - [ ] lambda
  - [ ] defun
- [ ] `internal/eval/special_test.go` - テスト
- [ ] `internal/eval/eval.go` - evalList, apply更新
- [ ] `internal/eval/builtins.go` - funcall, apply追加
- [ ] `internal/eval/eval_test.go` - 統合テスト追加

### 9.2 テスト例

```lisp
> 'x
X

> '(1 2 3)
(1 2 3)

> (if (> 5 3) 'yes 'no)
YES

> (if nil 'yes 'no)
NO

> (lambda (x) (* x x))
#<FUNCTION>

> ((lambda (x) (* x x)) 5)
25

> (defun square (x) (* x x))
SQUARE

> (square 5)
25

> (defun make-adder (n)
    (lambda (x) (+ x n)))
MAKE-ADDER

> (funcall (make-adder 10) 5)
15

> (apply + '(1 2 3))
6
```

## 10. 実装の順序

### ステップ1: Lambda型の追加
1. `internal/types/types.go`に`Lambda`構造体を追加
2. `String()`メソッドの実装

### ステップ2: quote の実装
1. `internal/eval/special.go`を作成
2. `evalQuote`の実装
3. `isSpecialForm`の実装
4. `evalList`の更新（特殊形式のチェック）
5. テスト作成

### ステップ3: if の実装
1. `evalIf`の実装
2. `isTrue`ヘルパー関数
3. テスト作成

### ステップ4: lambda の実装
1. `evalLambda`の実装
2. `parseParams`ヘルパー関数
3. `apply`の拡張（Lambdaサポート）
4. `applyLambda`の実装
5. テスト作成

### ステップ5: defun の実装
1. `evalDefun`の実装
2. テスト作成

### ステップ6: funcall と apply
1. `builtinFuncall`の実装
2. `builtinApply`の実装
3. `listToSlice`ヘルパー関数
4. グローバル環境への登録
5. テスト作成

### ステップ7: 統合テスト
1. クロージャのテスト
2. ネストした関数呼び出し
3. エラーケース

## 11. 注意点

### 11.1 本体が複数式の場合

M3では本体は1つの式のみサポートします。

```lisp
; M3では動かない
(lambda (x)
  (print x)
  (* x x))

; M4以降でprognを実装後に対応
```

### 11.2 可変長引数

`&rest`などの可変長引数はM13で実装予定です。

### 11.3 レキシカルスコープ

`Lambda`が`Env`を保持することで、レキシカルスコープ（クロージャ）を実現します。

## 12. クロージャの動作確認

```lisp
> (defun make-counter ()
    (lambda () 0))  ; M3ではカウンタは未実装（setqが必要）
MAKE-COUNTER

; クロージャの基本動作確認
> (defun make-adder (n)
    (lambda (x) (+ x n)))
MAKE-ADDER

> (define add10 (make-adder 10))
ADD10

> (funcall add10 5)
15

> (funcall (make-adder 20) 5)
25
```

---

**作成日**: 2025-10-13
**対象マイルストーン**: M3（関数とクロージャ）
