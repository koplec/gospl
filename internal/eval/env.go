package eval

import (
	"fmt"

	"github.com/koplec/gospl/internal/types"
)

// 変数の束縛の管理
type Environment struct {
	bindings map[string]types.Expr
	parent   *Environment //親環境、スコープチェーンに利用
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		bindings: make(map[string]types.Expr),
		parent:   parent,
	}
}

func NewGlobalEnvironment() *Environment {
	env := NewEnvironment(nil)

	env.Set("+", BuiltinFunc{Name: "+", Fn: builtinAdd})
	env.Set("-", BuiltinFunc{Name: "-", Fn: builtinSub})
	env.Set("*", BuiltinFunc{Name: "*", Fn: builtinMul})
	env.Set("/", BuiltinFunc{Name: "/", Fn: builtinDiv})

	return env
}

func (e *Environment) Set(name string, value types.Expr) {
	e.bindings[name] = value
}

func (e *Environment) Get(name string) (types.Expr, error) {
	//現在の環境で探す
	if val, ok := e.bindings[name]; ok {
		return val, nil
	}

	//親環境で探す
	if e.parent != nil {
		return e.parent.Get(name)
	}

	//なかった。。。
	return nil, fmt.Errorf("undefined variable: %s", name)
}
