package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"testing"
)

type exampleStruct struct {
	Name string                 `json:"name"`
	Age  int                    `json:"age"`
	Meta map[string]interface{} `json:"meta"`
}

func newExampleStruct() *exampleStruct {
	return &exampleStruct{
		Name: "suzhen",
		Age:  100,
		Meta: map[string]interface{}{
			"k1": "v1",
			"k2": 1,
			"k3": []string{"a", "b", "c"},
		},
	}
}

func TestEncodeMessage_JSON(t *testing.T) {
	o := newExampleStruct()
	bytes, err := Encode(1, JSON, o)
	std.AssertError(err, "Encode")
	fmt.Println(string(bytes))
}

func TestEncodeMessage_MSGPACK(t *testing.T) {
	o := newExampleStruct()
	bytes, err := Encode(1, MSGPACK, o)
	std.AssertError(err, "Encode")
	fmt.Println(string(bytes))
	buffer := std.NewByteBuffer()
	buffer.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 9, 90})
	buffer.Write(bytes)
	msg, err := Decode(buffer, 1024*1024)
	std.AssertError(err, "Decode")
	fmt.Println(*msg)
}

func TestUUID(t *testing.T) {
	uuid := std.GenRandomUUID()
	fmt.Println("uuid -> ", uuid, " len -> ", len(uuid))
}
