package ast

type Exp interface {
}

// nil
type NilExp struct {
	Line int
}

// true
type TrueExp struct {
	Line int
}

// false
type FalseExp struct {
	Line int
}

// ...
type VarargExp struct {
	Line int
}

// int
type IntegerExp struct {
	Line int
	Val  int64
}

// float
type FloatExp struct {
	Line int
	Val  int64
}

// string
type StringExp struct {
	Line int
	Str  string
}

// name
type NameExp struct {
	Line int
	Name string
}

type UnopExp struct {
	Line int // line of operator
	Op   int // operator
	Exp  Exp
}

type BinopExp struct {
	Line int
	Op   int
	Exp1 Exp
	Exp2 Exp
}

type ConcatExp struct {
	Line int
	Exps []Exp
}

// table
type TableConstructorExp struct {
	Line     int // line of '{'
	LastLine int // line of '}'
	KeyExps  []Exp
	ValExps  []Exp
}

type FuncDefExp struct {
	Line     int
	LastLine int
	ParList  []string
	IsVararg bool
	Block    *Block
}

type ParensExp struct {
	Exp Exp
}

type TableAccessExp struct {
	LastLine  int // line of ']'
	PreFixExp Exp
	KeyExp    Exp
}

type FuncCallExp struct {
	Line      int // line of '('
	LastLine  int // line of ')'
	PrefixExp Exp
	NameExp   *StringExp
	Args      []Exp
}
