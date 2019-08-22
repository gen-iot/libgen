package libgen

type Restrict interface {
	Type() RestrictType
	Validate(v interface{}) error
	IsRequired() bool
}

type baseRestrict struct {
	RestrictType RestrictType `json:"type"`
	Required     bool         `json:"required"`
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
