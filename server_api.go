package libgen

import (
	"gitee.com/Puietel/std"
	"reflect"
)

type ServerWriteHandle func(data []byte, userData interface{})
type ServerOnError func(err error, userData interface{})

type Server struct {
	API
	writeHandle ServerWriteHandle
	errHandle   ServerOnError
	apiFuncMap  map[GenCommand]reflect.Value
}

func NewServer(api API, writeHandle ServerWriteHandle, errHandle ServerOnError) *Server {
	s := new(Server)
	s.API = api
	s.writeHandle = writeHandle
	s.errHandle = errHandle
	return s
}

func (this *Server) makeApiFuncMap() {
	of := reflect.TypeOf(this.API)
	for i := 0; i < of.NumMethod(); i++ {
		method := of.Method(i)
	}
}

func (this *Server) OnDataRead(buf std.ReadableBuffer, userData interface{}) {
	for {
		msg, err := Decode(buf, 1024*1024*1)
		if err != nil {
			if err == ErrNeedMore {
				break
			}
			this.errHandle(err, userData)
			return
		}
		//decode msg
	}

}
