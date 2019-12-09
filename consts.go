package libgen

import "time"

const DefaultRoomName = "Default"
const DefaultFunctionName = "default"

type AppType int

const (
	LocalApp  AppType = 900
	RemoteApp AppType = 901
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
	Array    RestrictType = 1005
	Any      RestrictType = 2000
)
