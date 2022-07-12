package configs

type EndpointKind string

const (
	EndpointKindIxiaC EndpointKind = "ixia_c"
	EndpointKindDut   EndpointKind = "dut"
)

func (t *EndpointKind) IsValid() bool {
	switch *t {
	case EndpointKindIxiaC:
		return true
	case EndpointKindDut:
		return true
	}

	return false
}

type Endpoint struct {
	Name  *string       `yaml:"name"`
	Kind  *EndpointKind `yaml:"kind"`
	IxiaC *IxiaC        `yaml:"ixia_c"`
	Dut   *Dut          `yaml:"dut"`
}

func (v *Endpoint) SetDefaults() {
	SetDefaultString(&v.Name, "otg")
	if v.Kind == nil {
		k := EndpointKind("ixia_c")
		v.Kind = &k
	}
	if v.IxiaC == nil {
		v.IxiaC = &IxiaC{}
	}
	v.IxiaC.SetDefaults()
}
