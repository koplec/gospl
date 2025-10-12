// 組み込み関数
package eval

import (
	"fmt"

	"github.com/koplec/gospl/internal/types"
)

type BuiltinFn func([]types.Expr) (types.Expr, error)

type BuiltinFunc struct {
	Name string
	Fn   BuiltinFn
}

// 組み込み関数を呼び出し
func (b BuiltinFunc) Call(args []types.Expr) (types.Expr, error) {
	return b.Fn(args)
}

func (b BuiltinFunc) String() string {
	return fmt.Sprintf("#<BUILTIN %s>", b.Name)
}

// 加算
// func([]types.Expr) (types.Expr, error)という関数系自体を型として定義できたらかっこいいかもと思ったり
// 関数定義時に型を明示することもできる。var buitinAdd BuiltinFn = func(arg []types.Expr, error){
func builtinAdd(args []types.Expr) (types.Expr, error) {
	if len(args) == 0 {
		return types.Number{Value: 0}, nil
	}

	var sum float64
	for _, arg := range args {
		num, ok := arg.(types.Number)
		if !ok {
			return nil, fmt.Errorf("+ expects numbers, got %T", arg)
		}
		sum += num.Value
	}

	return types.Number{Value: sum}, nil
}

// 下のように型アサーションを使うほうほうもあるけど、冗長
var _ BuiltinFn = builtinAdd

// 減算
func builtinSub(args []types.Expr) (types.Expr, error) {
	// common lispに準拠し、引数がないときはエラーにする
	if len(args) == 0 {
		return nil, fmt.Errorf("- requires at least 1 argument")
	}

	first, ok := args[0].(types.Number)
	if !ok {
		return nil, fmt.Errorf("- expects numbers, got %T", args[0])
	}

	if len(args) == 1 {
		//単項マイナス (- 3) -> -3
		return types.Number{Value: -first.Value}, nil
	}

	result := first.Value
	for _, arg := range args[1:] {
		num, ok := arg.(types.Number)
		if !ok {
			return nil, fmt.Errorf("- expects numbers, got %T", num)
		}
		result -= num.Value
	}

	return types.Number{Value: result}, nil
}

// 乗算
func builtinMul(args []types.Expr) (types.Expr, error) {
	if len(args) == 0 {
		return types.Number{Value: 1.0}, nil
	}

	result := 1.0
	for _, arg := range args {
		num, ok := arg.(types.Number)
		if !ok {
			return nil, fmt.Errorf("* expects numbers, got %T", num)
		}
		result *= num.Value
	}

	return types.Number{Value: result}, nil
}

// 除算
func builtinDiv(args []types.Expr) (types.Expr, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("/ requires at least 1 argument")
	}

	first, ok := args[0].(types.Number)
	if !ok {
		return nil, fmt.Errorf("/ expects numbers, got %T", args[0])
	}

	if len(args) == 1 {
		// 逆数
		if first.Value == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return types.Number{Value: 1.0 / first.Value}, nil
	}

	result := first.Value
	for _, arg := range args[1:] {
		num, ok := arg.(types.Number)
		if !ok {
			return nil, fmt.Errorf("/ expects numbers, got %T", arg)
		}
		if num.Value == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result /= num.Value
	}

	return types.Number{Value: result}, nil
}
