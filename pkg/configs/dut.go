package configs

type Dut struct {
	Image *string `yaml:"image"`
}

func (v *Dut) SetDefaults() {
}
