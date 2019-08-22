package libgen

type StatusProperty struct {
	Name     string   `json:"name" validate:"required"`
	Restrict Restrict `json:"restrict" validate:"required"`
}

type ControlProperty struct {
	Name           string     `json:"name" validate:"required"`
	ParamRestricts []Restrict `json:"restricts" validate:"omitempty,dive,required"`
}

func NewRestrictGroup(name string, params ...Restrict) *ControlProperty {
	return &ControlProperty{
		Name:           name,
		ParamRestricts: params,
	}
}

func (this *ControlProperty) AddParamRestrict(r Restrict) {
	this.ParamRestricts = append(this.ParamRestricts, r)
}
