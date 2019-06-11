package libgen

import (
	"encoding/json"
	"errors"
)

// HEADER(FE FE) 2 | CMD 2| CONTENT_TYPE 1| DATA_LEN 4| DATA N|
type IOMsg struct {
	Cmd    uint16
	Format uint8
	Body   []uint8
}

type MsgFmt uint8

const (
	JSON MsgFmt = iota + 0x01
	MSGPACK
)

const kMinMsgLen = 2 + 2 + 1 + 4

var ErrNeedMore = errors.New("codec want read more bytes")
var ErrUnknownMsgFmt = errors.New("unknown msg format")

func Decode(buf ReadableBuffer) (*IOMsg, error) {
	for {
		if buf.ReadableLen() < kMinMsgLen {
			return nil, ErrNeedMore
		}
		header := buf.PeekUInt16(0)
		if header != 0xFEFE {
			buf.PopN(1)
			continue
		}
		cmd := buf.PeekUInt16(2)
		format := buf.PeekUInt8(4)
		dataLen := buf.PeekInt32(5)
		if dataLen < 0 || dataLen > int32(buf.ReadableLen()+kMinMsgLen) {
			return nil, ErrNeedMore
		}
		buf.PopN(kMinMsgLen)
		return &IOMsg{
			Cmd:    cmd,
			Format: format,
			Body:   buf.ReadN(int(dataLen)),
		}, nil
	}
}

func Encode(cmd uint16, format MsgFmt, data interface{}) []byte {
	buffer := NewByteBuffer()
	buffer.WriteUInt16(cmd)
	buffer.WriteUInt8(uint8(format))

}
