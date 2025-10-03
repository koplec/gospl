package types

import "fmt"

type Expr interface {
	String() string
}

type Number struct {
	Value float64
}

type Symbol struct {
	Name string
}

type Nil struct{}

type String struct {
	Value string
}

type Boolean struct {
	Value bool
}

type Cons struct {
	Car Expr
	Cdr Expr
}

func (n Number) String() string {
	// 整数なら%d表示
	if n.Value == float64(int64(n.Value)) {
		return fmt.Sprintf("%d", int64(n.Value))
	}

	return fmt.Sprintf("%g", n.Value) //小数点以下があるときは少数形式、それ以外は整数形式
}

func (n Nil) String() string {
	return "NIL"
}

func (s Symbol) String() string {
	return s.Name
}

// 文字列は""をつける
func (s String) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

func (b Boolean) String() string {
	if b.Value {
		return "T"
	}
	return "NIL"
}

func (c *Cons) String() string {
	return "(...)" //いったん簡単表示
}
