package libgen

//big ending
func Uint16ToArrBE(v uint16) []uint8 {
	return []uint8{uint8(v>>8) & 0xFF, uint8(v) & 0xFF}
}

//big ending
func ArrToUint16BE(arr []uint8) uint16 {
	if len(arr) != 2 {
		return 0
	}
	return (uint16(arr[0]) << 8) | uint16(arr[1])
}


//big ending
func Int16ToArrBE(v uint16) []uint8 {
	return []uint8{uint8(v>>8) & 0xFF, uint8(v) & 0xFF}
}

//big ending
func ArrToInt16BE(arr []uint8) int16 {
	if len(arr) != 2 {
		return 0
	}
	return (int16(arr[0]) << 8) | int16(arr[1])
}


//big ending
func Uint32ToArrBE(v uint32) []uint8 {
	return []uint8{uint8(v>>24) & 0xFF, uint8(v>>16) & 0xFF, uint8(v>>8) & 0xFF, uint8(v) & 0xFF}
}

//big ending
func ArrToUint32BE(arr []uint8) uint32 {
	if len(arr) != 4 {
		return 0
	}
	return (uint32(arr[0]) << 24) | (uint32(arr[1]) << 16) | (uint32(arr[2]) << 8) | uint32(arr[3])
}

//big ending
func Uint64ToArrBE(v uint64) []uint8 {
	return []uint8{
		uint8(v>>(8*7)) & 0xFF,
		uint8(v>>(8*6)) & 0xFF,
		uint8(v>>(8*5)) & 0xFF,
		uint8(v>>(8*4)) & 0xFF,
		uint8(v>>(8*3)) & 0xFF,
		uint8(v>>(8*2)) & 0xFF,
		uint8(v>>(8*1)) & 0xFF,
		uint8(v>>(8*0)) & 0xFF,
	}
}

//big ending
func ArrToUint64BE(arr []uint8) uint64 {
	if len(arr) != 8 {
		return 0
	}
	return (uint64(arr[0]) << (8 * 7)) |
		(uint64(arr[1]) << (8 * 6)) |
		(uint64(arr[2]) << (8 * 5)) |
		(uint64(arr[3]) << (8 * 4)) |
		(uint64(arr[4]) << (8 * 3)) |
		(uint64(arr[5]) << (8 * 2)) |
		(uint64(arr[6]) << (8 * 1)) |
		(uint64(arr[7]) << (8 * 0))
}
