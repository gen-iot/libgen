package libgen

import (
	"encoding/json"
	"gitee.com/Puietel/std"
	"github.com/pkg/errors"
	"time"
)

type Request interface {
	Cmd() GenCommand
	MsgId() string
	Format() MsgFmt
	Data() interface{}
}

type CRequest interface {
	GetResponse(v interface{}) error
}

type SRequest interface {
	Response(v interface{}) error
}

type Response interface {
	Ok() bool
	Error() error
}

type BaseResponse struct {
}

func (this *BaseResponse) Ok() bool {
	return true
}

func (this *BaseResponse) Error() error {
	return nil
}

type BaseRequest struct {
	command GenCommand
	id      string
	msgFmt  MsgFmt
}

var ErrInvalidIoMsg = errors.New("not a valid msg")

func (this *BaseRequest) getResponse(timeout time.Duration, bodyData interface{}, rsp interface{}) error {
	p := NewPromise()
	//TODO SendMsg

	//TODO save p to collector

	f, err := p.GetFuture()
	if err != nil {
		return err
	}
	msgInterface, err := f.Wait(timeout)
	if err != nil {
		return nil
	}
	msg, ok := msgInterface.(*IOMsg)
	if !ok {
		return ErrInvalidIoMsg
	}
	data := msg.Body[:]
	switch msg.Format {
	case MSGPACK:
		{
			return MsgpackUnmarshal(data, rsp)
		}
	case JSON:
		{
			return json.Unmarshal(data, rsp)
		}
	}
	return ErrUnknownMsgFmt
}

func NewBaseRequest(command GenCommand) *BaseRequest {
	return &BaseRequest{
		command: command,
		id:      std.GenRandomUUID(),
		msgFmt:  JSON,
	}
}

func (this *BaseRequest) Cmd() GenCommand {
	return this.command
}

func (this *BaseRequest) MsgId() string {
	return this.id
}

func (this *BaseRequest) Format() MsgFmt {
	return this.msgFmt
}
