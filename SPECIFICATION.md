# GoLisp 仕様書

## プロジェクトの目的

Go言語でLispインタープリタを実装することで、Lispの思想を理解する。
Common Lisp寄りの実装に、TypeScript風の型システムを組み合わせ、さらにGo言語との相互運用を可能にする。

## 設計方針

- **Common Lispのサブセット**: 基本的なCommon Lispコードが動作することを目指す（完全互換ではない）
- **段階的型付け（Gradual Typing）**: 型宣言はオプション。型なしでも動作し、型で制約することも可能
- **Go連携**: Lisp内部からGo関数を呼び出せる
- **Lisp-2**: Common Lispと同様に、関数と変数の名前空間を分離

## 最終ゴール

### コア機能

#### データ型
- **プリミティブ型**
  - 整数（integer）
  - 浮動小数点数（float）
  - 文字列（string）
  - シンボル（symbol）
  - 真偽値（boolean）: `t`, `nil`

- **複合型**
  - コンスセル（cons）
  - リスト（list）
  - ベクター/配列（vector/array）
  - ハッシュテーブル（hash-table）

#### 特殊形式
- `quote` / `'` - クォート
- `if` - 条件分岐
- `lambda` - 無名関数
- `defun` - 関数定義
- `defvar` / `defparameter` - グローバル変数定義
- `setq` / `setf` - 代入
- `let` / `let*` - ローカル変数束縛
- `progn` - 順次実行
- `cond` - 多分岐条件
- `when` / `unless` - 単純条件
- `dolist` / `dotimes` / `do` - ループ
- `defmacro` - マクロ定義
- `function` / `#'` - 関数オブジェクト参照

#### 組み込み関数

**算術演算**
- `+`, `-`, `*`, `/`
- `=`, `<`, `>`, `<=`, `>=`
- `min`, `max`, `abs`

**リスト操作**
- `cons`, `car`, `cdr`
- `list`, `append`, `reverse`
- `first`, `rest`, `nth`, `length`
- `mapcar`, `remove-if`

**述語**
- `null`, `atom`, `listp`, `consp`
- `numberp`, `symbolp`, `stringp`
- `eq`, `eql`, `equal`

**関数操作**
- `funcall`, `apply`

**配列/ベクター**
- `make-array`, `aref`, `vector`

**ハッシュテーブル**
- `make-hash-table`, `gethash`, `sethash`

**文字列**
- `make-string`, `string=`, `concatenate`

**その他**
- `print`, `format`
- `macroexpand`

### 型システム

#### 基本的な型

```lisp
;; プリミティブ型
integer
float
string
boolean
symbol

;; 複合型
cons
list
vector
hash-table

;; 関数型
(function (arg-type1 arg-type2 ...) return-type)

;; ユニオン型
(or type1 type2 ...)

;; 任意の型
t
```

#### 型宣言の構文

```lisp
;; 変数の型宣言
(defvar x integer)
(defvar y (or integer string))

;; 関数の型宣言（引数と戻り値）
(defun add ((x integer) (y integer)) integer
  (+ x y))

;; declare文を使った型宣言（Common Lisp互換）
(defun add (x y)
  (declare (type integer x y))
  (+ x y))

;; 型なしも許容（動的型付け）
(defun add (x y)
  (+ x y))
```

#### 型チェック

- **静的型チェック**: 型宣言がある場合、評価前（実行前）に型をチェック
- **型推論**: 可能な範囲で型を推論（リテラル、関数戻り値など）
- **段階的導入**: 一部の関数だけ型付けすることも可能
- **型なしコードとの共存**: 型アノテーションがないコードは動的型付けとして動作

### Go連携

```lisp
;; Go関数をLispから呼び出す
(go:call "fmt.Println" "Hello from Lisp")

;; Go関数をLisp関数として登録
(defun println (x)
  (go:call "fmt.Println" x))
```

Go側からは：
- Go関数をLisp環境に登録
- Lisp関数をGoから呼び出し
- Lisp値とGo値の相互変換

### パッケージシステム（簡易版）

```lisp
(defpackage :myapp
  (:use :common-lisp))

(in-package :myapp)

;; パッケージ修飾子
common-lisp:car
myapp::internal-function
```

