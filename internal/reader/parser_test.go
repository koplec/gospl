package reader

import (
	"testing"

	"github.com/koplec/gospl/internal/types"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42", 42.0},
		{"3.14", 3.14},
		{"-10", -10.0},
		{"0", 0.0},
		{"-3.14", -3.14},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parser := NewParser(tt.input)
			expr, err := parser.Parse()

			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			num, ok := expr.(types.Number)
			if !ok {
				t.Fatalf("expected Number, got %T", expr)
			}

			if num.Value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, num.Value)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello, world"`, "hello, world"},
		{`""`, ""},
		//{`"multiple\nlines`, "multiple\\nlines"}, //エスケープ
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parser := NewParser(tt.input)
			expr, err := parser.Parse()

			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			str, ok := expr.(types.String)
			if !ok {
				t.Fatalf("expected String, got %T", expr)
			}

			if str.Value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, str.Value)
			}
		})
	}
}

func TestParseSymbol(t *testing.T) {
	tests := []string{
		"x",
		"foo",
		"+",
		"-",
		"*",
		"/",
		"my-func",
		"defun",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			parser := NewParser(input)
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			sym, ok := expr.(types.Symbol)
			if !ok {
				t.Fatalf("expected Symbol, got %T", expr)
			}

			if sym.Name != input {
				//Fatalfだとテストが終わるから、Errorfする
				t.Errorf("expected %q, got %q", input, sym.Name)
			}

		})
	}
}

func TestParseBoolean(t *testing.T) {
	tests := []struct {
		input    string
		isTrue   bool
		expected string
	}{
		{"t", true, "T"},
		{"nil", false, "NIL"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parser := NewParser(tt.input)
			expr, err := parser.Parse()

			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			if tt.isTrue {
				b, ok := expr.(types.Boolean)
				if !ok {
					t.Fatalf("expected Boolean, got %T", expr)
				}
				if !b.Value {
					t.Errorf("expected true, got false")
				}
			} else {
				_, ok := expr.(*types.Nil)
				if !ok {
					t.Fatalf("expected Nil, got %T", expr)
				}
			}

			if expr.String() != tt.expected {
				//%sはそのまま出力だけど、%qはGo文字列リテラルとして出力　例えばダブルクォート　エスケープがつく
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}

		})
	}
}

func TestParseEmptyList(t *testing.T) {
	parser := NewParser("()")
	expr, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := expr.(*types.Nil)
	if !ok {
		t.Fatalf("expected Nil, got %T", expr)
	}
	if expr.String() != "NIL" {
		t.Errorf("expected NIL string , got %s", expr.String())
	}
}

func TestParseSimpleList(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(1)", "(1)"},
		{"(1   2 3)", "(1 2 3)"},
		{"(+  1      2)", "(+ 1 2)"},
		{"(foo bar baz)", "(foo bar baz)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parser := NewParser(tt.input)
			expr, err := parser.Parse()

			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			_, ok := expr.(*types.Cons)
			if !ok {
				t.Fatalf("expected Cons, got %T", expr)
			}

			//パースして文字列化したときに正規化されていて、構造が崩れていないことを確認する
			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseNestedList(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(1 (2 3))", "(1 (2 3))"},
		{"((1 2) (3 4))", "((1 2) (3 4))"},
		{"(1 (2 (3 4)))", "(1 (2 (3 4)))"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parser := NewParser(tt.input)
			expr, err := parser.Parse()

			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			if expr.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'x", "(quote x)"},
		{"'123", "(quote 123)"},
		{"'(1 2 3)", "(quote (1 2 3))"},
		{"'()", "(quote NIL)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parser := NewParser(tt.input)
			expr, err := parser.Parse()

			if err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			if tt.expected != expr.String() {
				t.Errorf("expected %s, got %s", tt.expected, expr.String())
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unclosed list", "(1 2 3"},
		{"unexpected closing paren", ")"},
		//これはテスト通らなかった。Parseの責務が一つの式を読む責務なので。
		//{"unexpected closing paren in list", "(1 2))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			expr, err := parser.Parse()
			if err == nil {
				t.Fatalf("expected err, got nil, expected string:%s", expr.String())
			}
		})
	}
}

func TestParseMultilineString(t *testing.T) {
	input := `"hello
world"`
	parser := NewParser(input)
	expr, err := parser.Parse()

	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}

	str, ok := expr.(types.String)
	if !ok {
		t.Fatalf("expected String, got %T", expr)
	}

	expected := "hello\nworld"
	if str.Value != expected {
		// やっぱり文字列の時にどうして%qを使うのかしら
		t.Errorf("expected %q, got %q", expected, str.Value)
	}
}
