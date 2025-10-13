package eval

import (
	"fmt"

	"github.com/koplec/gospl/internal/types"
)

func Eval(expr types.Expr, env *Environment) (types.Expr, error) {
	switch e := expr.(type) {
	case types.Number:
		//数値はそのまま返す
		return e, nil
	case types.String:
		//文字列もそのまま
		return e, nil
	case types.Boolean:
		//真偽値もそのまま
		return e, nil
	case *types.Nil:
		return e, nil

	case types.Symbol:
		//シンボルは環境から値を取得
		return env.Get(e.Name)
	case *types.Cons:
		//リストは関数適用
		return evalList(e, env)
	default:
		return nil, fmt.Errorf("unknown expression type:%T", expr)
	}
}

// リスト（関数適用）を評価
func evalList(list *types.Cons, env *Environment) (types.Expr, error) {
	// Goのnilポインタチェック（通常は発生しないはず、防衛的に記述）
	// 空リストは*types.Nil{}として表現されるので、Eval内のcase *types.Nilで対応しているため、個々には到達しないはず
	if list == nil {
		return nil, fmt.Errorf("cannot evaluate empty list")
	}

	//先頭要素を取得
	//例えば、(hoge bar baz)のhoge
	first := list.Car

	//シンボルなら、special formかどうかを確認
	if sym, ok := first.(types.Symbol); ok {
		if isSpecialForm(sym.Name) {
			//list.Cdrについて
			//もとのlistが(hoge bar baz)だったら(bar baz)が渡される
			//quoteの時は難しくて、(quote (a b c))だったら((a b c))が渡される。
			//(quote x)だったら(cons 'quote (cons 'x nil))という構造なので、
			// list.Cdrは(x)=(cons 'x nil)
			return evalSpecialForm(sym.Name, list.Cdr, env)
		}
	}

	//先頭要素(関数)を評価
	fn, err := Eval(first, env)
	if err != nil {
		return nil, err
	}

	//引数を評価
	args, err := evalArgs(list.Cdr, env)
	if err != nil {
		return nil, err
	}

	// 関数適用
	return apply(fn, args)
}

// 引数リストを評価
// 引数のsliceにする
func evalArgs(expr types.Expr, env *Environment) ([]types.Expr, error) {
	var args []types.Expr

	current := expr
	for {
		// nilなら終了
		// ここにnilはLispのnilであって、golangのnilでないことに注意
		if _, ok := current.(*types.Nil); ok {
			break
		}

		//Consでなければエラー
		cons, ok := current.(*types.Cons)
		if !ok {
			return nil, fmt.Errorf("invalid argument list")
		}

		//引数を評価
		arg, err := Eval(cons.Car, env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		current = cons.Cdr
	}

	return args, nil
}

// 関数を引数に適用
func apply(fn types.Expr, args []types.Expr) (types.Expr, error) {
	//まずは組み込む関数のみのサポート
	builtin, ok := fn.(BuiltinFunc)
	if !ok {
		return nil, fmt.Errorf("not a function: %v", fn)
	}

	return builtin.Call(args)
}
