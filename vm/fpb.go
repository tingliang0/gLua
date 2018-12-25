package vm

func Int2fb(x int) int {
	e := 0
	if x < 8 {
		return x
	}
	for x >= (8 << 4) {
		x = (x + 0xf) >> 4		// x = ceil(x / 16)
		e += 4
	}
	for x >= (8 << 1) {
		x = (x + 1) >> 1		// x = ceil(x / 2)
		e++
	}
	return ((e + 1) << 3) | (x - 8)
}

func Fb2int(x int) int {
	if x < 8 {
		return x
	} else {
		return ((x & 7) + 8) << uint((x>>3) - 1)
	}
}
