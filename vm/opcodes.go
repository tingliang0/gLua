package vm

import "gLua/api"

// 指令编码模式
const (
	IABC  = iota // Opcode:A:B:C
	IABx         // Opcode:A:Bx
	IAsBx        // Opcode:A:sBx
	IAx          // Opcode:Ax
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
	OP_TAILCALL
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
	OpArgN = iota // argument is not used
	OpArgU        // argument is used
	OpArgR        // argument is a register or a jump offset
	OpArgK        // argument is a constant or register/constant
)

// 指令定义
type opcode struct {
	testFlag byte // operator is a test(next instruction must be a jump)
	setAFlag byte // instruction set register A
	argBMode byte // B arg mode
	argCMode byte // C arg mode
	opMode   byte // op mode
	name     string
	action   func(i Instruction, vm api.LuaVM)
}

// 所有指令
var opcodes = []opcode{
	//     T  A     B      C     mode   name       action
	opcode{0, 1, OpArgR, OpArgN, IABC, "MOVE    ", move},
	opcode{0, 1, OpArgK, OpArgN, IABx, "LOADK   ", loadK},
	opcode{0, 1, OpArgN, OpArgN, IABx, "LOADKX  ", loadKx},
	opcode{0, 1, OpArgU, OpArgU, IABC, "LOADBOOL", loadBool},
	opcode{0, 1, OpArgU, OpArgN, IABC, "LOADNIL ", loadNil},
	opcode{0, 1, OpArgU, OpArgN, IABC, "GETUPVAL", getUpval},
	opcode{0, 1, OpArgU, OpArgK, IABC, "GETTABUP", getTabUp},
	opcode{0, 1, OpArgR, OpArgK, IABC, "GETTABLE", getTable},
	opcode{0, 0, OpArgK, OpArgK, IABC, "SETTABUP", setTabUp},
	opcode{0, 0, OpArgU, OpArgN, IABC, "SETUPVAL", setUpval},
	opcode{0, 0, OpArgK, OpArgK, IABC, "SETTABLE", setTable},
	opcode{0, 1, OpArgU, OpArgU, IABC, "NEWTABLE", newTable},
	opcode{0, 1, OpArgR, OpArgK, IABC, "SELF    ", self},
	opcode{0, 1, OpArgK, OpArgK, IABC, "ADD     ", add},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SUB     ", sub},
	opcode{0, 1, OpArgK, OpArgK, IABC, "MUL     ", mul},
	opcode{0, 1, OpArgK, OpArgK, IABC, "MOD     ", mod},
	opcode{0, 1, OpArgK, OpArgK, IABC, "POW     ", pow},
	opcode{0, 1, OpArgK, OpArgK, IABC, "DIV     ", div},
	opcode{0, 1, OpArgK, OpArgK, IABC, "IDIV    ", idiv},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BAND    ", band},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BOR     ", bor},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BXOR    ", bxor},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SHL     ", shl},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SHR     ", shr},
	opcode{0, 1, OpArgR, OpArgN, IABC, "UNM     ", unm},
	opcode{0, 1, OpArgR, OpArgN, IABC, "BNOT    ", bnot},
	opcode{0, 1, OpArgR, OpArgN, IABC, "NOT     ", not},
	opcode{0, 1, OpArgR, OpArgN, IABC, "LEN     ", len},
	opcode{0, 1, OpArgR, OpArgR, IABC, "CONCAT  ", concat},
	opcode{0, 0, OpArgR, OpArgN, IAsBx, "JMP     ", jmp},
	opcode{1, 0, OpArgK, OpArgK, IABC, "EQ      ", eq},
	opcode{1, 0, OpArgK, OpArgK, IABC, "LT      ", lt},
	opcode{1, 0, OpArgK, OpArgK, IABC, "LE      ", le},
	opcode{1, 0, OpArgN, OpArgU, IABC, "TEST    ", test},
	opcode{1, 1, OpArgR, OpArgU, IABC, "TESTSET ", testSet},
	opcode{0, 1, OpArgU, OpArgU, IABC, "CALL    ", call},
	opcode{0, 1, OpArgU, OpArgU, IABC, "TAILCALL", tailCall},
	opcode{0, 0, OpArgU, OpArgN, IABC, "RETURN  ", _return},
	opcode{0, 1, OpArgR, OpArgN, IAsBx, "FORLOOP ", forLoop},
	opcode{0, 1, OpArgR, OpArgN, IAsBx, "FORPREP ", forPrep},
	opcode{0, 0, OpArgN, OpArgU, IABC, "TFORCALL", nil},
	opcode{0, 1, OpArgR, OpArgN, IAsBx, "TFORLOOP", nil},
	opcode{0, 0, OpArgU, OpArgU, IABC, "SETLIST ", setList},
	opcode{0, 1, OpArgU, OpArgN, IABx, "CLOSURE ", closure},
	opcode{0, 1, OpArgU, OpArgN, IABC, "VARARG  ", vararg},
	opcode{0, 0, OpArgU, OpArgU, IAx, "EXITAARG", nil},
}
