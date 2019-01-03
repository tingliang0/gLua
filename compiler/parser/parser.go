package parser

import (
	. "gLua/compiler/ast"
	. "gLua/compiler/lexer"
)

func Parse(chunk, chunkName string) *Block {
	lexer := NewLexer(chunk, chunkName)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_EOF)
	return block
}
