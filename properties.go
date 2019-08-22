package libgen

type StatusProperty struct {
	Restrict Restrict `json:"restrict" validate:"omitempty,required"`
}

type ControlProperty struct {
	CtrlFuncName   string     `json:"ctrlFuncName"`
	ParamRestricts []Restrict `json:"restricts" validate:"omitempty,dive,required"`
}

func NewControlProperty(funcName string, params ...Restrict) *ControlProperty {
	return &ControlProperty{
		CtrlFuncName:   funcName,
		ParamRestricts: params,
	}
}

func (this *ControlProperty) AddParamRestrict(r Restrict) {
	this.ParamRestricts = append(this.ParamRestricts, r)
}
