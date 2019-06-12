package libgen

import (
	"errors"
	"libgen/liblpc/backend"
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

func Decode(buf ReadableBuffer, maxBodyLen int) (*IOMsg, error) {
	backend.Assert(maxBodyLen > 0, "maxBodyLen must > 0")
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
		if dataLen < 0 {
			buf.PopN(1)
			continue
		}
		if dataLen > int32(buf.ReadableLen()+kMinMsgLen) {
			return nil, ErrNeedMore
		}
		if int(dataLen) > maxBodyLen {
			buf.PopN(kMinMsgLen + int(dataLen))
			continue
		}
		buf.PopN(kMinMsgLen)
		return &IOMsg{
			Cmd:    cmd,
			Format: format,
			Body:   buf.ReadN(int(dataLen)),
		}, nil
	}
}

func Encode(cmd uint16, format MsgFmt, data interface{}) ([]byte, error) {
	buffer := NewByteBuffer()
	buffer.WriteUInt16(0xFEFE)
	buffer.WriteUInt16(cmd)
	buffer.WriteUInt8(uint8(format))
	dataLen := 0
	var dataBytes []byte = nil
	var err error = nil
	if data != nil {
		dataBytes, err = serializeData(format, data)
		if err != nil {
			return nil, err
		}
		dataLen = len(dataBytes)
	}
	buffer.WriteInt32(int32(dataLen))
	if len(dataBytes) != 0 {
		buffer.Write(dataBytes)
	}
	return buffer.ToArray(), nil
}
