package libgen

type Restrict interface {
	Name() string
	Type() RestrictType
	Validate(v interface{}) error
	IsRequired() bool
}

type baseRestrict struct {
	RestrictName string       `json:"name" validate:"required"`
	RestrictType RestrictType `json:"type" validate:"required"`
	Required     bool         `json:"required"`
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
