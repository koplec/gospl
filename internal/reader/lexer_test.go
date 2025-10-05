package reader

import "testing"

func TestLexer_SingleToken(t *testing.T) {
	lexer := NewLexer("(")

	token, err := lexer.NextToken()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.Type != LPAREN {
		t.Errorf("expected LPAREN, got %v", token.Type)
	}

	if token.Value != "(" {
		t.Errorf("expected '(', got %v", token.Value)
	}

	if token.Pos.Line != 1 || token.Pos.Column != 1 {
		t.Errorf("expected position(1,1), got (%d,%d)", token.Pos.Line, token.Pos.Column)
	}
}

func TestLexer_BasicTokens(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  TokenType
		wantValue string
	}{
		{"left paren", "(", LPAREN, "("},
		{"right paren", ")", RPAREN, ")"},
		{"quote", "'", QUOTE, "'"},
		{"positive number", "123", NUMBER, "123"},
		{"negative number", "-123", NUMBER, "-123"},
		{"float", "3.14", NUMBER, "3.14"},
		{"symbol", "foo", SYMBOL, "foo"},
		{"operator plus", "+", SYMBOL, "+"},
		{"operator minus alone", "-", SYMBOL, "-"},
		{"string", `"hello"`, STRING, "hello"},
		{"empty string", `""`, STRING, ""},
	}

	for _, tt := range tests {
		//Goのサブテストの機能 testnameごとに実行とかできる
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			token, err := lexer.NextToken()

			if err != nil {
				t.Fatalf("unexpected error : %v", err)
			}

			if token.Type != tt.wantType {
				t.Errorf("Type: expected:%v, got %v", tt.wantType, token.Type)
			}

			if token.Value != tt.wantValue { //GOの文字列は値型として比較されるので、==で比較できる
				t.Errorf("Value: expected %q, got %q", tt.wantValue, token.Value)
			}
		})

	}
}

func TestLexer_MultipleTokens(t *testing.T) {
	input := "(+ 1 2)"
	lexer := NewLexer(input)

	expected := []struct {
		tokenType TokenType
		value     string
	}{
		{LPAREN, "("},
		{SYMBOL, "+"},
		{NUMBER, "1"},
		{NUMBER, "2"},
		{RPAREN, ")"},
		{EOF, ""},
	}

	for i, want := range expected {
		token, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error:%v", i, err)
		}

		if token.Type != want.tokenType {
			t.Errorf("token %d: Type = %v, want %v", i, token.Type, want.tokenType)
		}

		if token.Value != want.value {
			t.Errorf("token %d: Value = %q, want %q", i, token.Value, want.value)
		}
	}
}

func TestLexer_WhitespaceSkipTest(t *testing.T) {
	input := "   (   +    1   )"
	lexer := NewLexer(input)

	expected := []struct {
		tokenType TokenType
		value     string
	}{
		{LPAREN, "("},
		{SYMBOL, "+"},
		{NUMBER, "1"},
		{RPAREN, ")"},
		{EOF, ""},
	}

	for i, want := range expected {
		token, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error:%v", i, err)
		}

		if token.Type != want.tokenType {
			t.Errorf("token %d: Type = %v, want %v", i, token.Type, want.tokenType)
		}

		if token.Value != want.value {
			t.Errorf("token %d: Value = %q, want %q", i, token.Value, want.value)
		}
	}
}

// 改行を含んだ文字列も行ける
func TestLexer_MultilineString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue string
	}{
		{
			name:      "string with newline",
			input:     "\"hello\nworld\"",
			wantValue: "hello\nworld",
		},
		{
			name:      "string with multiple newlines",
			input:     "\"line1\nline2\nline3\"",
			wantValue: "line1\nline2\nline3",
		},
		{
			name: "multiline string literal",
			input: `"first line
second line
third line"`,
			wantValue: "first line\nsecond line\nthird line",
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			token, err := lexer.NextToken()

			if err != nil {
				t.Fatalf("test case %d unexpected error: %v", i, err)
			}

			if token.Type != STRING {
				t.Errorf("test case %d Type = %v, want STRING", i, token.Type)
			}

			if token.Value != tt.wantValue {
				t.Errorf("test case %d Value=%q, want=%q", i, token.Value, tt.wantValue)
			}
		})
	}
}

// 位置情報のテスト
func TestLexer_Position(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []struct {
			typ  TokenType
			line int
			col  int
		}
	}{
		{
			name:  "single line positions",
			input: "(+ 1)",
			tokens: []struct {
				typ  TokenType
				line int
				col  int
			}{
				{LPAREN, 1, 1},
				{SYMBOL, 1, 2},
				{NUMBER, 1, 4},
				{RPAREN, 1, 5},
				{EOF, 1, 6},
			},
		},
		{
			name: "multiline positions",
			input: `(+
1
2)`,
			tokens: []struct {
				typ  TokenType
				line int
				col  int
			}{
				{LPAREN, 1, 1},
				{SYMBOL, 1, 2},
				{NUMBER, 2, 1},
				{NUMBER, 3, 1},
				{RPAREN, 3, 2},
				{EOF, 3, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			for i, want := range tt.tokens {
				token, err := lexer.NextToken()
				if err != nil {
					t.Fatalf("token %d: unexpected error: %v", i, err)
				}

				if token.Pos.Line != want.line ||
					token.Pos.Column != want.col {
					t.Errorf("token %d: position=(%d, %d), want (%d, %d)",
						i, token.Pos.Line, token.Pos.Column, want.line, want.col)
				}
			}
		})
	}
}

func TestLexer_Quote(t *testing.T) {
	input := "'(1 2 3)"
	lexer := NewLexer(input)

	expected := []TokenType{
		QUOTE,
		LPAREN, NUMBER, NUMBER, NUMBER, RPAREN,
		EOF,
	}

	for i, want := range expected {
		token, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if token.Type != want {
			t.Errorf("token %d: Type=%v, want=%v", i, token.Type, want)
		}
	}
}

func TestLexer_ComplexSymbols(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue string
	}{
		{"symbol with hyphen", "foo-bar", "foo-bar"},
		{"symbol with number", "var1", "var1"},
		{"symbol with underscore", "foo_bar", "foo_bar"},
		{"defun", "defun", "defun"},
		{"lambda", "lambda", "lambda"},
		{"asterisc operator", "*", "*"},
		{"division operator", "/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			token, err := lexer.NextToken()

			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			if token.Type != SYMBOL {
				t.Errorf("Type=%v, want SYMBOL", token.Type)
			}

			if token.Value != tt.wantValue {
				t.Errorf("Value=%q, want %q", token.Value, tt.wantValue)
			}
		})
	}
}

// 認識できない文字
// common lispのマクロ文字に対応
func TestLexer_UnrecoginizedCharacter(t *testing.T) {
	tests := []string{
		"@", "#", "$", "%",
	}

	for _, input := range tests {
		t.Run("illegal char: "+input, func(t *testing.T) {
			lexer := NewLexer(input)
			token, err := lexer.NextToken()

			if err == nil {
				t.Errorf("expected error for input %q", input)
			}
			if token.Type != ILLEGAL {
				t.Errorf("expected ILLEGAL token, got %v", token.Type)
			}
		})
	}
}

// エラーケース

func TestLexer_UnterminatedString(t *testing.T) {
	lexer := NewLexer(`"hello`)
	token, err := lexer.NextToken()

	if err == nil {
		t.Fatal("expected error for unterminated string")
	}

	if token.Type != ILLEGAL {
		t.Errorf("expected ILLEGAL token, got %v", token.Type)
	}
}
