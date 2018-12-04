package eval

import (
	"testing"
	"context"
	"fmt"
)

func TestExp_Eval(t *testing.T) {

	// + - * /

	// _floor(a)  	地板除
	// _ceil(a)		向上取整
	// _pow(a,b)		幂次
	// _max(a,b)		取最大的
	// _min(a,b)		取最小的
	// _gt(a,b)		如果a 大于 b 得 1 否则 得 0
	// _gte(a,b)		如果a 大于等于 b 得 1 否则 得 0
	// _lt(a,b)		如果a 小于 b 得 1 否则 得 0
	// _lte(a,b)		如果a 小于等于 b 得 1 否则 得 0
	// _eq(a,b)		如果 a 等于 b 得 1 否则 得 0

	//
	data := []byte(`
		{ 
			"exp11": { "index": "exp1", "exp": "add(3,5)" },
			"exp12": { "index": "exp1", "exp": "minus(3, 5)" },
			"exp13": { "index": "exp1", "exp": "minus(5, 3)" },
			"exp14": { "index": "exp1", "exp": "add(div(3, 5), 4)" },

"exp21": { "index": "exp1", "exp": "floor(3.8)" },
"exp22": { "index": "exp1", "exp": "ceil(3.8)" },
"exp23": { "index": "exp1", "exp": "pow(2, 8)" },
"exp24": { "index": "exp1", "exp": "max(3, 8)" },
"exp25": { "index": "exp1", "exp": "min(3, 8)" },
"exp26": { "index": "exp1", "exp": "gt(3, 8)" },
"exp27": { "index": "exp1", "exp": "gte(3, 8)" },
"exp28": { "index": "exp1", "exp": "lt(3, 8)" },
"exp29": { "index": "exp1", "exp": "lte(3, 8)" },
"exp210": { "index": "exp1", "exp": "eq(3, 8)" },

"exp31": { "index": "exp1", "exp": "findExp('exp11')" },
"exp32": { "index": "exp1", "exp": "mapKeyVal( min(3.8, 10), map{ 3 8 'exp25', 8 200 'exp12' })" },

"exp41": { "index": "exp1", "exp": "min(ctx.level, ctx.pearls)" },
"exp42": { "index": "exp1", "exp": "pow(ctx.level, findEval(exp32)" }
		}
	`)
	//data = []byte(`
	//	{
	//		"exp11": { "index": "exp1", "exp": "(2 + 1)" }
	//	}
	//`)
	err := RegisterEvalExpress(data)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}
	fmt.Println("register eval finished!")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "level", 7)
	ctx = context.WithValue(ctx, "pearls", 3)
	vals := []BaseVal{
		//findExpMust("'exp11'").Eval(ctx),
		//findEval("exp12").Eval(ctx),
		//findEval("exp13").Eval(ctx),
		//findEval("exp14").Eval(ctx),
		//
		//findEval("exp21").Eval(ctx),
		//findEval("exp22").Eval(ctx),
		//findEval("exp23").Eval(ctx),
		//findEval("exp24").Eval(ctx),
		//findEval("exp25").Eval(ctx),
		//findEval("exp26").Eval(ctx),
		//findEval("exp27").Eval(ctx),
		//findEval("exp28").Eval(ctx),
		//findEval("exp29").Eval(ctx),
		//findEval("exp210").Eval(ctx),
		//
		//findEval("exp31").Eval(ctx),
		//findEval("exp32").Eval(ctx),
		//
		//findEval("exp41").Eval(ctx),
		//findEval("exp42").Eval(ctx),
	}

	fmt.Println("show:")
	fmt.Printf("%s\n", vals)
}
