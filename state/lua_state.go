package state

import "gLua/binchunk"

// 相当于CPU
type luaState struct {
	stack *luaStack           // register
	proto *binchunk.Prototype // code & const
	pc    int                 // pc register
}

func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
