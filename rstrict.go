package libgen

import "github.com/gen-iot/std"

type Restrict interface {
	Name() string
	Type() RestrictType
	Validate(v interface{}) error
	IsRequired() bool
	DefaultValue() (exist bool, v interface{})
	SetDefaultValue(enable bool, v interface{}) Restrict
	SetRequired(required bool) Restrict
}

type baseRestrict struct {
	RestrictName string       `json:"name" validate:"required"`
	RestrictType RestrictType `json:"type" validate:"required"`
	Required     bool         `json:"required"`
	DefaultVal   interface{}  `json:"default"`
	HasDefault   bool         `json:"hasDefault"`
	drivenChild  Restrict     `json:"-"`
}

func newBaseRestrict(child Restrict, name string, tp RestrictType) baseRestrict {
	std.Assert(child != nil, "child restrict is nil")
	return baseRestrict{
		RestrictName: name,
		RestrictType: tp,
		Required:     true,
		drivenChild:  child,
	}
}

func (this *baseRestrict) Name() string {
	return this.RestrictName
}

func (this *baseRestrict) Type() RestrictType {
	return this.RestrictType
}

func (this *baseRestrict) Validate(v interface{}) error {
	return errValidateNotSupport
}

func (this *baseRestrict) IsRequired() bool {
	return this.Required
}

func (this *baseRestrict) DefaultValue() (exist bool, v interface{}) {
	return this.HasDefault, this.DefaultVal
}

func (this *baseRestrict) SetDefaultValue(enable bool, v interface{}) Restrict {
	this.HasDefault = enable
	if this.HasDefault {
		this.DefaultVal = v
		err := this.drivenChild.Validate(v)
		std.AssertError(err, "validate default value failed")
	} else {
		this.DefaultVal = nil
	}
	this.Required = !this.HasDefault
	return this.drivenChild
}
func (this *baseRestrict) SetRequired(required bool) Restrict {
	this.Required = required
	return this.drivenChild
}
