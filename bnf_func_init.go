package eval

import (
	"sync"
	"math"
	"context"
	"fmt"
)

type EvalFunc func(ctx context.Context, args []BaseVal) BaseVal
var funcLock sync.RWMutex
var funcMap map[string]EvalFunc

func getFunc(name string) (f EvalFunc, exist bool) {
	funcLock.RLock()
	defer funcLock.RUnlock()

	f, exist = funcMap[name]
	return
}

func RegisterFunc(name string, f EvalFunc) {
	funcLock.Lock()
	defer funcLock.Unlock()

	if funcMap == nil {
		funcMap = make(map[string]EvalFunc)
	}

	funcMap[name] = f
}

func init() {

	RegisterFunc("add", add)
	RegisterFunc("minus", minus)
	RegisterFunc("mul", mul)
	RegisterFunc("div", div)
	RegisterFunc("floor", floor)
	RegisterFunc("ceil", ceil)
	RegisterFunc("pow", pow)
	RegisterFunc("max", max)
	RegisterFunc("min", min)
	RegisterFunc("gt", gt)
	RegisterFunc("gte", gte)
	RegisterFunc("lt", lt)
	RegisterFunc("lte", lte)
	RegisterFunc("eq", eq)

	RegisterFunc("findExp", findExpFunc)
	RegisterFunc("mapKeyVal", mapKeyVal)

}

func add(ctx context.Context, args []BaseVal) BaseVal {
	v1 := *args[0].Number
	v2 := *args[1].Number
	result := v1 * v2
	return BaseVal{Number:&result}
}

func minus(ctx context.Context, args []BaseVal) BaseVal {
	v1 := *args[0].Number
	v2 := *args[1].Number
	result := v1 - v2
	return BaseVal{Number:&result}
}

func mul(ctx context.Context, args []BaseVal) BaseVal {
	v1 := *args[0].Number
	v2 := *args[1].Number
	result := v1 * v2
	return BaseVal{Number:&result}
}

func div(ctx context.Context, args []BaseVal) BaseVal {
	v1 := *args[0].Number
	v2 := *args[1].Number
	if v2 == 0 {
		panicFor("除数不能为0")
	}
	result := v1 / v2
	return BaseVal{Number:&result}
}

// floor(a)  	地板除
func floor(ctx context.Context, args []BaseVal) BaseVal {
	v := *args[0].Number
	result := math.Floor(v)
	return BaseVal{Number:&result}
}

// ceil(a)		向上取整
func ceil(ctx context.Context, args []BaseVal) BaseVal {
	v := *args[0].Number
	result := math.Ceil(v)
	return BaseVal{Number:&result}
}

// pow(a,b)		幂次
func pow(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[1].Number
	result := math.Pow(a, b)
	return BaseVal{Number:&result}
}

// max(a,b)		取最大的
func max(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[1].Number
	result := math.Max(a, b)
	return BaseVal{Number:&result}
}

// min(a,b)		取最小的
func min(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[0].Number
	result := math.Min(a, b)
	return BaseVal{Number:&result}
}

// gt(a,b)		如果a 大于 b 得 1 否则 得 0
func gt(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[1].Number
	var result float64
	if a > b {
		result = 1
	} else {
		result = 0
	}
	return BaseVal{Number:&result}
}

// gte(a,b)		如果a 大于等于 b 得 1 否则 得 0
func gte(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[1].Number
	var result float64
	if a >= b {
		result = 1
	} else {
		result = 0
	}
	return BaseVal{Number:&result}
}

// lt(a,b)		如果a 小于 b 得 1 否则 得 0
func lt(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[1].Number
	var result float64
	if a < b {
		result = 1
	} else {
		result = 0
	}

	return BaseVal{Number:&result}
}

// lte(a,b)		如果a 小于等于 b 得 1 否则 得 0
func lte(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[1].Number
	var result float64
	if a <= b {
		result = 1
	} else {
		result = 0
	}
	return BaseVal{Number:&result}
}

// eq(a,b)		如果 a 等于 b 得 1 否则 得 0
func eq(ctx context.Context, args []BaseVal) BaseVal {
	a := *args[0].Number
	b := *args[1].Number
	var result float64
	if a == b {
		result = 1
	} else {
		result = 0
	}

	return BaseVal{Number:&result}
}

/* ===========  */

func findExpFunc(ctx context.Context, args []BaseVal) BaseVal {
	expName := *args[0].StrData
	v, exist := findExp(expName)
	if exist {
		panicFor(fmt.Sprintf("not exist exp for key %s", expName))
		return BaseVal{}
	}

	return v.Eval(ctx)
}

// 映射某个值到某个key
func mapKeyVal(ctx context.Context, args []BaseVal) BaseVal {
	v := *args[0].Number
	intervals := *args[0].MapValue
	mapkey := intervals.Find(v)
	return BaseVal{StrData:&mapkey}
}