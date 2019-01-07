package stdlib

import (
	. "gLua/api"
	"gLua/number"
	"math"
	"math/rand"
)

var mathLib = map[string]GoFunction{
	"random":     mathRandom,
	"randomseed": mathRandomSeed,
	"max":        mathMax,
	"min":        mathMin,
	"exp":        mathExp,
	"log":        mathLog,
	"deg":        mathDeg,
	"rad":        mathRad,
	"sin":        mathSin,
	"cos":        mathCos,
	"tan":        mathTan,
	"asin":       mathAsin,
	"acos":       mathAcos,
	"atan":       mathAtan,
	"ceil":       mathCeil,
	"floor":      mathFloor,
	"abs":        mathAbs,
	"sqrt":       mathSqrt,
	"fmod":       mathFmod,
	"ult":        mathUlt,
	"tointeger":  mathToInt,
	"type":       mathType,
}

func mathRandom(ls LuaState) int {
	var low, up int64
	switch ls.GetTop() {
	case 0:
		ls.PushNumber(rand.Float64())
		return 1
	case 1:
		low = 1
		up = ls.CheckInteger(1)
	case 2:
		low = ls.CheckInteger(1)
		up = ls.CheckInteger(2)
	default:
		return ls.Error2("wrong number of arguments")
	}

	ls.ArgCheck(low <= up, 1, "interval is empty")
	ls.ArgCheck(low >= 0 || up <= math.MaxInt64+low, 1, "interval too large")
	if up-low == math.MaxInt64 {
		ls.PushInteger(low + rand.Int63())
	} else {
		ls.PushInteger(low + rand.Int63n(up-low+1))
	}
	return 1
}

func mathRandomSeed(ls LuaState) int {
	x := ls.CheckNumber(1)
	rand.Seed(int64(x))
	return 0
}

func mathMax(ls LuaState) int {
	n := ls.GetTop()
	imax := 1
	ls.ArgCheck(n >= 1, 1, "value expected")
	for i := 2; i <= n; i++ {
		if ls.Compare(imax, i, LUA_OPLT) {
			imax = i
		}
	}
	ls.PushValue(imax)
	return 1
}

func mathMin(ls LuaState) int {
	n := ls.GetTop()
	imin := 1
	ls.ArgCheck(n >= 1, 1, "value expected")
	for i := 2; i <= n; i++ {
		if ls.Compare(i, imin, LUA_OPLT) {
			imin = i
		}
	}
	ls.PushValue(imin)
	return 1
}

func mathExp(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(math.Exp(x))
	return 1
}

func mathLog(ls LuaState) int {
	x := ls.CheckNumber(1)
	var res float64

	if ls.IsNoneOrNil(2) {
		res = math.Log(x)
	} else {
		base := ls.ToNumber(2)
		if base == 2 {
			res = math.Log2(x)
		} else if base == 10 {
			res = math.Log10(x)
		} else {
			res = math.Log(x) / math.Log(base)
		}
	}

	ls.PushNumber(res)
	return 1
}

func mathDeg(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(x * 180 / math.Pi)
	return 1
}

func mathRad(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(x * math.Pi / 180)
	return 1
}

func mathSin(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(math.Sin(x))
	return 1
}

func mathCos(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(math.Cos(x))
	return 1
}

func mathTan(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(math.Tan(x))
	return 1
}

func mathAsin(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(math.Asin(x))
	return 1
}

func mathAcos(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(math.Acos(x))
	return 1
}

func mathAtan(ls LuaState) int {
	y := ls.CheckNumber(1)
	x := ls.OptNumber(2, 1.0)
	ls.PushNumber(math.Atan2(y, x))
	return 1
}

func mathCeil(ls LuaState) int {
	if ls.IsInteger(1) {
		ls.SetTop(1) /* integer is its own ceil */
	} else {
		x := ls.CheckNumber(1)
		_pushNumInt(ls, math.Ceil(x))
	}
	return 1
}

func mathFloor(ls LuaState) int {
	if ls.IsInteger(1) {
		ls.SetTop(1) /* integer is its own floor */
	} else {
		x := ls.CheckNumber(1)
		_pushNumInt(ls, math.Floor(x))
	}
	return 1
}

func mathFmod(ls LuaState) int {
	if ls.IsInteger(1) && ls.IsInteger(2) {
		d := ls.ToInteger(2)
		if uint64(d)+1 <= 1 { /* special cases: -1 or 0 */
			ls.ArgCheck(d != 0, 2, "zero")
			ls.PushInteger(0) /* avoid overflow with 0x80000... / -1 */
		} else {
			ls.PushInteger(ls.ToInteger(1) % d)
		}
	} else {
		x := ls.CheckNumber(1)
		y := ls.CheckNumber(2)
		ls.PushNumber(x - math.Trunc(x/y)*y)
	}

	return 1
}

func mathModf(ls LuaState) int {
	if ls.IsInteger(1) {
		ls.SetTop(1)     /* number is its own integer part */
		ls.PushNumber(0) /* no fractional part */
	} else {
		x := ls.CheckNumber(1)
		i, f := math.Modf(x)
		_pushNumInt(ls, i)
		if math.IsInf(x, 0) {
			ls.PushNumber(0)
		} else {
			ls.PushNumber(f)
		}
	}

	return 2
}

func mathAbs(ls LuaState) int {
	if ls.IsInteger(1) {
		x := ls.ToInteger(1)
		if x < 0 {
			ls.PushInteger(-x)
		}
	} else {
		x := ls.CheckNumber(1)
		ls.PushNumber(math.Abs(x))
	}
	return 1
}

func mathSqrt(ls LuaState) int {
	x := ls.CheckNumber(1)
	ls.PushNumber(math.Sqrt(x))
	return 1
}

func mathUlt(ls LuaState) int {
	m := ls.CheckInteger(1)
	n := ls.CheckInteger(2)
	ls.PushBoolean(uint64(m) < uint64(n))
	return 1
}

func mathToInt(ls LuaState) int {
	if i, ok := ls.ToIntegerX(1); ok {
		ls.PushInteger(i)
	} else {
		ls.CheckAny(1)
		ls.PushNil() /* value is not convertible to integer */
	}
	return 1
}

func mathType(ls LuaState) int {
	if ls.Type(1) == LUA_TNUMBER {
		if ls.IsInteger(1) {
			ls.PushString("integer")
		} else {
			ls.PushString("float")
		}
	} else {
		ls.CheckAny(1)
		ls.PushNil()
	}
	return 1
}

func _pushNumInt(ls LuaState, d float64) {
	if i, ok := number.FloatToInteger(d); ok { /* does 'd' fit in an integer? */
		ls.PushInteger(i) /* result is integer */
	} else {
		ls.PushNumber(d) /* result is float */
	}
}

func OpenMathLib(ls LuaState) int {
	ls.NewLib(mathLib)
	ls.PushNumber(math.Pi)
	ls.SetField(-2, "pi")
	ls.PushNumber(math.Inf(1))
	ls.SetField(-2, "huge")
	ls.PushInteger(math.MaxInt64)
	ls.SetField(-2, "maxinteger")
	ls.PushInteger(math.MinInt64)
	ls.SetField(-2, "mininteger")
	return 1
}
