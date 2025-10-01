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

- **実行時型チェック**: 型宣言がある場合、実行時に型をチェック
- **型推論**: 可能な範囲で型を推論
- **段階的導入**: 一部の関数だけ型付けすることも可能

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

### M4: 基本的な型システム

**目標**: 型宣言と実行時型チェックができる

**実装内容**:
- 型の内部表現
  - Type インターフェース
  - プリミティブ型（Integer, Float, String, Boolean, Symbol）
  - 複合型（List, Cons）
- `declare`フォームの実装
- 実行時型チェック
  - 関数引数の型チェック
  - 戻り値の型チェック
- 型エラーの報告
- 型述語の実装
  - `numberp`, `symbolp`, `listp`, `stringp`, `atom`, `null`

**動作例**:
```lisp
> (defun typed-add ((x integer) (y integer)) integer
    (+ x y))
TYPED-ADD

> (typed-add 1 2)
3

> (typed-add 1 "hello")
Error: Type error: expected integer, got string

> (numberp 42)
T

> (symbolp 'foo)
T
```

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

### M7: 型推論とユニオン型

**目標**: より高度な型機能を実装

**実装内容**:
- ユニオン型 `(or type1 type2 ...)`
- 関数型 `(function (arg-types...) return-type)`
- リスト型 `(list element-type)`
- 簡易的な型推論
  - リテラルからの型推論
  - 関数戻り値の型推論
- 型アノテーション付きの変数定義
- 型エラーメッセージの改善

**動作例**:
```lisp
> (defvar x (or integer string))
X

> (setq x 42)
42

> (setq x "hello")
"hello"

> (setq x 3.14)
Error: Type error: expected (or integer string), got float

> (defun maybe-number ((x (or integer nil))) (or integer nil)
    (if x (+ x 1) nil))
MAYBE-NUMBER

> (maybe-number 5)
6

> (maybe-number nil)
NIL
```

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

---

最終更新: 2025-10-01