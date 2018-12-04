目的: golang 中支持 数学公式的求值

```
// BNF:
// exp ::= math | mathFunc | func | val
// math ::= exp ( "+" | "-" | "*" | "/" ) exp
// mathFunc ::= _floor(a) | _ceil(a) | _pow(a, b) | _max(a, b) | _min(a,b) | _gt(a, b) | _gte(a,b) | _lt(a,b) | _lte(a,b) | _eq(a,b)         # a, b 属于 float
// func ::= findEval(key) | mapValKey(a, intervals)       # a 属于float
// val ::= string | float | intervals
// intervals ::=  "map{"  interval* "}"
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
// 程序运行中的参数：
//          ctx.xxxx   // xxx --> level , weaponMultiple
```