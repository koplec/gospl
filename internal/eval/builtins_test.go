package eval

import (
	"testing"

	"github.com/koplec/gospl/internal/types"
)

func TestBuiltinAdd(t *testing.T) {
	tests := []struct {
		name string
		args []types.Expr
		want types.Number
	}{
		{"noargs", []types.Expr{}, types.Number{Value: 0}},
		{
			"single",
			[]types.Expr{
				types.Number{Value: 5},
			},
			types.Number{Value: 5},
		},
		{
			"multiple",
			[]types.Expr{
				types.Number{Value: 1},
				types.Number{Value: 2},
				types.Number{Value: 3},
			},
			types.Number{Value: 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builtinAdd(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			num, ok := result.(types.Number)
			if !ok {
				t.Fatalf("unexpected type: %T", num)
			}
			if num.Value != tt.want.Value {
				t.Errorf("got %v, want %v", num.Value, tt.want.Value)
			}
		})
	}
}

func TestBuiltinAdd_TypeError(t *testing.T) {
	result, err := builtinAdd([]types.Expr{
		types.Number{Value: 1},
		types.String{Value: "hello"},
	})
	if err == nil {
		t.Fatal("expected type error")
	}
	if result != nil {
		t.Errorf("unexpected result:%v", result)
	}
}

func TestBuiltinSub(t *testing.T) {
	tests := []struct {
		name string
		args []types.Expr
		want types.Number
	}{
		{
			"single",
			[]types.Expr{
				types.Number{Value: 5},
			},
			types.Number{Value: -5},
		},
		{
			"multiple",
			[]types.Expr{
				types.Number{Value: 1},
				types.Number{Value: 2},
				types.Number{Value: 3},
			},
			types.Number{Value: -4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builtinSub(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			num, ok := result.(types.Number)
			if !ok {
				t.Fatalf("unexpected type: %T", num)
			}
			if num.Value != tt.want.Value {
				t.Errorf("got %v, want %v", num.Value, tt.want.Value)
			}
		})
	}
}

func TestBuiltinSub_No_Args(t *testing.T) {
	result, err := builtinSub([]types.Expr{})
	if err == nil {
		t.Fatalf("expected no error result:%v", result)
	}
}

func TestBuiltinMul(t *testing.T) {
	tests := []struct {
		name string
		args []types.Expr
		want types.Number
	}{
		{"noargs", []types.Expr{}, types.Number{Value: 1}},
		{
			"single",
			[]types.Expr{
				types.Number{Value: 5},
			},
			types.Number{Value: 5},
		},
		{
			"multiple",
			[]types.Expr{
				types.Number{Value: 1},
				types.Number{Value: 2},
				types.Number{Value: 5},
			},
			types.Number{Value: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builtinMul(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			num, ok := result.(types.Number)
			if !ok {
				t.Fatalf("unexpected type: %T", num)
			}
			if num.Value != tt.want.Value {
				t.Errorf("got %v, want %v", num.Value, tt.want.Value)
			}
		})
	}
}

func TestBuiltinDiv(t *testing.T) {
	tests := []struct {
		name string
		args []types.Expr
		want types.Number
	}{
		{
			"single",
			[]types.Expr{
				types.Number{Value: 5},
			},
			types.Number{Value: 1.0 / 5.0},
		},
		{
			"multiple",
			[]types.Expr{
				types.Number{Value: 1},
				types.Number{Value: 2},
				types.Number{Value: 5},
			},
			types.Number{Value: 0.1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builtinDiv(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			num, ok := result.(types.Number)
			if !ok {
				t.Fatalf("unexpected type: %T", num)
			}
			if num.Value != tt.want.Value {
				t.Errorf("got %v, want %v", num.Value, tt.want.Value)
			}
		})
	}
}

func TestBuitinDiv_No_Args(t *testing.T) {
	result, err := builtinDiv([]types.Expr{})
	if err == nil {
		t.Fatalf("unexpected nil")
	}
	if result != nil {
		t.Errorf("unexpected result:%v", result)
	}
}

func TestBuitlinDiv_Zero_Division(t *testing.T) {
	result, err := builtinDiv([]types.Expr{
		types.Number{Value: 5},
		types.Number{Value: 0},
	})
	if err == nil {
		t.Fatal("expected division by zero error")
	}
	if result != nil {
		t.Errorf("unexpected result: %v", result)
	}
}
