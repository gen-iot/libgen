package libgen

type Encoding int

const (
	MSGPACK Encoding = iota + 10
	JSON
)

type Message interface {
	Command() int
	Encoding() Encoding
	Data() []byte
}

func SysSendMsg(msg Message) error {
	if msg == nil {
		return ErrIllegalParam
	}
	return nil
}

func SysRecvMsg() (Message, error) {
	return nil, nil
}

