# GoLisp アーキテクチャ設計書（M1）

## 目的

このドキュメントは、GoLisp M1（基本的なS式とREPL）の技術選定とアーキテクチャ設計を記録します。

## 1. データ構造の設計

### 1.1 S式の表現

**採用方式: interface + 具体型**

```go
// internal/types/types.go
package types

// Expr はすべてのLisp式の基底インターフェース
type Expr interface {
    String() string
}

// Number は数値を表現（M1では整数と浮動小数点を統一）
type Number struct {
    Value float64
}

// Symbol はシンボルを表現
type Symbol struct {
    Name string
}

// Cons はコンスセルを表現（リストの基本構造）
type Cons struct {
    Car Expr  // リストの先頭要素
    Cdr Expr  // リストの残り（ConsまたはNil）
}

// String は文字列を表現
type String struct {
    Value string
}

// Boolean は真偽値を表現
type Boolean struct {
    Value bool
}

// Nil は空リストとfalseを表現
type Nil struct{}
```

**設計判断:**
- **型安全性**: interface{}ではなく、具体的な型を定義
- **拡張性**: 新しい型を追加しやすい
- **Common Lisp準拠**: コンスセルベースのリスト表現

### 1.2 特殊な値

```go
var (
    Nil  = &Nil{}       // 空リスト / false
    True = Boolean{Value: true}  // t
)
```

### 1.3 リストの構造

**純粋なコンスセル方式を採用**

```
(1 2 3) の内部表現:

Cons{
  Car: Number{1},
  Cdr: Cons{
    Car: Number{2},
    Cdr: Cons{
      Car: Number{3},
      Cdr: Nil
    }
  }
}
```

**利点:**
- Common Lispの意味論に忠実
- `car`/`cdr`の実装が自然
- improper list `(1 . 2)` もサポート可能

## 2. Lexer（字句解析器）の設計

### 2.1 実装方式

**手書きLexerを採用**

```go
// internal/reader/lexer.go
package reader

type TokenType int

const (
    LPAREN   TokenType = iota  // (
    RPAREN                      // )
    NUMBER                      // 123, 3.14
    STRING                      // "hello"
    SYMBOL                      // foo, +, defun
    QUOTE                       // '
    EOF
    ILLEGAL
)

type Token struct {
    Type  TokenType
    Value string
    Pos   Position
}

type Position struct {
    Line   int
    Column int
}

type Lexer struct {
    input   string
    pos     int     // 現在の位置
    line    int
    column  int
}

func NewLexer(input string) *Lexer
func (l *Lexer) NextToken() (Token, error)
```

**設計判断:**
- エラーメッセージに位置情報を含める
- 1文字先読み（peek）で判断
- シンプルで理解しやすい実装

### 2.2 トークン化のルール

| 入力 | トークン | 備考 |
|------|---------|------|
| `(` | LPAREN | |
| `)` | RPAREN | |
| `123` | NUMBER("123") | 整数 |
| `3.14` | NUMBER("3.14") | 浮動小数点 |
| `"hello"` | STRING("hello") | エスケープは後で実装 |
| `foo` | SYMBOL("foo") | |
| `+` | SYMBOL("+") | 演算子もシンボル |
| `'` | QUOTE | `(quote x)`の糖衣構文 |
| `;comment` | スキップ | 行末まで |
| 空白/改行 | スキップ | |

## 3. Parser（構文解析器）の設計

### 3.1 実装方式

**再帰下降パーサを採用**

```go
// internal/reader/parser.go
package reader

type Parser struct {
    lexer   *Lexer
    current Token
}

func NewParser(input string) *Parser
func (p *Parser) Parse() (types.Expr, error)
func (p *Parser) parseExpr() (types.Expr, error)
func (p *Parser) parseList() (types.Expr, error)
```

### 3.2 文法（BNF風）

```
expr   ::= atom | list | quoted
atom   ::= NUMBER | STRING | SYMBOL
list   ::= '(' expr* ')'
quoted ::= "'" expr
```

### 3.3 パース例

**入力:** `(+ 1 2)`

```
1. Lexer:
   [LPAREN, SYMBOL("+"), NUMBER("1"), NUMBER("2"), RPAREN]

2. Parser:
   parseExpr() → parseList()
     - '(' を消費
     - parseExpr() → Symbol{"+"}
     - parseExpr() → Number{1}
     - parseExpr() → Number{2}
     - ')' を消費

3. 結果:
   Cons{
     Car: Symbol{"+"},
     Cdr: Cons{
       Car: Number{1},
       Cdr: Cons{
         Car: Number{2},
         Cdr: Nil
       }
     }
   }
```

## 4. Printer（出力）の設計

### 4.1 実装方式

