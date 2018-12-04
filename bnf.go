package eval

import (
	"context"
	"fmt"
)

type Evaler interface {
	Eval(ctx context.Context) BaseVal
}

// 基本元素，不可以再求值的
type BaseVal struct {
	MapValue *Intervals `  "map{" { @@ "," } @@ "}"`
	Number   *float64   `| @(Float|Int)`
	StrData  *string    `| @String `
}

func (bv BaseVal) String() string {
	s := ""
	if bv.Number != nil {
		s += fmt.Sprintf("n: %f ", *bv.Number)
	}
	if bv.StrData != nil {
		s += fmt.Sprintf("s: %s ", *bv.StrData)
	}
	if bv.MapValue != nil {
		s += fmt.Sprintf("m: %v ", *bv.MapValue)
	}
	return s
}

type Value struct {
	// 下面几个字段，之所以用指针是为了，判断是哪个字段被初始化了
	//SubExpress *Express  ` "(" @@ ")"` // 表达式递归
	BaseVal
	FuncCall *FuncCall `| @@ ` // 调用函数操作  xxx(args...)
	CtxValue *CtxValue `| @@ ` // 上下文参数   ctx.xxx
}

func (v Value) Eval(ctx context.Context) BaseVal {
	if v.Number != nil {
		return BaseVal{Number: v.Number}
	}
	if v.StrData != nil {
		return BaseVal{StrData: v.StrData}
	}
	if v.MapValue != nil {
		return BaseVal{MapValue: v.MapValue}
	}
	if v.FuncCall != nil {
		return v.FuncCall.Eval(ctx)
	}
	if v.CtxValue != nil {
		return v.CtxValue.Eval(ctx)
	}
	//if v.SubExpress != nil {
	//	return v.SubExpress.Eval(ctx)
	//}

	panicFor("Value")
	return BaseVal{}
}

type Intervals []Interval

func (Intervals Intervals) Find(val float64) string {
	for _, interval := range Intervals {
		if interval.A <= val && val < interval.B {
			return interval.Key
		}
	}
	panic(fmt.Sprintf("not support val(%f) in interval(%#v)", val, Intervals))
}

// 表示 从 [A, B) 映射到 Key 字符串， 后续会使用Key来找到新的表达式
type Interval struct {
	A   float64 `@Float`
	B   float64 `@Float`
	Key string  `@String`
}

type FuncCall struct {
	Name string  `@Ident` // 数学函数， 加个下划线前缀表示内置的
	Args []Value `"(" {@@ ","} @@ ")"`
}

// 函数调用
func (fc FuncCall) Eval(ctx context.Context) BaseVal {
	f, exist := getFunc(fc.Name)
	if !exist {
		panicFor(fmt.Sprintf("not exist func %s", fc.Name))
		return BaseVal{}
	}
	args := make([]BaseVal, 0, len(fc.Args))
	for _, arg := range fc.Args {
		args = append(args, arg.Eval(ctx))
	}
	return f(ctx, args)
}

type CtxValue struct {
	FieldName string `  "ctx." @String`
}

func (cv CtxValue) Eval(ctx context.Context) BaseVal {
	val := ctx.Value(cv.FieldName)
	if val == nil {
		panicFor(fmt.Sprintf("not exist field %s", cv.FieldName))
		return BaseVal{}
	}
	switch val := val.(type) {
	case string:
		return BaseVal{StrData: &val}
	case float64:
		return BaseVal{Number: &val}
	}
	panicFor(fmt.Sprintf("not support ctx.%s field type", cv.FieldName))
	return BaseVal{}
}

// 乘法除法优先级比较高，所以嵌入到底层
type Term struct {
	OpFactors []FactorOp ` @@*`
	Val       Value      ` @@`
}

func (t Term) Eval(ctx context.Context) BaseVal {
	if len(t.OpFactors) == 0 {
		return t.Val.Eval(ctx)
	}
	residue := *Term{OpFactors: t.OpFactors[1:], Val: t.Val}.Eval(ctx).Number

	op1 := t.OpFactors[0]
	v1 := op1.Val.Eval(ctx).Number
	switch op1.OP {
	case "*":
		r := *v1 * residue
		return BaseVal{Number: &r}
	case "/":
		r := *v1 / residue
		return BaseVal{Number: &r}
	default:
		panicFor(fmt.Sprintf("unsupport op %s", op1.OP))
		return BaseVal{}
	}
}

type FactorOp struct {
	Val Value  `@@`
	OP  string `@("*" | "/")`
}

type SimpleTerm struct {
	TermOps []TermOp ` @@* `
	Term    Term     ` @@`
}

func (st SimpleTerm) Eval(ctx context.Context) BaseVal {
	if len(st.TermOps) == 0 {
		return st.Term.Eval(ctx)
	}
	residue := *SimpleTerm{TermOps: st.TermOps[1:], Term: st.Term}.Eval(ctx).Number
	op1 := st.TermOps[0]
	v1 := *op1.Term.Eval(ctx).Number
	switch op1.OP {
	case "+":
		r := v1 + residue
		return BaseVal{Number: &r}
	case "-":
		r := v1 - residue
		return BaseVal{Number: &r}
	default:
		panicFor(fmt.Sprintf("not support op \"%s\"", op1.OP))
		return BaseVal{}
	}
}

type TermOp struct {
	Term Term   `@@`
	OP   string `@('+' | '-')`
}

type Express struct {
	SimpleTerm
}
