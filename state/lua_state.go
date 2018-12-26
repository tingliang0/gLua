package state

// 相当于CPU
type luaState struct {
	stack *luaStack // register
}

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}

func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
}

func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
}
