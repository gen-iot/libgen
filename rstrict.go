package libgen

type Restrict interface {
	Type() RestrictType
	Validate(v interface{}) error
}

type baseRestrict struct {
	RestrictType RestrictType `json:"type"`
}

func (this *baseRestrict) Type() RestrictType {
	return this.RestrictType
}

func (this *baseRestrict) Validate(v interface{}) error {
	return errValidateNotSupport
}
