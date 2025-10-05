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
