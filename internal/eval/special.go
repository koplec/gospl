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
	return nil, fmt.Errorf("if not implemented yet")
}

func evalDefun(args types.Expr, evn *Environment) (types.Expr, error) {
	return nil, fmt.Errorf("defun not implemented yet")
}

func evalLambda(args types.Expr, env *Environment) (types.Expr, error) {
	return nil, fmt.Errorf("lambda not implemented yet")
}