### コメント

```lisp
;; 行コメント

#|
  ブロックコメント
  複数行対応
|#
```

## マイルストーン

### M1: 基本的なS式とREPL

**目標**: S式を読み込んで表示できるREPLを作る

**実装内容**:
- データ構造の定義
  - Cons, Symbol, Number, String, Boolean
  - S式を表すインターフェース
- トークナイザ（Lexer）
  - 文字列をトークン列に分割
  - `(`, `)`, 数値、シンボル、文字列リテラル
- パーサ（Reader）
  - トークン列からS式を構築
  - ネストしたリストの解析
- プリンタ
  - S式を文字列表現に変換
- REPL
  - Read-Eval-Print-Loop の骨格
  - この段階ではEvalは恒等関数（そのまま返す）

**動作例**:
```
> (+ 1 2)
(+ 1 2)
> (list 1 2 3)
(list 1 2 3)
```

### M2: 基本的な評価器と算術

**目標**: 数値と基本的な算術演算が動く

**実装内容**:
- 評価器（Evaluator）の基本実装
- 自己評価オブジェクト（数値、文字列）
- 環境（Environment）の実装
  - 変数/シンボルの束縛管理
  - スコープの実装
- 基本的な算術関数
  - `+`, `-`, `*`, `/`
  - 可変長引数の処理
- 組み込み関数の登録メカニズム

**動作例**:
```lisp
> (+ 1 2)
3
> (* 3 4)
12
> (+ 1 2 3 4 5)
15
```

### M3: 関数とクロージャ

**目標**: 関数定義と呼び出しができる

**実装内容**:
- `lambda`式の実装
  - 引数リストと本体の解析
  - クロージャ（環境を捕捉）
- 関数呼び出しの実装
  - 引数の評価
  - 関数適用
- `defun`の実装
  - グローバル関数定義
  - 関数名前空間への登録
- `funcall`, `apply`の実装

**動作例**:
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

### M4a: 基本的な型システム（実行時型チェック）

**目標**: 型宣言と実行時型チェックができる（まず動くものを作る）

**実装内容**:
- 型の内部表現
  - Type インターフェース
  - プリミティブ型（Integer, Float, String, Boolean, Symbol）
  - 複合型（List, Cons）
- 型アノテーションのパース
  - 関数の引数型・戻り値型の解析
  - 変数の型宣言の解析
  - `(defun add ((x integer) (y integer)) integer ...)` 形式のサポート
- **実行時型チェック**の実装
  - 関数呼び出し時に引数の型をチェック
  - 関数の戻り値の型をチェック
  - 変数への代入時に型をチェック
- 型エラーの報告（実行時）
- 型述語の実装
  - `numberp`, `symbolp`, `listp`, `stringp`, `atom`, `null`
- `declare`フォームの実装（Common Lisp互換）

**処理の流れ**:
```
入力: (add 1 "hello")
  ↓
1. Parser: S式に変換
  ↓
2. Evaluator: 評価開始
  ↓
3. 関数呼び出し
  - 引数を評価: 1, "hello"
  - 型チェック: integer, string ← ここでエラー検出
  - エラー: 引数2はintegerを期待
```

**動作例**:
```lisp
> (defun add ((x integer) (y integer)) integer
    (+ x y))
ADD

> (add 1 2)
3

> (add 1 "hello")
Error: Type error (at runtime)
  Function: add
  Argument 2: expected integer, got string
  At: (add 1 "hello")

> (defun square ((x integer)) integer
    (* x x))
SQUARE

> (square 5)
25

> (defvar result integer)
RESULT

> (setq result 42)
42

> (setq result "invalid")
Error: Type error (at runtime)
  Variable: result
  Expected: integer
  Got: string

> (numberp 42)
T

> (symbolp 'foo)
T

;; 型なしコードも動作（動的型付け）
> (defun untyped (x) (* x x))
UNTYPED

> (untyped 5)
25

> (untyped "hi")
Error: Runtime error: * expects numbers
```

**実装のポイント**:
- 型情報は関数定義時に保存
- 実際のチェックは実行時（関数呼び出し時）
- 実装が簡単で、すぐに動作確認できる