```go
// internal/types/printer.go
package types

func (n Number) String() string {
    return fmt.Sprintf("%g", n.Value)
}

func (s Symbol) String() string {
    return s.Name
}

func (c *Cons) String() string {
    // リスト形式で出力: (a b c)
    result := "("
    for c != nil {
        if cons, ok := c.Cdr.(*Cons); ok {
            result += c.Car.String() + " "
            c = cons
        } else if _, ok := c.Cdr.(*Nil); ok {
            result += c.Car.String()
            break
        } else {
            // improper list: (a . b)
            result += c.Car.String() + " . " + c.Cdr.String()
            break
        }
    }
    return result + ")"
}

func (s String) String() string {
    return fmt.Sprintf("\"%s\"", s.Value)
}

func (b Boolean) String() string {
    if b.Value {
        return "T"
    }
    return "NIL"
}

func (*Nil) String() string {
    return "NIL"
}
```

## 5. REPL の設計

### 5.1 実装方式

```go
// internal/repl/repl.go
package repl

import (
    "bufio"
    "fmt"
    "os"
    "your-module-name/internal/reader"
    "your-module-name/internal/types"
)

func Start() {
    scanner := bufio.NewScanner(os.Stdin)

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

        // Eval (M1ではスキップ、そのまま返す)
        result := expr

        // Print
        fmt.Println(result.String())
    }
}
```

**M1の動作:**
- 入力をそのまま出力（Evalは恒等関数）
- エラー時も継続実行
- Ctrl+D で終了

## 6. エラーハンドリング

### 6.1 エラー型の設計

```go
// internal/reader/error.go
package reader

type ErrorType int

const (
    LexError ErrorType = iota
    ParseError
)

type Error struct {
    Type    ErrorType
    Pos     Position
    Message string
}

func (e *Error) Error() string {
    return fmt.Sprintf("%s at line %d, col %d: %s",
        e.Type, e.Pos.Line, e.Pos.Column, e.Message)
}
```

### 6.2 エラーハンドリング戦略

1. **エラーの伝播**: すべての関数は `(result, error)` を返す
2. **位置情報の保持**: エラーメッセージに行・列番号を含める
3. **パニックを避ける**: Goのidiomatic wayに従う
4. **REPL での継続**: エラーが発生してもREPLは継続

## 7. モジュール構成

### 7.1 ディレクトリ構造

```
golisp/
├── docs/
│   ├── SPECIFICATION.md
│   ├── design/
│   │   └── M1-architecture.md
│   └── notes.md
├── go.mod
├── main.go
└── internal/
    ├── types/
    │   ├── types.go
    │   └── printer.go
    ├── reader/
    │   ├── lexer.go
    │   ├── parser.go
    │   └── error.go
    └── repl/
        └── repl.go
```

### 7.2 依存関係

```
main.go
  ↓
internal/repl/repl.go
  ↓
internal/reader/parser.go → internal/reader/lexer.go
  ↓                               ↓
internal/types/types.go ← internal/types/printer.go
```

**依存ルール:**
- `internal/types` は他のどのパッケージにも依存しない
- `internal/reader` は `internal/types` のみに依存
- `internal/repl` は `internal/reader` と `internal/types` に依存
- 循環依存を避ける

## 8. M1完成の定義

### 8.1 動作する機能

- [ ] S式の読み込み
- [ ] ネストしたリストの解析
- [ ] 数値、文字列、シンボルの認識
- [ ] REPLの起動と終了
- [ ] エラーメッセージの表示

### 8.2 テスト例

```lisp
> 123
123

> "hello"
"hello"

> foo
FOO

> (+ 1 2)
(+ 1 2)

> (list 1 2 3)
(LIST 1 2 3)

> '(1 2 3)
(QUOTE (1 2 3))

> (nested (list structure))
(NESTED (LIST STRUCTURE))
```

### 8.3 未実装の機能（M2以降）

- 式の評価（Eval）
- 算術演算
- 変数の束縛
- 関数定義

## 9. 今後の拡張ポイント

### 9.1 M2への準備

M2では評価器を実装するため、以下を追加予定:

```go
// eval/eval.go
func Eval(expr types.Expr, env *Environment) (types.Expr, error)

// eval/env.go
type Environment struct {
    bindings map[string]types.Expr
    parent   *Environment
}
```

### 9.2 型システムへの準備

M4で型システムを追加するため、以下の拡張が必要:

```go
// types/type_system.go
type Type interface {
    Name() string
}

// 各Exprに型情報を追加
type TypedExpr interface {
    Expr
    Type() Type
}
```

## 10. 参考資料

- **Make a Lisp**: https://github.com/kanaka/mal
- **Build Your Own Lisp**: http://www.buildyourownlisp.com/
- **Writing An Interpreter In Go**: https://interpreterbook.com/

---

**作成日**: 2025-10-04
**対象マイルストーン**: M1（基本的なS式とREPL）
