// スペシャルフォーム
package eval

import (
	"fmt"

	"github.com/koplec/gospl/internal/types"
)

const (
	SpecialFormQuote  = "quote"
	SpecialFormIf     = "if"
	SpecialFormLambda = "lambda"
	SpecialFormDefun  = "defun"
)

func isSpecialForm(name string) bool {
	switch name {
	case SpecialFormDefun, SpecialFormIf, SpecialFormLambda, SpecialFormQuote:
		return true
	default:
		return false
	}
}

// 特殊形式の評価
func evalSpecialForm(name string, args types.Expr, env *Environment) (types.Expr, error) {
	switch name {
	case SpecialFormQuote:
		return evalQuote(args)
	case SpecialFormDefun:
		return evalDefun(args, env)
	case SpecialFormLambda:
		return evalLambda(args, env)
	case SpecialFormIf:
		return evalIf(args, env)
	default:
		return nil, fmt.Errorf("unknown special form:%s", name)
	}
}

// 引数を評価せずにそのまま返す
// argsは常に引数のリストであって、引数そのものでないことに注意
func evalQuote(args types.Expr) (types.Expr, error) {
	//quoteは引数を一つだけとる
	cons, ok := args.(*types.Cons)
	if !ok {
		return nil, fmt.Errorf("quote requires exactly 1 argument")
	}

	//引数が1つだけか確認 cdrがNILなら引数は一つ
	//例えば、(quote x)の場合、argsは(x) というリストなので、cons.CdrはNIL
	//例えば、(quote (a b c))の場合は、argsは((a b c))なので、
	//(cons   (cons a (cons b (cons c nil)))  nil)なので、cons.CdrはNILになる
	if _, ok := cons.Cdr.(*types.Nil); !ok {
		return nil, fmt.Errorf("quote requires exactly 1 argument")
	}

	//評価せずに返す
	return cons.Car, nil
}

func evalIf(args types.Expr, env *Environment) (types.Expr, error) {
	//(if condition then-expr else-expr )
	/// else-exprは省略可能

	cons, ok := args.(*types.Cons)
	if !ok {
		return nil, fmt.Errorf("if requires at least 2 arguments")
	}

	// 条件式
	condition := cons.Car

	// then, else
	restCons, ok := cons.Cdr.(*types.Cons)
	if !ok {
		return nil, fmt.Errorf("if requires at least 2 arguments")
	}
	thenExpr := restCons.Car

	var elseExpr types.Expr = &types.Nil{}
	if elseCons, ok := restCons.Cdr.(*types.Cons); ok {
		elseExpr = elseCons.Car

		// 引数が３つより多い場合はエラー
		if _, ok := elseCons.Cdr.(*types.Nil); !ok {
			return nil, fmt.Errorf("if requires at most 3 arguments")
		}
	} else if _, ok := restCons.Cdr.(*types.Nil); !ok {
		//例えば
		// else節がない場合、引数２つ(if t 'aa)みたいなとき
		// これは(cons 'if (cons t (cons 'aa nil)))
		// このときargsに渡されるのは (const t (cons aa 'nil))
		// すなわちrestConsは、(cons 'aa nil))
		// するとrestCons.Cdrは、NILになる
		//では、
		// (if t 1 . 2)の場合、不正であるとき
		// (cons 'if (cons t (cons 1 2)))なので
		// argsは、(cons t (cons 1 2))
		// restConsは(cons 1 2)
		// するとrestCons.CdrはNILにならない
		return nil, fmt.Errorf("if: invalid argument list")
	}

	// 条件式の評価
	condResult, err := Eval(condition, env)
	if err != nil {
		return nil, err
	}

	//条件式の真偽判定
	// NILとBoolean{Value:false}に注意
	isFalse := false
	if _, ok := condResult.(*types.Nil); ok {
		isFalse = true
	} else if b, ok := condResult.(types.Boolean); ok && !b.Value {
		isFalse = true
	}

	if !isFalse { //すなわち if true
		return Eval(thenExpr, env)
	} else {
		return Eval(elseExpr, env)
	}
}

func evalDefun(args types.Expr, evn *Environment) (types.Expr, error) {
	return nil, fmt.Errorf("defun not implemented yet")
}

func evalLambda(args types.Expr, env *Environment) (types.Expr, error) {
	return nil, fmt.Errorf("lambda not implemented yet")
}
