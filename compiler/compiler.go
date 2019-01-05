package compiler

import (
	"gLua/binchunk"
	"gLua/compiler/codegen"
	"gLua/compiler/parser"
)

func Compile(chunk, chunkName string) *binchunk.Prototype {
	ast := parser.Parse(chunk, chunkName)
	return codegen.GenProto(ast)
}
