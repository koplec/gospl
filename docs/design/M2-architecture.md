# GoLisp アーキテクチャ設計書（M2）

## 目的

このドキュメントは、GoLisp M2（基本的な評価器と算術）の技術選定とアーキテクチャ設計を記録します。

## 1. M2の目標

**数値と基本的な算術演算が動く**

```lisp
> (+ 1 2)
3

> (* 3 4)
12

> (+ 1 2 3 4 5)
15
```

## 2. 評価器（Evaluator）の設計

### 2.1 基本概念

評価器は、S式（types.Expr）を受け取り、評価した結果を返します。

```
入力: (+ 1 2)
  ↓
S式: Cons{Car: Symbol{"+"}, Cdr: Cons{Car: Number{1}, Cdr: Cons{Car: Number{2}, Cdr: Nil}}}
  ↓
評価: +関数を1と2に適用
  ↓
結果: Number{3}
```

### 2.2 実装方式

```go
// internal/eval/eval.go
package eval

import (
    "fmt"
    "github.com/koplec/gospl/internal/types"
)

// Eval はS式を評価して結果を返す
func Eval(expr types.Expr, env *Environment) (types.Expr, error) {
    switch e := expr.(type) {
    case types.Number:
        // 数値は自己評価（そのまま返す）
        return e, nil

    case types.String:
        // 文字列は自己評価
        return e, nil

    case types.Boolean:
        // 真偽値は自己評価
        return e, nil

    case *types.Nil:
        // nilは自己評価
        return e, nil

    case types.Symbol:
        // シンボルは変数参照（環境から値を取得）
        return env.Get(e.Name)

    case *types.Cons:
        // リストは関数適用
        return evalList(e, env)

    default:
        return nil, fmt.Errorf("unknown expression type: %T", expr)
    }
}

// evalList はリスト（関数適用）を評価
func evalList(list *types.Cons, env *Environment) (types.Expr, error) {
    // 空リストはエラー
    if list == nil {
        return nil, fmt.Errorf("cannot evaluate empty list")
    }

    // 先頭要素（関数）を評価
    fn, err := Eval(list.Car, env)
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

// evalArgs は引数リストを評価
func evalArgs(expr types.Expr, env *Environment) ([]types.Expr, error) {
    var args []types.Expr

    current := expr
    for {
        // nilなら終了
        if _, ok := current.(*types.Nil); ok {
            break
        }

        // Consでなければエラー
        cons, ok := current.(*types.Cons)
        if !ok {
            return nil, fmt.Errorf("invalid argument list")
        }

        // 引数を評価
        arg, err := Eval(cons.Car, env)
        if err != nil {
            return nil, err
        }
        args = append(args, arg)

        current = cons.Cdr
    }

    return args, nil
}

// apply は関数を引数に適用
func apply(fn types.Expr, args []types.Expr) (types.Expr, error) {
    // M2では組み込み関数のみサポート
    builtin, ok := fn.(BuiltinFunc)
    if !ok {
        return nil, fmt.Errorf("not a function: %v", fn)
    }

    return builtin.Call(args)
}
```

### 2.3 評価の流れ

**例: `(+ 1 2)`**

```
1. Eval((+ 1 2), env)
   ↓
2. evalList(Cons{+, (1 2)}, env)
   ↓
3. Eval(+, env) → BuiltinFunc(+)
   ↓
4. evalArgs((1 2), env) → [Number{1}, Number{2}]
   ↓
5. apply(BuiltinFunc(+), [Number{1}, Number{2}])
   ↓
6. Number{3}
```

## 3. 環境（Environment）の設計

### 3.1 役割

環境は、変数名とその値の対応を管理します。

- 変数の束縛を保存
- 変数の値を取得
- スコープの管理（親環境への参照）

### 3.2 実装方式

```go
// internal/eval/env.go
package eval

import (
    "fmt"
    "github.com/koplec/gospl/internal/types"
)

// Environment は変数の束縛を管理
type Environment struct {
    bindings map[string]types.Expr
    parent   *Environment  // 親環境（スコープチェーン用、M3以降で使用）
}

// NewEnvironment は新しい環境を作成
func NewEnvironment(parent *Environment) *Environment {
    return &Environment{
        bindings: make(map[string]types.Expr),
        parent:   parent,
    }
}

// Set は変数に値を束縛
func (e *Environment) Set(name string, value types.Expr) {
    e.bindings[name] = value
}

// Get は変数の値を取得
func (e *Environment) Get(name string) (types.Expr, error) {
    // 現在の環境で探す
    if value, ok := e.bindings[name]; ok {
        return value, nil
    }

    // 親環境で探す
    if e.parent != nil {
        return e.parent.Get(name)
    }

    // 見つからない
    return nil, fmt.Errorf("undefined variable: %s", name)
}
```

