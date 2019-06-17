package libgen

import (
	"errors"
	"gitee.com/Puietel/std"
	"log"
)

type MsgFmt uint8

const (
	JSON MsgFmt = iota + 0x01
	MSGPACK
)

// HEADER(FE FE) 2 |DATA_LEN 4| DATA N|

const kHeaderLen = 2
const kDataLen = 4

const kHeaderOffset = 0
const kDataLenOffset = kHeaderOffset + kHeaderLen
const kDataOffset = kDataLenOffset + kDataLen

const kMinMsgLen = kDataOffset

var ErrNeedMore = errors.New("codec want read more bytes")

type rpcRawMsg struct {
	Id         string `json:"msgId"`
	MethodName string `json:"methodName"`
	Type       int    `json:"type"` // req or ack
	Err        string `json:"err"`  // fast path for ack error
	Data       []byte `json:"data"` // req = param
}

func Decode(buf std.ReadableBuffer, maxBodyLen int) (*rpcRawMsg, error) {
	std.Assert(maxBodyLen > 0, "maxBodyLen must > 0")
	for {
		if buf.ReadableLen() < kMinMsgLen {
			return nil, ErrNeedMore
		}
		header := buf.PeekUInt16(kHeaderOffset)
		if header != 0xFEFE {
			buf.PopN(1)
			continue
		}
		dataLen := buf.PeekInt32(kDataLenOffset)
		if dataLen < 0 {
			buf.PopN(1)
			continue
		}
		if dataLen > int32(buf.ReadableLen()-kMinMsgLen) {
			return nil, ErrNeedMore
		}
		if int(dataLen) > maxBodyLen {
			buf.PopN(kMinMsgLen + int(dataLen))
			continue
		}
		buf.PopN(kDataOffset)
		data := buf.ReadN(int(dataLen))
		outMsg := new(rpcRawMsg)
		err := std.MsgpackUnmarshal(data, outMsg)
		if err != nil {
			log.Println("unmarshal rpc msg failed -> ", err)
			continue
		}
		return outMsg, nil
	}
}

func Encode(msg *rpcRawMsg) ([]byte, error) {
	std.Assert(len(msg.Id) == 32, "msgId.Len != 32")
	buffer := std.NewByteBuffer()
	datas, err := std.MsgpackMarshal(msg)
	if err != nil {
		return nil, err
	}
	buffer.WriteUInt16(0xFEFE)
	buffer.WriteInt32(int32(len(datas)))
	buffer.Write(datas)
	return buffer.ToArray(), nil
}
