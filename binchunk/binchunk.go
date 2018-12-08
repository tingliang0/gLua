package binchunk

const (
	LUA_SIGNATURE = "\x1bLua"
	LUAC_VERSION = 0x53
	LUAC_FORMAT = 0
	LUAC_DATA = "\x19\x93\r\n\x1a\n"
	CINT_SIZE = 4
	CSZIET_SIZE = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE = 8
	LUAC_INT = 0x5678
	LUAC_NUM = 370.5
)

// 常量类型
const (
	TAG_NIL = 0x00
	TAG_BOOLEAN = 0x01
	TAG_NUMBER = 0x03
	TAG_INTEGER = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR = 0x14
)

// Upvalue
type Upvalue struct {
	Instack byte
	Idx byte
}

// 局部变量
type LocVar struct {
	VarName string
	StartPC uint32
	EndPC uint32
}

// 二进制chunk头部
type header struct {
	signature [4]byte			// 签名, 1B4C7561
	version byte				// 版本号, 53
	format byte					// 格式号, 00
	luacData [6]byte			// LUAC_DATA, 19930D0A1A0A
	cintSize byte				// cint类型字节数，04
	sizetSize byte				// size_t类型字节数，08
	instructionSize byte		// 虚拟机指令占用字节数，04
	luaIntegerSize byte			// 整数字节数，08
	luaNumberSize byte			// 浮点数字节数，08
	luacInt int64				// 8字节保存0x5678，用于检测大小端
	luacNum float64				// 8字节保存浮点数370.5，用于检测浮点数格式
}

// 函数原型
type Prototype struct {
	Source string				// 源文件名
	LineDefined uint32			// 函数在源文件开始行号
	LastLineDefined uint32		// 函数在源文件结束行号
	NumParams byte				// 固定参数个数
	IsVararg byte				// 是否Vararg函数，0表示否，1表示是，主函数是Vararg函数
	MaxStackSize byte			// 运行函数时需要用到虚拟寄存器数量
	Code []uint32				// 指令表，每个指令占4个字节
	Constants []interface{}		// 常量表
	Upvalues []Upvalue			// Upvalue表
	Protos []*Prototype			// 子函数原型表
	LineInfo []uint32			// 行号表，行号表中的行号和指令表中的指令一一对应，记录每条指令在源代码中对应的行号（调试用）
	LocVars []LocVar			// 局部变量表（调试用）
	UpvalueNames []string		// Upvalue名，和Upvalues表一一对应。（调试用）
}

// binaryChunk文件结构
type binaryChunk struct {
	header						// 头部
	sizeUpvalues byte			// 主函数upvalue数量
	mainFunc *Prototype			// 主函数原型
}

// []byte -> prototype
func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()		// 检测头部
	reader.readByte()			// 跳过upvalue数量
	return reader.readProto("")
}
