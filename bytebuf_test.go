package libgen

import (
	"libgen/liblpc/backend"
	"testing"
)

func TestByteBuffer(t *testing.T) {
	buffer := NewByteBuffer()
	buffer.Write([]byte{0x01, 0x0e})
	backend.Assert(buffer.PeekInt8(1) == int8(0x0e), "bad")
	backend.Assert(buffer.ReadInt16() == 0x010e, "bad")
	backend.Assert(buffer.ReadableLen() == 0, "bad")
}
