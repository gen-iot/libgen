package libgen

import (
	"bytes"
	"encoding/json"
	"github.com/vmihailenco/msgpack"
)

func serializeData(format MsgFmt, data interface{}) ([]byte, error) {
	switch format {
	case JSON:
		return json.Marshal(data)
	case MSGPACK:
		return MsgpackMarshal(data)
	}
	return nil, ErrUnknownMsgFmt
}

func MsgpackMarshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := msgpack.NewEncoder(buf)
	encoder.UseJSONTag(true)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}

func MsgpackUnmarshal(data []byte, v interface{}) error {
	rd := bytes.NewReader(data)
	decoder := msgpack.NewDecoder(rd)
	decoder.UseJSONTag(true)
	return decoder.Decode(v)
}

