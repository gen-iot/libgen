package libgen

import (
	"time"
)

type TimeoutPolicy interface {
	IsTimeout(current time.Time) bool
	Update(t time.Time)
	LastTime() time.Time
}

type CustomTimeout struct {
	statusReported bool
}

func (this *CustomTimeout) IsTimeout(time.Time) bool {
	return this.statusReported
}

func (this *CustomTimeout) Update(t time.Time) {
	if !this.statusReported {
		this.statusReported = true
	}
}

func (*CustomTimeout) LastTime() time.Time {
	return time.Time{}
}

type DefaultTimeout struct {
	timeout    time.Duration
	lastUpdate time.Time
}

func (this *DefaultTimeout) IsTimeout(current time.Time) bool {
	return current.Sub(this.lastUpdate) > this.timeout
}

func (this *DefaultTimeout) Update(t time.Time) {
	this.lastUpdate = t
}

func (this *DefaultTimeout) LastTime() time.Time {
	return this.lastUpdate
}

type NeverTimeout struct {
}

func (*NeverTimeout) IsTimeout(current time.Time) bool {
	return false
}

func (*NeverTimeout) Update(t time.Time) {
}

func (*NeverTimeout) LastTime() time.Time {
	return time.Time{}
}

func NewDefaultTimeout(dur time.Duration) TimeoutPolicy {
	out := new(DefaultTimeout)
	out.timeout = dur
	out.Update(time.Now())
	return out
}

func newTimeoutPolicy(tp TimeoutType, config []byte) TimeoutPolicy {
	switch tp {
	case Timeout:
		var timeout time.Duration = 0
		jsob, e := NewJsonObjectFromBytes(config)
		if e != nil {
			timeout = DefaultTimeoutPolicyTimeout
		} else {
			timeout = time.Duration(jsob.GetIntOr("timeout", int(DefaultTimeoutPolicyTimeout)))
		}
		return NewDefaultTimeout(timeout)
	case Custom:
		return new(CustomTimeout)
	case Never:
		return new(NeverTimeout)
	default:
		return new(NeverTimeout)
	}
}
