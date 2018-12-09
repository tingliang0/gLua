package vm

// 指令编码模式
const (
	IABC = iota					// Opcode:A:B:C
	IABx						// Opcode:A:Bx
	IAsBx						// Opcode:A:sBx
	IAx							// Opcode:Ax
)

// 指令操作码
const (
	OP_MOVE = iota
	OP_LOADK
	OP_LOADKX
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETTABUP
	OP_GETTABLE
	OP_SETTABUP
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE
	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_MOD
	OP_POW
	OP_DIV
	OP_IDIV
	OP_BAND
	OP_BOR
	OP_BXOR
	OP_SHL
	OP_SHR
	OP_UNM
	OP_BNOT
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TALLCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORCALL
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSURE
	OP_VARARG
	OP_EXTRAARG
)

// 指令操作数类型
const (
	OpArgN = iota				// argument is not used
	OpArgU						// argument is used
	OpArgR						// argument is a register or a jump offset
	OpArgK						// argument is a constant or register/constant
)

// 指令定义
type opcode struct {
	testFlag byte				// operator is a test(next instruction must be a jump)
	setAFlag byte				// instruction set register A
	argBMode byte				// B arg mode
	argCMode byte				// C arg mode
	opMode byte					// op mode
	name string
}

// 所有指令
var opcodes = []opcode {
	//     T  A     B      C     mode   name
	opcode{0, 1, OpArgR, OpArgN, IABC, "MOVE    "},
	opcode{0, 1, OpArgK, OpArgN, IABx, "LOADK   "},
	opcode{0, 1, OpArgN, OpArgN, IABx, "LOADKX  "},
	opcode{0, 1, OpArgU, OpArgU, IABC, "LOADBOOL"},
	opcode{0, 1, OpArgU, OpArgN, IABC, "LOADNIL "},
	opcode{0, 1, OpArgU, OpArgN, IABC, "GETUPVAL"},
	opcode{0, 1, OpArgU, OpArgK, IABC, "GETTABUP"},
	opcode{0, 1, OpArgR, OpArgK, IABC, "GETTABLE"},
	opcode{0, 0, OpArgK, OpArgK, IABC, "SETTABUP"},
	opcode{0, 0, OpArgU, OpArgN, IABC, "SETUPVAL"},
	opcode{0, 0, OpArgK, OpArgK, IABC, "SETTABLE"},
	opcode{0, 1, OpArgU, OpArgU, IABC, "NEWTABLE"},
	opcode{0, 1, OpArgR, OpArgK, IABC, "SELF    "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "ADD     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SUB     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "MUL     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "MOD     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "POW     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "DIV     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "IDIV    "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BAND    "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BOR     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BXOR    "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SHL     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SHR     "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "UNM     "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "BNOT    "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "NOT     "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "LEN     "},
	opcode{0, 1, OpArgR, OpArgR, IABC, "CONCAT  "},
	opcode{0, 0, OpArgR, OpArgN, IAsBx,"JMP     "},
	opcode{1, 0, OpArgK, OpArgK, IABC, "EQ      "},
	opcode{1, 0, OpArgK, OpArgK, IABC, "LT      "},
	opcode{1, 0, OpArgK, OpArgK, IABC, "LE      "},
	opcode{1, 0, OpArgN, OpArgU, IABC, "TEST    "},
	opcode{1, 1, OpArgR, OpArgU, IABC, "TESTSET "},
	opcode{0, 1, OpArgU, OpArgU, IABC, "CALL    "},
	opcode{0, 1, OpArgU, OpArgU, IABC, "TALLCALL"},
	opcode{0, 0, OpArgU, OpArgN, IABC, "RETURN  "},
	opcode{0, 1, OpArgR, OpArgN, IAsBx,"FORLOOP "},
	opcode{0, 1, OpArgR, OpArgN, IAsBx,"FORPREP "},
	opcode{0, 0, OpArgN, OpArgU, IABC, "TFORCALL"},
	opcode{0, 1, OpArgR, OpArgN, IAsBx,"TFORLOOP"},
	opcode{0, 0, OpArgU, OpArgU, IABC, "SETLIST "},
	opcode{0, 1, OpArgU, OpArgN, IABx, "CLOSURE "},
	opcode{0, 1, OpArgU, OpArgN, IABC, "VARARG  "},
	opcode{0, 0, OpArgU, OpArgU, IAx,  "EXITAARG"},
}

type Instruction uint32
const MAXARG_Bx = 1 << 18 - 1
const MAXARG_sBx = MAXARG_Bx >> 1

func (i Instruction) Opcode() int {
	return int(i & 0x3F)
}

func (i Instruction) ABC() (a, b, c int) {
	a = int(i >> 6 & 0xFF)
	c = int(i >> 14 & 0x1FF)
	b = int(i >> 23 & 0x1FF)
	return
}

func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 6 & 0xFF)
	bx = int(i >> 14)
	return
}

func (i Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

func (i Instruction) Ax() int {
	return int(i >> 6)
}

func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

func (i Instruction) BMode() byte {
	return opcodes[i.Opcode()].argBMode
}

func (i Instruction) CMode() byte {
	return opcodes[i.Opcode()].argCMode
}
