package state

import (
	"fmt"
	. "gLua/api"
	"gLua/stdlib"
	"io/ioutil"
)

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

func (self *luaState) ToString2(idx int) string {
	if self.CallMeta(idx, "__tostring") {
		if !self.IsString(-1) {
			self.Error2("'__tostring' must return a string")
		}
	} else {
		switch self.Type(idx) {
		case LUA_TNUMBER:
			if self.IsInteger(idx) {
				self.PushString(fmt.Sprintf("%d", self.ToInteger(idx)))
			} else {
				self.PushString(fmt.Sprintf("%g", self.ToNumber(idx)))
			}
		case LUA_TSTRING:
			self.PushValue(idx)
		case LUA_TBOOLEAN:
			if self.ToBoolean(idx) {
				self.PushString("true")
			} else {
				self.PushString("false")
			}
		case LUA_TNIL:
			self.PushString("nil")
		default:
			tt := self.GetMetafield(idx, "__name")
			var kind string
			if tt == LUA_TSTRING {
				kind = self.CheckString(-1)
			} else {
				kind = self.TypeName2(idx)
			}
			self.PushString(fmt.Sprintf("%s: %p", kind, self.ToPointer(idx)))
			if tt != LUA_TNIL {
				self.Remove(-2)
			}
		}
	}
	return self.CheckString(-1)
}

func (self *luaState) LoadString(s string) int {
	return self.Load([]byte(s), s, "bt")
}

func (self *luaState) LoadFileX(filename, mode string) int {
	if data, err := ioutil.ReadFile(filename); err == nil {
		return self.Load(data, "@"+filename, mode)
	}
	return LUA_ERRFILE
}

func (self *luaState) LoadFile(filename string) int {
	return self.LoadFileX(filename, "bt")
}

func (self *luaState) DoFile(filename string) bool {
	return self.LoadFile(filename) == LUA_OK &&
		self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

func (self *luaState) DoString(str string) bool {
	return self.LoadString(str) == LUA_OK && self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

func (self *luaState) ArgError(arg int, extraMsg string) int {
	return self.Error2("bad argument #%d (%s)", arg, extraMsg)
}

func (self *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		self.ArgError(arg, extraMsg)
	}
}

func (self *luaState) CheckAny(arg int) {
	if self.Type(arg) == LUA_TNONE {
		self.ArgError(arg, "value expected")
	}
}

func (self *luaState) CheckType(arg int, t LuaType) {
	if self.Type(arg) != t {
		self.tagError(arg, t)
	}
}

func (self *luaState) CheckInteger(arg int) int64 {
	i, ok := self.ToIntegerX(arg)
	if !ok {
		self.intError(arg)
	}
	return i
}

func (self *luaState) CheckNumber(arg int) float64 {
	f, ok := self.ToNumberX(arg)
	if !ok {
		self.tagError(arg, LUA_TNUMBER)
	}
	return f
}

func (self *luaState) CheckString(arg int) string {
	s, ok := self.ToStringX(arg)
	if !ok {
		self.tagError(arg, LUA_TSTRING)
	}
	return s
}

func (self *luaState) OptInteger(arg int, def int64) int64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckInteger(arg)
}

func (self *luaState) OptNumber(arg int, def float64) float64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckNumber(arg)
}

func (self *luaState) OptString(arg int, def string) string {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckString(arg)
}

func (self *luaState) tagError(arg int, tag LuaType) {
	self.typeError(arg, self.TypeName(LuaType(tag)))
}

func (self *luaState) OpenLibs() {
	libs := map[string]GoFunction{
		"_G":        stdlib.OpenBaseLib,
		"math":      stdlib.OpenMathLib,
		"table":     stdlib.OpenTableLib,
		"string":    stdlib.OpenStringLib,
		"utf8":      stdlib.OpenUTF8Lib,
		"os":        stdlib.OpenOSLib,
		"package":   stdlib.OpenPackageLib,
		"coroutine": stdlib.OpenCoroutineLib,
	}

	for name, fun := range libs {
		self.RequireF(name, fun, true)
		self.Pop(1)
	}
}

func (self *luaState) RequireF(modname string, openf GoFunction, glb bool) {
	self.GetSubTable(LUA_REGISTRYINDEX, "_LOADED")
	self.GetField(-1, modname)
	if !self.ToBoolean(-1) {
		self.Pop(1)
		self.PushGoFunction(openf)
		self.PushString(modname)
		self.Call(1, 1)
		self.PushValue(-1)
		self.SetField(-3, modname)
	}
	self.Remove(-2)
	if glb {
		self.PushValue(-1)
		self.SetGlobal(modname)
	}
}

func (self *luaState) GetSubTable(idx int, fname string) bool {
	if self.GetField(idx, fname) == LUA_TTABLE {
		return true
	}
	self.Pop(1)
	idx = self.stack.absIndex(idx)
	self.NewTable()
	self.PushValue(-1)
	self.SetField(idx, fname)
	return false
}

func (self *luaState) typeError(arg int, tname string) int {
	var typeArg string /* name for the type of the actual argument */
	if self.GetMetafield(arg, "__name") == LUA_TSTRING {
		typeArg = self.ToString(-1) /* use the given type name */
	} else if self.Type(arg) == LUA_TLIGHTUSERDATA {
		typeArg = "light userdata" /* special name for messages */
	} else {
		typeArg = self.TypeName2(arg) /* standard name */
	}
	msg := tname + " expected, got " + typeArg
	self.PushString(msg)
	return self.ArgError(arg, msg)
}

func (self *luaState) GetMetafield(obj int, event string) LuaType {
	if !self.GetMetatable(obj) { /* no metatable? */
		return LUA_TNIL
	}

	self.PushString(event)
	tt := self.RawGet(-2)
	if tt == LUA_TNIL { /* is metafield nil? */
		self.Pop(2) /* remove metatable and metafield */
	} else {
		self.Remove(-2) /* remove only metatable */
	}
	return tt /* return metafield type */
}

func (self *luaState) intError(arg int) {
	if self.IsNumber(arg) {
		self.ArgError(arg, "number has no integer representation")
	} else {
		self.tagError(arg, LUA_TNUMBER)
	}
}

// [-0, +(0|1), e]
// http://www.lua.org/manual/5.3/manual.html#luaL_callmeta
func (self *luaState) CallMeta(obj int, event string) bool {
	obj = self.AbsIndex(obj)
	if self.GetMetafield(obj, event) == LUA_TNIL { /* no metafield? */
		return false
	}

	self.PushValue(obj)
	self.Call(1, 1)
	return true
}

func (self *luaState) NewLib(l FuncReg) {
	self.NewLibTable(l)
	self.SetFuncs(l, 0)
}

func (self *luaState) NewLibTable(l FuncReg) {
	self.CreateTable(0, len(l))
}

// [-nup, +0, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_setfuncs
func (self *luaState) SetFuncs(l FuncReg, nup int) {
	self.CheckStack2(nup, "too many upvalues")
	for name, fun := range l { /* fill the table with given functions */
		for i := 0; i < nup; i++ { /* copy upvalues to the top */
			self.PushValue(-nup)
		}
		// r[-(nup+2)][name]=fun
		self.PushGoClosure(fun, nup) /* closure with those upvalues */
		self.SetField(-(nup + 2), name)
	}
	self.Pop(nup) /* remove upvalues */
}
