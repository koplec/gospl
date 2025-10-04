package reader

import "fmt"

type TokenType int

const (
	LPAREN TokenType = iota // (
	RPAREN                  // )
	NUMBER                  // 123, 3.14
	STRING                  // "hello"
	SYMBOL                  // foo, +, defun
	QUOTE                   // '
	EOF
	ILLEGAL
)

// Lexerがソースコードを読んで生成するトークン
type Token struct {
	Type  TokenType
	Value string
	Pos   Position
}

// Lexerが読んでいるソースコード上の位置
type Position struct {
	Line   int
	Column int
}

// 字句解析機
// プログラムのソースコードをTOKENに分解する
type Lexer struct {
	input  string
	pos    int // inputの中で現在読んでいる位置、input全体の中でどの位置か
	line   int // 現在読んでいる行番号（1始まり）
	column int // 現在読んでいる列番号（1始まり）
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) NextToken() (Token, error) {
	//まず空白が1文字以上あるときは、スキップする
	l.skipWhitespace()

	//もし入力が終わっているときはEOFを返す
	if l.pos >= len(l.input) {
		return Token{Type: EOF,
			Value: "",
			Pos:   Position{Line: l.line, Column: l.column},
		}, nil
	}

	ch := l.input[l.pos]
	pos := l.currentPos()

	switch ch {
	case '(':
		l.advance()
		return Token{Type: LPAREN, Value: "(", Pos: pos}, nil
	case ')':
		l.advance()
		return Token{Type: RPAREN, Value: ")", Pos: pos}, nil
	case '\'': //quote
		l.advance()
		return Token{Type: QUOTE, Value: "'", Pos: pos}, nil
	case '"':
		return l.readString()
	}

	//数値リテラルの判定、　数字または-で始まる場合は先に数字があるはず
	if isDigit(ch) || (ch == '-' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1])) {
		return l.readNumber()
	}

	//それ以外はシンボルとして読む
	if isSymbolStart(ch) {
		return l.readSymbol()
	}
	// ここまで当たらないということはエラー
	return Token{
		Type:  ILLEGAL,
		Value: string(ch),
		Pos:   pos,
	}, fmt.Errorf("unexpected character: %c at line %d, column %d", ch, pos.Line, pos.Column)
}

/**
 * 空白をスキップする
 * 空白とは、スペース、タブ、改行など
 * 空白が1文字以上続くときはずっとスキップする
 * 空白が1文字もないときは何もしない
 */
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.pos++
			if ch == '\n' {
				l.line++ // 縦に移動
				l.column = 1
			} else {
				l.column++ //右に移動
			}
		} else {
			// 空白以外の文字が来たときは何もしない
			break
		}
	}
}

func (l *Lexer) currentPos() Position {
	return Position{Line: l.line, Column: l.column}
}

// 1文字読み進める
// posだけでなく、columnも進める
// 改行はされないことに注意（なので、l.lineは進めない）
func (l *Lexer) advance() {
	l.pos++
	l.column++
}

func (l *Lexer) readString() (Token, error) {
	pos := l.currentPos()
	l.advance() // 最初の"をスキップする

	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		// 改行が来てcommon lispのように複数行にわたる文字列リテラルを許す
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++ //下位行以外の文字を読むたびに列を進める
		}
		l.pos++
	}

	//文字列が長いときエラー
	if l.pos >= len(l.input) {
		return Token{Type: ILLEGAL, Value: "", Pos: pos}, fmt.Errorf("unterminated string at line %d, column %d", pos.Line, pos.Column)
	}

	//閉じていることを確認する
	//// for loopで確認できているはずなので、省略してもよさそう
	ch := l.input[l.pos]
	if ch != '"' {
		return Token{Type: ILLEGAL, Value: "", Pos: pos}, fmt.Errorf("unclosed string")
	}

	// 閉じる"を見つけた
	value := l.input[start:l.pos]
	l.advance() //閉じるをスキップ, value取得前に行うと"まで読んでしまう
	return Token{Type: STRING, Value: value, Pos: pos}, nil
}

func (l *Lexer) readNumber() (Token, error) {
	pos := l.currentPos()
	start := l.pos

	//マイナス記号があれば先へ
	if l.input[l.pos] == '-' {
		l.advance()
	}

	//整数部分を読む
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.advance()
	}

	//小数点があれば浮動小数点
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		l.advance()
		//小数部分を読む
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.advance()
		}
	}

	value := l.input[start:l.pos]
	return Token{Type: NUMBER, Value: value, Pos: pos}, nil
}

func (l *Lexer) readSymbol() (Token, error) {
	pos := l.currentPos()
	start := l.pos

	//シンボルの最初の文字を読む
	l.advance()

	//シンボルの残りの文字を読む
	for l.pos < len(l.input) && isSymbolChar(l.input[l.pos]) {
		l.advance()
	}

	value := l.input[start:l.pos]
	return Token{Type: SYMBOL, Value: value, Pos: pos}, nil
}

// helper関数
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// commonlispのシンボルで使える文字を先頭にしたらsymbolとする
func isSymbolStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '+' || ch == '-' || ch == '*' || ch == '/' ||
		ch == '=' || ch == '<' || ch == '>' || ch == '!'
}

func isSymbolChar(ch byte) bool {
	return isSymbolStart(ch) || isDigit(ch) ||
		// ch == '-' // hyphenの意味だけど、isSymbolStartにこの文字の判定は含まれるのでは？
		ch == '_'
}
