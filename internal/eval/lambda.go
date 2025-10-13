// typesに定義すると、Envが定義されているevalへの循環参照になるので、
// いったん、eval packageに定義
package eval

import "github.com/koplec/gospl/internal/types"

type Lambda struct {
	Params []string   //仮引数のリスト
	Body   types.Expr //関数本体はS式
	Env    *Environment
}

func (l *Lambda) String() string {
	return "#<FUNCTION>"
}
