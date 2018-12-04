package eval

import (
	"context"
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
	"log"
)


type NonExistEvaler struct {
	Name string
}

func (n NonExistEvaler) Eval(ctx context.Context) BaseVal {
	return panicFor(n.Name)
}

func panicFor(domain string) BaseVal {
	panic(fmt.Sprintf("[%s] shouldn't run to here", domain))
	return BaseVal{}
}

// BNF:
// exp ::= math | mathFunc | func | val
// math ::= exp ( "+" | "-" | "*" | "/" ) exp
// func ::= findEval(key) | mapValKey(a, intervals)       # a 属于float
// val ::= string | float | intervals
// intervals ::=  "map{"  interval { "," interval }* "}"
// interval ::=  float float key   // a b key [a, b) --> key
// key ::= string

// 注册表达式
// 支持的数学函数:
// floor(a)  	地板除
// ceil(a)		向上取整
// pow(a,b)		幂次
// max(a,b)		取最大的
// min(a,b)		取最小的
// gt(a,b)		如果a 大于 b 得 1 否则 得 0
// gte(a,b)		如果a 大于等于 b 得 1 否则 得 0
// lt(a,b)		如果a 小于 b 得 1 否则 得 0
// lte(a,b)		如果a 小于等于 b 得 1 否则 得 0
// eq(a,b)		如果 a 等于 b 得 1 否则 得 0
// 支持的再求值函数：
// 		findEval(keyName)  --> Exp	// 这个主要是结合下面这个 mapValKey， 用来，算出某个值后， 映射到某个新的表达式进行再次求值
//		mapValKey(val, intervals) --> keyName  // 这个主要是用来， 将数值映射到某个key上， 譬如， 有个规则， [3,8) --> "level1Exp",  [8,200) --> "level2Exp"
//
// data:  {
// 		"key": { "index": "xxxx", "exp": `1+4+_max(4,2) + _min(findEval('key1'), findEval(mapValKey(3, map{ 3 8 "level1", 8 200 "level2" })))`},
//      "level1": { "index": "xxxxx", "exp": `3` },
//		"level2": { "index": "nnnn", "exp": `pow(3, 9)` }
//		}

var parser *participle.Parser
func init() {
	var err error
	grammar := `
		Interval = Number Number String .
		MapValue = "map{" { Interval "," }  Interval "}" .
		Express = (String | Number | FuncCall | Express ) .
		math = Express [ "+" | "-" | "*" | "/" ] Express .
		FuncCall = Ident "(" { Express "," } [ Express ] ")" .
		Ident = (alpha | "_") { "_" | alpha | digit } .
		String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
		Number = [ "-" | "+" ] ("." | digit) { "." | digit } .
		EOL = ( "\n" | "\r" ) { "\n" | "\r" }.
		Whitespace = ( " " | "\t" ) { " " | "\t" } .

		alpha = "a"…"z" | "A"…"Z" .
		digit = "0"…"9" .
		any = "\u0000"…"\uffff" .
	`
	log.Println(grammar)
	basicLexer := lexer.Must(ebnf.New(grammar))
	parser, err = participle.Build(&Value{}, participle.Lexer(basicLexer))
	if err != nil {
		panic(fmt.Sprintf("build exp err: %s\n", err))
	}
	fmt.Println("parser build success")
}