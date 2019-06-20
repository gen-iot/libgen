package libgen

import (
	"encoding/json"
	"time"
)

type TimeoutType int

const (
	Timeout TimeoutType = iota
	Custom
	Never
)

const (
	DefaultTimeoutPolicyTimeout = time.Second * 60
)

type RestrictType int

const (
	I32Range RestrictType = 1000
	I32Enum  RestrictType = 1001
	I32Any   RestrictType = 1002
	StrEnum  RestrictType = 1003
	StrAny   RestrictType = 1004
	Any      RestrictType = 2000
)

type PropertyType string

const (
	Status  PropertyType = "status"
	Command PropertyType = "command"
)

type JsonObject map[string]interface{}

func NewJsonObjectFromBytes(data []byte) (JsonObject, error) {
	out := make(JsonObject)
	err := json.Unmarshal(data, &out)
	return out, err
}

func (this JsonObject) GetIntOr(key string, dft int) int {
	i, ok := this[key]
	if !ok {
		return dft
	}
	i2, ok := i.(int)
	if !ok {
		return dft
	}
	return i2
}
