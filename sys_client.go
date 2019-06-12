package libgen

import (
	"fmt"
	"gitee.com/SuzhenProjects/liblpc"
	"net"
	"os"
	"sync"
	"time"
)

var gClientConn net.Conn
var rdBuf = make([]byte, 1024*1024*2)
var rdCache = NewByteBuffer()
var initOnce = sync.Once{}

const clientFd = uintptr(3)

func Init() {
	initOnce.Do(doInit)
}

func doInit() {
	fmt.Println("LIBGEN CLIENT INIT")
	file := os.NewFile(clientFd, "")
	c, err := net.FileConn(file)
	liblpc.PanicIfError(err)
	gClientConn = c
	fmt.Println("LIBGEN CLIENT INIT SUCCESS")
}

func SendMsg(timeout time.Duration, cmd uint16, format MsgFmt, data interface{}) error {
	err := gClientConn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	bytes, err := Encode(cmd, format, data)
	if err != nil {
		return err
	}
	dataLen := len(bytes)
	nWrite := 0
	for {
		nw, err := gClientConn.Write(bytes)
		if err != nil {
			return err
		}
		nWrite += nw
		if nWrite == dataLen {
			break
		}
		bytes = bytes[nWrite:]
	}
	return nil
}

func RecvMsg(timeout time.Duration, maxDataSize int) ([]*IOMsg, error) {
	err := gClientConn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}

	n, err := gClientConn.Read(rdBuf)
	if err != nil {
		return nil, nil
	}
	rdCache.Write(rdBuf[:n])
	out := make([]*IOMsg, 0, 16)
	for {
		msg, err := Decode(rdCache, maxDataSize)
		if err != nil {
			break
		}
		out = append(out, msg)
	}
	return out, nil
}
