package state

func (self *luaState) TypeName2(idx int) string {
	return self.TypeName(self.Type(idx))
}

func (self *luaState) Len2(idx int) int64 {
	self.Len(idx)
	i, isNum := self.ToIntegerX(-1)
	if !isNum {
		self.Error2("object length is not an integer")
	}
	self.Pop(1)
	return i
}

func (self *luaState) CheckStack2(sz int, msg string) {
	if !self.CheckStack(sz) {
		if msg != "" {
			self.Error2("stack overflow (%s)", msg)
		} else {
			self.Error2("stack overflow")
		}
	}
}

func (self *luaState) Error2(fmt string, a ...interface{}) int {
	self.PushFString(fmt, a...)
	return self.Error()
}
