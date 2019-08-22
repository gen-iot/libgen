package libgen

type StatusProperty struct {
	Restrict Restrict `json:"restrict" validate:"omitempty,required"`
}

type CommandProperty struct {
	Command        string     `json:"command"`
	ParamRestricts []Restrict `json:"restricts" validate:"omitempty,dive,required"`
}

func NewCommandProperty(funcName string, params ...Restrict) *CommandProperty {
	return &CommandProperty{
		Command:        funcName,
		ParamRestricts: params,
	}
}

func (this *CommandProperty) AddParamRestrict(r Restrict) {
	this.ParamRestricts = append(this.ParamRestricts, r)
}
