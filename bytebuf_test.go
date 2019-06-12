package libgen

import (
	"gitee.com/SuzhenProjects/liblpc"
	"testing"
)

func TestByteBuffer(t *testing.T) {
	buffer := NewByteBuffer()
	buffer.Write([]byte{0x01, 0x0e})
	liblpc.Assert(buffer.PeekInt8(1) == int8(0x0e), "bad")
	liblpc.Assert(buffer.ReadInt16() == 0x010e, "bad")
	liblpc.Assert(buffer.ReadableLen() == 0, "bad")

	buffer.WriteUInt16(0x010e)
	liblpc.Assert(buffer.PeekInt8(1) == int8(0x0e), "bad")
	liblpc.Assert(buffer.ReadInt16() == 0x010e, "bad")
	liblpc.Assert(buffer.ReadableLen() == 0, "bad")

	buffer.WriteInt32(0x010e)
	liblpc.Assert(buffer.PeekInt8(buffer.ReadableLen()-1) == int8(0x0e), "bad")
	liblpc.Assert(buffer.ReadInt32() == 0x010e, "bad")
	liblpc.Assert(buffer.ReadableLen() == 0, "bad")
}
