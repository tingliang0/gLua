package ast

type Stat interface {
}

// ;
type Empty struct {
}

// break
type BreakStat struct {
	Line int
}

// ::
type LabelStat struct {
	Name string
}

// goto
type GotoStat struct {
	Name string
}

// do
type DoStat struct {
	Block *Block
}

// functioncall
type FuncCallStat = FuncCallExp

// while
type WhileStat struct {
	Exp   Exp
	Block *Block
}

// repeat
type RepeatStat struct {
	Block *Block
	Exp   Exp
}

// if
type IfStat struct {
	Exps   []Exp
	Blocks []*Block
}

// for
type ForNumStat struct {
	LineOfFor int
	LineOfDo  int
	VarName   string
	InitExp   Exp
	LimitExp  Exp
	StepExp   Exp
	Block     *Block
}

// for in
type ForInStat struct {
	LineOfDo int
	NameList []string
	ExpList  []Exp
	Block    *Block
}

// local
type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}

// var
type AssignStat struct {
	LastLine int
	VarList  []Exp
	ExpList  []Exp
}

// local function
type LocalFuncDefStat struct {
	Name string
	Exp  *FuncDefExp
}