### M4b: 静的型チェック（実行前エラー検出）

**目標**: 評価前に型エラーを検出する（PHP/Pythonより安全に）

**実装内容**:
- **静的型チェッカー**の実装
  - 評価前に型チェックパスを実行
  - S式全体を走査して型を検証
- リテラルの型推論
  - `1` → integer
  - `"hello"` → string
  - `'symbol` → symbol
- 変数の型追跡
  - 環境に型情報を保持
  - 変数参照時に型を返す
- 関数呼び出しの型検証
  - 引数の型を事前に検証
  - 戻り値の型を返す
- より詳細なエラーメッセージ

**処理の流れ**:
```
入力: (add 1 "hello")
  ↓
1. Parser: S式に変換
  ↓
2. Type Checker: 型チェックパス ← 新規追加
  - addの型シグネチャ: (integer, integer) -> integer
  - 引数1の型: integer ✓
  - 引数2の型: string ✗
  - エラー: "Expected integer, got string"
  ↓
3. Evaluator: 評価（型チェックが通った場合のみ実行）
```

**動作例**:
```lisp
> (defun add ((x integer) (y integer)) integer
    (+ x y))
ADD

> (add 1 2)
3

> (add 1 "hello")
Error: Type error (before execution)
  Function: add
  Argument 2: expected integer, got string
  At: (add 1 "hello")

> (defun square ((x integer)) integer
    (* x x))
SQUARE

> (defun use-square ()
    (square "oops"))
Error: Type error (before execution)
  Function: square
  Expected argument: integer
  Got: string
  At: (square "oops")
  In function: use-square

;; 関数を定義しただけではエラーにならない
> (defun bad-func ()
    (add 1 "hello"))
BAD-FUNC

;; 呼び出そうとするとエラー
> (bad-func)
Error: Type error (before execution)
  Function: add
  Expected argument 2: integer
  Got: string
  At: (add 1 "hello")
  In function: bad-func
```

**実装のポイント**:
- M4aの実装を活用（型情報の管理は同じ）
- 評価器の前に型チェッカーを追加
- REPLでは入力ごとに型チェック → 評価の順で実行

### M5: リスト操作と制御構造

**目標**: 基本的なリスト操作と条件分岐ができる

**実装内容**:
- `quote` / `'`の実装
- リスト操作関数
  - `car`, `cdr`, `cons`, `list`
  - `first`, `rest`, `nth`, `length`
  - `append`, `reverse`
- 条件分岐
  - `if`
  - `cond`
  - `when`, `unless`
- 真偽値の扱い
  - `nil`はfalse、それ以外はtrue
- 比較関数
  - `eq`, `eql`, `equal`
  - `=`, `<`, `>`, `<=`, `>=`
- `progn`
- `setq`, `defvar`

**動作例**:
```lisp
> '(1 2 3)
(1 2 3)

> (car '(1 2 3))
1

> (cdr '(1 2 3))
(2 3)

> (cons 0 '(1 2 3))
(0 1 2 3)

> (if (> 5 3) 'yes 'no)
YES

> (cond
    ((< 5 3) 'less)
    ((> 5 3) 'greater)
    (t 'equal))
GREATER

> (defvar *x* 10)
*X*

> *x*
10

> (setq *x* 20)
20
```

### M6: より多くのCommon Lisp関数

**目標**: よく使われるCommon Lisp関数を実装

**実装内容**:
- 高階関数
  - `mapcar`
  - `remove-if`, `remove-if-not`
  - `find-if`
- 数値関数
  - `min`, `max`, `abs`
  - `floor`, `ceiling`, `round`
  - `mod`, `rem`
- リスト関数
  - `member`, `assoc`
  - `push`, `pop`（マクロとして）
- 述語
  - `consp`, `listp`
  - `evenp`, `oddp`
  - `zerop`, `plusp`, `minusp`
- 文字列関数（基本）
  - `string=`, `string<`
  - `concatenate`
- 入出力
  - `print`, `princ`, `prin1`
  - `format`（簡易版）

