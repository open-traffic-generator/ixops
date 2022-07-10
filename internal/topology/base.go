package topology

type TopologyType string

const (
	OTG_DUT_OTG       TopologyType = "otg-dut-otg"
	OTG_DUT_DUT_2_OTG TopologyType = "otg-dut-dut-2otg"
)

const (
	TopologyFile         string = "topo.yaml"
	AristaInitConfigFile string = "init_dut.txt"
)
