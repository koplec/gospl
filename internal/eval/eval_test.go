package eval

import (
	"testing"

	"github.com/koplec/gospl/internal/reader"
)

func TestEval(t *testing.T) {
	env := NewGlobalEnvironment()

	tests := []struct {
		name  string
		input string
		want  string //結果の文字列表現
	}{
		{"number", "42", "42"},
		{"string", `"hello"`, `"hello"`},
		{"add", "(+ 1 2)", "3"},
		{"nested", "(+ (* 2 3) (/ 10 5))", "8"},
		{"subtract", "(- 10 3)", "7"},
		{"multiply", "(* 3 4)", "12"},
		{"divide", "(/ 10 2)", "5"},
		{"unary minus", "(- 5)", "-5"},
		{"minus", "-10", "-10"},
		{"quote symbol", "(quote x)", "x"},
		{"quote list", "(quote (1 2 3))", "(1 2 3)"},
		{"quote nested list", "(quote ((1 2) (3 4)))", "((1 2) (3 4))"},
		{"quote expression", "(quote (+ 1 2 3))", "(+ 1 2 3)"},
		{"quote number", "(quote 123)", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := reader.NewParser(tt.input)
			expr, err := parser.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			result, err := Eval(expr, env)
			if err != nil {
				t.Fatalf("eval error: %v", err)
			}
			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

// エラーケースのテスト
func TestEval_Errors(t *testing.T) {
	env := NewGlobalEnvironment()

	tests := []struct {
		name  string
		input string
	}{
		{"undefined variable", "x"},
		{"type error", `(+ 1 "hello")`},
		{"division by zero", "(/ 5 0)"},
		{"not a function", "(42 1 2)"},
		{"quote no args", "(quote)"},
		{"quote too many args", "(quote x y)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := reader.NewParser(tt.input)
			expr, err := parser.Parse()
			// parseはできる
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			_, err = Eval(expr, env)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}