### 3.3 グローバル環境の初期化

```go
// internal/eval/env.go

// NewGlobalEnvironment はグローバル環境を作成
func NewGlobalEnvironment() *Environment {
    env := NewEnvironment(nil)

    // 組み込み関数を登録
    env.Set("+", BuiltinFunc{Name: "+", Fn: builtinAdd})
    env.Set("-", BuiltinFunc{Name: "-", Fn: builtinSub})
    env.Set("*", BuiltinFunc{Name: "*", Fn: builtinMul})
    env.Set("/", BuiltinFunc{Name: "/", Fn: builtinDiv})

    return env
}
```

## 4. 組み込み関数の設計

### 4.1 組み込み関数の型

```go
// internal/eval/builtins.go
package eval

import (
    "fmt"
    "github.com/koplec/gospl/internal/types"
)

// BuiltinFunc は組み込み関数を表現
type BuiltinFunc struct {
    Name string
    Fn   func([]types.Expr) (types.Expr, error)
}

// Call は組み込み関数を呼び出し
func (b BuiltinFunc) Call(args []types.Expr) (types.Expr, error) {
    return b.Fn(args)
}

// String はBuiltinFuncの文字列表現
func (b BuiltinFunc) String() string {
    return fmt.Sprintf("#<BUILTIN %s>", b.Name)
}
```

### 4.2 算術関数の実装

```go
// internal/eval/builtins.go

// builtinAdd は加算 (+)
func builtinAdd(args []types.Expr) (types.Expr, error) {
    if len(args) == 0 {
        return types.Number{Value: 0}, nil
    }

    var sum float64
    for _, arg := range args {
        num, ok := arg.(types.Number)
        if !ok {
            return nil, fmt.Errorf("+ expects numbers, got %T", arg)
        }
        sum += num.Value
    }

    return types.Number{Value: sum}, nil
}

// builtinSub は減算 (-)
func builtinSub(args []types.Expr) (types.Expr, error) {
    if len(args) == 0 {
        return nil, fmt.Errorf("- requires at least 1 argument")
    }

    first, ok := args[0].(types.Number)
    if !ok {
        return nil, fmt.Errorf("- expects numbers, got %T", args[0])
    }

    if len(args) == 1 {
        // 単項マイナス: (- 5) → -5
        return types.Number{Value: -first.Value}, nil
    }

    result := first.Value
    for _, arg := range args[1:] {
        num, ok := arg.(types.Number)
        if !ok {
            return nil, fmt.Errorf("- expects numbers, got %T", arg)
        }
        result -= num.Value
    }

    return types.Number{Value: result}, nil
}

// builtinMul は乗算 (*)
func builtinMul(args []types.Expr) (types.Expr, error) {
    if len(args) == 0 {
        return types.Number{Value: 1}, nil
    }

    result := 1.0
    for _, arg := range args {
        num, ok := arg.(types.Number)
        if !ok {
            return nil, fmt.Errorf("* expects numbers, got %T", arg)
        }
        result *= num.Value
    }

    return types.Number{Value: result}, nil
}

// builtinDiv は除算 (/)
func builtinDiv(args []types.Expr) (types.Expr, error) {
    if len(args) == 0 {
        return nil, fmt.Errorf("/ requires at least 1 argument")
    }

    first, ok := args[0].(types.Number)
    if !ok {
        return nil, fmt.Errorf("/ expects numbers, got %T", args[0])
    }

    if len(args) == 1 {
        // 逆数: (/ 2) → 0.5
        if first.Value == 0 {
            return nil, fmt.Errorf("division by zero")
        }
        return types.Number{Value: 1.0 / first.Value}, nil
    }

    result := first.Value
    for _, arg := range args[1:] {
        num, ok := arg.(types.Number)
        if !ok {
            return nil, fmt.Errorf("/ expects numbers, got %T", arg)
        }
        if num.Value == 0 {
            return nil, fmt.Errorf("division by zero")
        }
        result /= num.Value
    }

    return types.Number{Value: result}, nil
}
```

## 5. types.Exprへの追加

BuiltinFuncを`types.Expr`として扱えるようにする必要があります。

**オプション1: evalパッケージ内で完結（推奨）**

`BuiltinFunc`は`eval`パッケージ内で定義し、`types.Expr`インターフェースを満たすようにする。

```go
// internal/eval/builtins.go

// BuiltinFunc は types.Expr を実装
func (b BuiltinFunc) String() string {
    return fmt.Sprintf("#<BUILTIN %s>", b.Name)
}
```

