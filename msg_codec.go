package libgen

import (
	"errors"
)

// HEADER(FE FE) 2 | CMD 2| DATA_LEN 4| DATA N|
type IOMsg struct {
	Cmd  uint16
	Data []byte
}

type GenMsg struct {
	AppId string
	IOMsg
}

const kMinMsgLen = 2 + 2 + 4

var ErrNeedMore = errors.New("codec want read more bytes")
var ErrDecode = errors.New("decode failed")

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
		dataLen := buf.PeekInt32(4)
		if dataLen > int32(buf.ReadableLen()+kMinMsgLen) {
			return nil, ErrNeedMore
		}
		return &IOMsg{
			Cmd:  cmd,
			Data: buf.ReadN(kMinMsgLen, int(kMinMsgLen+dataLen)),
		}, nil
	}
}


