package libgen

import (
	"errors"
	"gitee.com/Puietel/std"
)

// HEADER(FE FE) 2 | ID 32| CMD 2| CONTENT_TYPE 1| DATA_LEN 4| DATA N|
type IOMsg struct {
	Id     string
	Cmd    uint16
	Format uint8
	Body   []uint8
}

type MsgFmt uint8

const (
	JSON MsgFmt = iota + 0x01
	MSGPACK
)

const kMinMsgLen = 2 + 32 + 2 + 1 + 4

var ErrNeedMore = errors.New("codec want read more bytes")
var ErrUnknownMsgFmt = errors.New("unknown msg format")

func Decode(buf std.ReadableBuffer, maxBodyLen int) (*IOMsg, error) {
	std.Assert(maxBodyLen > 0, "maxBodyLen must > 0")
	for {
		if buf.ReadableLen() < kMinMsgLen {
			return nil, ErrNeedMore
		}
		header := buf.PeekUInt16(0)
		if header != 0xFEFE {
			buf.PopN(1)
			continue
		}
		msgId := string(buf.PeekN(2, 32))
		cmd := buf.PeekUInt16(2 + 32)
		format := buf.PeekUInt8(2 + 32 + 2)
		dataLen := buf.PeekInt32(2 + 32 + 2 + 1)
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
			Id:     msgId,
			Cmd:    cmd,
			Format: format,
			Body:   buf.ReadN(int(dataLen)),
		}, nil
	}
}

func Encode(msgId string, cmd uint16, format MsgFmt, data interface{}) ([]byte, error) {
	std.Assert(len(msgId) == 32, "msgId.Len != 32")
	buffer := std.NewByteBuffer()
	buffer.WriteUInt16(0xFEFE)
	buffer.Write([]byte(msgId))
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
