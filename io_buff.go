package libgen

import (
	"container/list"
)

type ReadableBuffer interface {
	ReadInt32(offset int) int32
	ReadUInt32(offset int) uint32
	ReadInt16(offset int) int16
	ReadUInt16(offset int) uint16
	ReadInt8(offset int) int8
	ReadUInt8(offset int) uint8

	PeekInt32(offset int) int32
	PeekUInt32(offset int) uint32
	PeekInt16(offset int) int16
	PeekUInt16(offset int) uint16
	PeekInt8(offset int) int8
	PeekUInt8(offset int) uint8

	ReadableLen() int

	ReadN(offset int, n int) []uint8
	PopN(n int)
}

type ByteBuffer struct {
	data *list.List
}

func NewIoBuffer() *ByteBuffer {
	return &ByteBuffer{
		data: list.New(),
	}
}

func (this *ByteBuffer) ToArray() []uint8 {
	ret := make([]uint8, 0, this.ReadableLen())
	for ele := this.data.Front(); ele != nil; ele = ele.Next() {
		ret = append(ret, ele.Value.(uint8))
	}
	return ret
}

func (this *ByteBuffer) Write(arr []uint8) {
	for _, v := range arr {
		this.data.PushBack(v)
	}
}

//checkout ReadableLen before call this
func (this *ByteBuffer) PeekN(offset, n int) []uint8 {
	arr := make([]uint8, 0, n)
	var ele = this.data.Front()
	for i := 0; i < offset; ele = ele.Next() {
		if ele == nil {
			return arr
		}
		i++
	}
	for i, it := 0, ele; i < n && it != nil; it = it.Next() {
		arr = append(arr, it.Value.(uint8))
		i++
	}
	return arr
}

//checkout ReadableLen before call this
func (this *ByteBuffer) ReadN(n int) []uint8 {
	arr := make([]uint8, 0, n)
	dataN := this.data.Len()
	for i := 0; i < n && i < dataN; i++ {
		front := this.data.Front()
		if front == nil {
			break
		}
		arr = append(arr, front.Value.(uint8))
		this.data.Remove(front)
	}
	return arr
}

//array len
func (this *ByteBuffer) ReadableLen() int {
	return this.data.Len()
}

func (this *ByteBuffer) PopN(n int) {
	for i := 0; i < n; i++ {
		front := this.data.Front()
		if front != nil {
			this.data.Remove(this.data.Front())
		}
	}
}
