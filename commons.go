package libgen

import (
	"time"
)

type AppType int

const (
	LocalApp  AppType = 900
	RemoteApp AppType = 901
)

func AppType2Str(t AppType) string {
	switch t {
	case LocalApp:
		return "LOCAL"
	case RemoteApp:
		return "REMOTE"
	default:
		return "UNKNOWN"
	}
}

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
