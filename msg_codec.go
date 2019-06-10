package libgen

import (
	"errors"
)

// HEADER(FE FE) 2 | CMD 2| CONTENT_TYPE 1| DATA_LEN 4| DATA N|
type IOMsg struct {
	Cmd    uint16
	Format uint8
	Body   []uint8
}

type GenMsg struct {
	AppId string
	IOMsg
}

const kMinMsgLen = 2 + 2 + 1 + 4

var ErrNeedMore = errors.New("codec want read more bytes")

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
		contentType := buf.PeekUInt8(4)
		dataLen := buf.PeekInt32(5)
		if dataLen > int32(buf.ReadableLen()+kMinMsgLen) {
			return nil, ErrNeedMore
		}
		buf.PopN(kMinMsgLen)
		return &IOMsg{
			Cmd:    cmd,
			Format: contentType,
			Body:   buf.ReadN(int(dataLen)),
		}, nil
	}
}