**動作例**:
```lisp
> (mapcar (lambda (x) (* x x)) '(1 2 3 4 5))
(1 4 9 16 25)

> (remove-if (lambda (x) (< x 5)) '(1 3 5 7 9))
(5 7 9)

> (max 3 1 4 1 5 9)
9

> (format t "Hello, ~a!~%" "World")
Hello, World!
NIL
```

### M7: 高度な型推論とユニオン型

**目標**: より高度な型機能を実装（TypeScript風）

**実装内容**:
- ユニオン型 `(or type1 type2 ...)`
  - 静的チェック対応
  - 複数の型を許容
- 関数型 `(function (arg-types...) return-type)`
- リスト型 `(list element-type)`
- **高度な型推論**
  - 関数戻り値の型推論
  - 変数への代入からの型推論
  - 式全体の型推論
- 型エラーメッセージの改善
- （オプション）Type Narrowing
  - 条件分岐による型の絞り込み

**動作例**:
```lisp
;; ユニオン型
> (defvar x (or integer string))
X

> (setq x 42)
42

> (setq x "hello")
"hello"

> (setq x 3.14)
Error: Type error (before execution)
  Variable: x
  Expected: (or integer string)
  Got: float

;; 関数戻り値の型推論
> (defun square ((x integer)) integer
    (* x x))
SQUARE

> (defun use-square ((x integer))
    (square x))  ; squareの戻り値がintegerと推論される
USE-SQUARE

> (defun bad-call ()
    (square "oops"))
Error: Type error (before execution)
  Function: square
  Expected argument: integer
  Got: string

;; 変数への代入からの型推論（型アノテーションなし）
> (defvar y 42)        ; yの型はintegerと推論
Y

> (setq y "hello")     ; Error: yはinteger型
Error: Type error
  Variable: y (inferred type: integer)
  Got: string

;; ユニオン型の例
> (defun maybe-inc ((x (or integer nil))) (or integer nil)
    (if x (+ x 1) nil))
MAYBE-INC

> (maybe-inc 5)
6

> (maybe-inc nil)
NIL

> (maybe-inc "bad")
Error: Type error (before execution)
  Expected: (or integer nil)
  Got: string

;; Type Narrowing（オプション機能）
> (defun process ((x (or integer string)))
    (if (numberp x)
        (+ x 1)        ; xはここではinteger
        (length x)))   ; xはここではstring
PROCESS

> (process 5)
6

> (process "hello")
5
```

**実装の優先順位**:
1. ユニオン型の基本サポート（必須）
2. 関数戻り値の型推論（必須）
3. 変数の型推論（推奨）
4. Type Narrowing（オプション、実装が複雑）

### M8: Go連携機能

**目標**: LispとGoの相互運用を実現

**実装内容**:
- Lisp値とGo値の相互変換
  - Number ↔ int, float64
  - String ↔ string
  - List ↔ []interface{}
  - Boolean ↔ bool
- Go関数をLisp環境に登録
  - 型情報の保持
  - 引数の自動変換
- Lisp関数をGoから呼び出し
- 組み込みGo関数の提供
  - `go:call` - Go関数の呼び出し
  - `go:import` - Goパッケージのインポート

**使用例（Lisp側）**:
```lisp
> (go:call "fmt.Println" "Hello from Lisp")
Hello from Lisp
NIL

> (defun greet (name)
    (go:call "fmt.Printf" "Hello, %s!\n" name))
GREET

> (greet "World")
Hello, World!
NIL
```

**使用例（Go側）**:
```go
// Go関数をLispに登録
env.RegisterGoFunc("add", func(a, b int) int {
    return a + b
})

// Lispコードを実行
result := env.Eval("(add 1 2)")
fmt.Println(result) // 3

// Lisp関数をGoから呼び出し
result := env.CallLispFunc("square", 5)
fmt.Println(result) // 25
```

### M9: マクロシステム

**目標**: マクロによるメタプログラミングを可能にする

**実装内容**:
- `defmacro`の実装
  - マクロ定義
  - マクロ展開時の環境