**オプション2: typesパッケージに追加**

後で型システムを実装する際に、型情報が必要になる可能性があるため、`types`パッケージに移動することも検討。

M2では**オプション1**で進め、必要に応じて後でリファクタリングします。

## 6. REPLの更新

```go
// internal/repl/repl.go
package repl

import (
    "bufio"
    "fmt"
    "os"

    "github.com/koplec/gospl/internal/eval"
    "github.com/koplec/gospl/internal/reader"
)

func Start() {
    scanner := bufio.NewScanner(os.Stdin)
    env := eval.NewGlobalEnvironment()  // グローバル環境を作成

    fmt.Println("Gospl REPL")

    for {
        fmt.Print("> ")

        if !scanner.Scan() {
            break
        }

        input := scanner.Text()

        // Read
        parser := reader.NewParser(input)
        expr, err := parser.Parse()
        if err != nil {
            fmt.Printf("Parse error: %v\n", err)
            continue
        }

        // Eval (M2で追加)
        result, err := eval.Eval(expr, env)
        if err != nil {
            fmt.Printf("Eval error: %v\n", err)
            continue
        }

        // Print
        fmt.Println(result.String())
    }
}
```

## 7. モジュール構成

### 7.1 ディレクトリ構造

```
gospl/
├── cmd/
│   └── gospl/
│       └── main.go
├── internal/
│   ├── types/
│   │   ├── types.go
│   │   └── printer.go
│   ├── reader/
│   │   ├── lexer.go
│   │   ├── lexer_test.go
│   │   ├── parser.go
│   │   └── parser_test.go
│   ├── eval/              ← 新規追加
│   │   ├── eval.go
│   │   ├── eval_test.go
│   │   ├── env.go
│   │   ├── env_test.go
│   │   ├── builtins.go
│   │   └── builtins_test.go
│   └── repl/
│       └── repl.go
└── docs/
    ├── SPECIFICATION.md
    └── design/
        ├── M1-architecture.md
        └── M2-architecture.md
```

### 7.2 依存関係

```
cmd/gospl/main.go
  ↓
internal/repl/repl.go
  ↓
internal/eval/eval.go → internal/eval/env.go
  ↓                      ↓
internal/eval/builtins.go
  ↓
internal/types/types.go
  ↑
internal/reader/parser.go → internal/reader/lexer.go
```

## 8. M2完成の定義

### 8.1 実装する機能

- [X] `internal/eval/eval.go` - 評価器
- [X] `internal/eval/env.go` - 環境
- [X] `internal/eval/builtins.go` - 組み込み関数
- [X] テストファイル作成
- [X] REPLの更新

### 8.2 テスト例

```lisp
> 42
42

> (+ 1 2)
3

> (+ 1 2 3 4 5)
15

> (* 3 4)
12

> (- 10 3)
7

> (/ 10 2)
5

> (+ (* 2 3) (/ 10 5))
8
```

### 8.3 エラー処理例

```lisp
> (+ 1 "hello")
Eval error: + expects numbers, got types.String

> (/ 5 0)
Eval error: division by zero

> undefined-var
Eval error: undefined variable: undefined-var
```

## 9. 実装の順序

### ステップ1: 環境の実装
1. `internal/eval/env.go` を作成
2. `Environment` 構造体の実装
3. `Set` / `Get` メソッドの実装
4. テスト作成

### ステップ2: 組み込み関数の実装
1. `internal/eval/builtins.go` を作成
2. `BuiltinFunc` 型の定義
3. 算術関数の実装（`+`, `-`, `*`, `/`）
4. `NewGlobalEnvironment` の実装
5. テスト作成

### ステップ3: 評価器の実装
1. `internal/eval/eval.go` を作成
2. `Eval` 関数の実装
3. `evalList` の実装
4. `evalArgs` の実装
5. `apply` の実装
6. テスト作成

### ステップ4: REPLの更新
1. `internal/repl/repl.go` を更新
2. グローバル環境の初期化
3. `Eval` の呼び出しを追加
4. 動作確認

## 10. 注意点

### 10.1 Common Lispとの違い

- **数値**: M2では整数と浮動小数点を区別せず、すべて`float64`で扱う
- **可変長引数**: `(+)` → `0`, `(*)` → `1` のように、引数なしも許可
- **単項演算**: `(- 5)` → `-5`, `(/ 2)` → `0.5`

### 10.2 今後の拡張

M3以降で実装予定:
- ユーザー定義関数（lambda, defun）
- 特殊形式（quote, if など）
- クロージャ

---

**作成日**: 2025-10-12
**対象マイルストーン**: M2（基本的な評価器と算術）
