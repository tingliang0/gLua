package codegen

import . "gLua/binchunk"
import . "gLua/compiler/ast"

func GenProto(chunk *Block) *Prototype {
	fd := &FuncDefExp{
		IsVararg: true,
		Block:    chunk,
		LastLine: chunk.LastLine,
	}
	fi := newFuncInfo(nil, fd)
	fi.addLocVar("_ENV")
	cgFuncDefExp(fi, fd, 0)
	return toProto(fi.subFunc[0])
}