- バッククォート `` ` ``
- カンマ `,`（unquote）
- カンマアット `,@`（unquote-splicing）
- `macroexpand`, `macroexpand-1`
- 衛生的マクロの簡易サポート（gensym）
- よく使うマクロの実装
  - `let`, `let*`
  - `push`, `pop`
  - `when`, `unless`（マクロ版）
  - `and`, `or`

**動作例**:
```lisp
> (defmacro when (condition &rest body)
    `(if ,condition (progn ,@body) nil))
WHEN

> (macroexpand '(when (> x 0) (print x) (print "positive")))
(IF (> X 0) (PROGN (PRINT X) (PRINT "positive")) NIL)

> (defmacro with-gensyms (syms &rest body)
    `(let ,(mapcar (lambda (s) `(,s (gensym))) syms)
       ,@body))
WITH-GENSYMS
```

### M10: データ構造（配列、ハッシュテーブル）

**目標**: より多くのデータ構造をサポート

**実装内容**:
- ベクター/配列
  - `#(1 2 3)` リテラル構文
  - `make-array`, `vector`
  - `aref`, `(setf aref)`
  - `vector-push`, `vector-pop`
- ハッシュテーブル
  - `make-hash-table`
  - `gethash`, `(setf gethash)`
  - `remhash`, `clrhash`
  - `hash-table-count`
- 文字列操作
  - `make-string`
  - `string-upcase`, `string-downcase`
  - `subseq`（文字列とリスト両方）
- 型システムとの統合
  - `(vector element-type)`
  - `(hash-table key-type value-type)`

**動作例**:
```lisp
> #(1 2 3 4 5)
#(1 2 3 4 5)

> (defvar arr (make-array 5))
ARR

> (setf (aref arr 0) 42)
42

> (aref arr 0)
42

> (defvar ht (make-hash-table))
HT

> (setf (gethash 'name ht) "Alice")
"Alice"

> (gethash 'name ht)
"Alice"
T
```

### M11: パッケージシステム

**目標**: 名前空間の管理を可能にする

**実装内容**:
- `defpackage`
  - パッケージ定義
  - `:use`オプション
  - `:export`オプション
- `in-package`
  - カレントパッケージの切り替え
- パッケージ修飾子
  - `package:symbol` - 外部シンボル
  - `package::symbol` - 内部シンボル
- 組み込みパッケージ
  - `common-lisp`（`cl`）
  - `keyword`
- パッケージ関連関数
  - `find-package`
  - `package-name`
  - `export`, `import`

**動作例**:
```lisp
> (defpackage :myapp
    (:use :common-lisp)
    (:export :main :run))
#<PACKAGE MYAPP>

> (in-package :myapp)
#<PACKAGE MYAPP>

> (defun main ()
    (cl:format t "Hello from myapp!~%"))
MAIN

> (in-package :cl-user)
#<PACKAGE CL-USER>

> (myapp:main)
Hello from myapp!
NIL
```

## プロジェクト構造（予定）

```
golisp/
├── SPECIFICATION.md          # この文書
├── README.md                 # プロジェクト概要
├── go.mod                    # Go modules
├── main.go                   # エントリーポイント
├── repl/
│   └── repl.go              # REPL実装
├── reader/
│   ├── lexer.go             # トークナイザ
│   └── parser.go            # パーサ
├── types/
│   ├── types.go             # S式のデータ型
│   └── type_system.go       # 型システム
├── typechecker/
│   ├── checker.go           # 静的型チェッカー
│   └── inference.go         # 型推論
├── eval/
│   ├── eval.go              # 評価器
│   ├── env.go               # 環境
│   └── builtins.go          # 組み込み関数
├── special/
│   └── special.go           # 特殊形式
├── macro/
│   └── macro.go             # マクロシステム
├── interop/
│   └── go_interop.go        # Go連携
├── package/
│   └── package.go           # パッケージシステム
└── stdlib/
    ├── list.go              # リスト関数
    ├── number.go            # 数値関数
    ├── io.go                # 入出力
    └── string.go            # 文字列関数
```

## 参考資料

- **Common Lisp**
  - Common Lisp HyperSpec: http://www.lispworks.com/documentation/HyperSpec/Front/
  - Practical Common Lisp: https://gigamonkeys.com/book/

- **型システム**
  - Typed Racket: https://docs.racket-lang.org/ts-guide/
  - TypeScript: https://www.typescriptlang.org/

- **Lisp実装**
  - Make a Lisp: https://github.com/kanaka/mal
  - Build Your Own Lisp: http://www.buildyourownlisp.com/

## 開発の進め方

1. 各マイルストーンを順番に実装
2. 各マイルストーン完了時にテストコードを作成
3. REPLで動作確認しながら進める
4. Common Lispの既存コードで動作検証

## 型システムの設計詳細

### 段階的実装アプローチ

GoLispの型システムは、TypeScriptと同様の「段階的型付け（Gradual Typing）」を採用しています。
実装も段階的に進めることで、学習効果を高めます。

#### M4a: 実行時型チェック

```
入力: (add 1 "hello")
  ↓
1. Lexer: トークン化
  ↓
2. Parser: S式に変換
  ↓
3. Evaluator: 評価開始
  - 関数 add を呼び出し
  - 引数の型をチェック ← ここでエラー検出
  - エラー: "Expected integer, got string"
```

**実装が簡単な理由**:
- 型チェックは関数呼び出し時のみ
- 既存の評価器に型チェックを追加するだけ
- S式全体を解析する必要なし

#### M4b: 静的型チェック（実行前）

```
入力: (add 1 "hello")
  ↓
1. Lexer: トークン化
  ↓
2. Parser: S式に変換
  ↓
3. Type Checker: 型チェックパス ← 新規追加
  - addの型シグネチャを確認: (integer, integer) -> integer
  - 引数1の型: integer ✓
  - 引数2の型: string ✗
  - エラー: "Expected integer, got string"
  ↓
4. Evaluator: 評価（型チェックが通った場合のみ）
```

**M4aとの違い**:
- 評価前に専用の型チェックパスを実行
- S式全体を走査して型を検証
- エラーが実行前に検出される

#### 型情報の管理

```go
// 各関数は型シグネチャを持つ
type FunctionType struct {
    ParamTypes  []Type
    ReturnType  Type
}

// 変数も型情報を持つ
type Environment struct {
    bindings  map[string]Value
    types     map[string]Type  // 変数名 → 型
}
```

### 型推論のレベル

| レベル | 内容 | M4a | M4b | M7 |
|-------|------|-----|-----|----|
| リテラル推論 | `42` → integer | ❌ | ✅ | ✅ |
| 変数推論 | `(defvar x 42)` → x: integer | ❌ | ✅ | ✅ |
| 関数戻り値推論 | `(square 5)` → integer | ❌ | ❌ | ✅ |
| 式全体の推論 | `(+ 1 2)` → integer | ❌ | ❌ | ✅ |
| ユニオン型の絞り込み | `(if (numberp x) ...)` | ❌ | ❌ | △ |

### 型なしコードとの共存

```lisp
;; 型付きコード
(defun typed-func ((x integer)) integer
  (+ x 1))

;; 型なしコード
(defun untyped-func (x)
  (+ x 1))

;; 型付き → 型なし: OK（型情報は失われる）
(untyped-func 5)  ; 実行時チェックのみ

;; 型なし → 型付き: 戻り値の型チェックが必要
(typed-func (untyped-func 5))  ; untypedの戻り値が不明
                                ; → 実行時チェックにフォールバック
```

### 段階的実装の利点

| 段階 | 実装難易度 | エラー検出タイミング | 学べること |
|------|-----------|---------------------|-----------|
| M4a | 🟢 簡単 | 実行時 | 型システムの基礎、型情報の管理 |
| M4b | 🟡 中程度 | 実行前 | 静的解析、型推論の基礎 |
| M7 | 🔴 難しい | 実行前 | 高度な型推論、ユニオン型 |

この設計により、以下を実現します：

1. **段階的学習**: 簡単なものから始めて徐々に高度な機能へ
2. **早期フィードバック**: M4aで型システムの基礎がすぐ動く
3. **実行前エラー検出**: M4bでPHP/Pythonより安全に
4. **TypeScript風の使用感**: M7で現代的な型システム
5. **実装可能な範囲**: 完全な型推論より現実的

---

最終更新: 2025-10-02