package ast

type Block struct {
	LastLine int
	Stats    []Stat
	ReExps   []Exp
}
