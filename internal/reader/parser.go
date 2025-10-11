package reader

import (
	"fmt"
	"strconv"

	"github.com/koplec/gospl/internal/types"
)

// Parser(構文解析器)
// Lexer（字句解析器）が生成したトークン列をS式に変換する
// 入力文字列 -> [Lexer] -> トークン列 -> [Parser] -> S式(types.Expr)
// Parserは読むだけで評価はしない
// 例えば(+ 1 2)を読んでも3にならない
type Parser struct {
	lexer   *Lexer
	current Token // 現在見ているトークン
}

// Parserを生成する
// 入力文字から最初のトークンを読んだ状態にする
func NewParser(input string) *Parser {
	p := &Parser{
		lexer: NewLexer(input),
	}

	// 最初のトークンを読み込む
	p.advance()
	return p
}

// エントリーポイント, 一つの式をパースする
func (p *Parser) Parse() (types.Expr, error) {
	return p.parseExpr()
}

// 1つの式をparseして、次のトークンに進む
func (p *Parser) parseExpr() (types.Expr, error) {
	switch p.current.Type {
	case NUMBER:
		//トークンの値をfloat64に変換
		value, err := strconv.ParseFloat(p.current.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", p.current.Value)
		}
		//次のトークンへは進んでおく
		if err := p.advance(); err != nil {
			return nil, err
		}
		return types.Number{Value: value}, nil
	case STRING:
		value := p.current.Value

		//次のトークンへは進んでおく
		if err := p.advance(); err != nil {
			return nil, err
		}
		return types.String{Value: value}, nil
	case SYMBOL:
		value := p.current.Value

		//次のトークンへは進んでおく
		if err := p.advance(); err != nil {
			return nil, err
		}

		//SYMBOLトークンをExprに変換するのがこの関数の目的だから
		//BOOLEANにもここで変換が必要
		switch value {
		case "t":
			return types.Boolean{Value: true}, nil
		case "nil":
			return &types.Nil{}, nil //common lisp風にfalseじゃなくてnil
		default:
			return types.Symbol{
				Name: value,
			}, nil
		}
	case LPAREN:
		// (が来たから　)がくるまで式を読み続ける
		//そのためにparseList()を呼ぶ
		return p.parseList()
	case RPAREN:
		// ここに到達してはダメ
		return nil, fmt.Errorf("unexpected ')' at position %d:%d",
			p.current.Pos.Line, p.current.Pos.Column)
	case QUOTE:
		// 'をスキップする
		if err := p.advance(); err != nil {
			return nil, err
		}

		//次の式をパース
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		// (quote expr)の形に変換
		quote := types.Symbol{Name: "quote"} //quoteシンボルを作成

		//'expr = (quote expr) = (quote . (expr . nil))
		return &types.Cons{
			Car: quote,
			Cdr: &types.Cons{
				Car: expr,
				Cdr: &types.Nil{},
			},
		}, nil
	case EOF:
		return nil, fmt.Errorf("unexpected end of input")
	default:
		return nil, fmt.Errorf("unexpected token '%s' (type:%v) at position %d:%d",
			p.current.Value, p.current.Type, p.current.Pos.Line, p.current.Pos.Column)
	}

}

func (p *Parser) parseList() (types.Expr, error) {
	//現在のトークンは'('
	if err := p.advance(); err != nil { //(をスキップする
		return nil, err
	}

	//空リストの場合 NIL
	if p.current.Type == RPAREN {
		//')'を読み飛ばす
		if err := p.advance(); err != nil {
			return nil, err
		}
		//ほかの実装では参照を返していないのに、参照を返すのは、Lispの場合、Listは実質的に参照の一覧であるから
		return &types.Nil{}, nil
	}

	//リストの要素を読んでいく
	var car *types.Cons
	var cdr *types.Cons

	//')'でない限りループ
	for p.current.Type != RPAREN && p.current.Type != EOF {
		//１つの式をパース
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		//リストは参照！
		cons := &types.Cons{Car: expr, Cdr: &types.Nil{}}

		if car == nil { //最初の要素
			//ここでは、carもcdrも同じ構造cons=(expr, nil)を指し示す
			// car -> (１番目expr . nil)
			// cdr -> (１番目expr . nil)
			car = cons
			cdr = cons
		} else { //２番目以降
			cdr.Cdr = cons
			cdr = cons
			//ここから変わっていく
			// cdr.Cdrはひとつ前のcdrの構造を新しくする
			// 例えば１番目の場合 cdr -> (1番目expr . nil)だったものが
			// cdr -> (１番目expr . (2番目expr . nil ))になる
			// そのあと新たにcdr -> (2番目expr . nil)を指し示すようになる
			// carは常に変わらないことに注意
		}
	}

	// ')'が最後までなかったらエラー
	if p.current.Type != RPAREN {
		return nil, fmt.Errorf("expected ')', got EOF")
	}

	//')'をスキップ
	if err := p.advance(); err != nil {
		return nil, err
	}

	// carが常に先頭を指し示すから、carを返す
	return car, nil
}

// advanceで次のトークンを読み込む
func (p *Parser) advance() error {
	token, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.current = token
	return nil
}
