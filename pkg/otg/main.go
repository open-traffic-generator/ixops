package otg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/open-traffic-generator/tests/helpers/otg"
	"github.com/open-traffic-generator/tests/helpers/table"
)

type OtgConfig struct {
	Pps   int64
	Count int32
	Size  int32
	SIp   string
	DIp   string
	SPort int32
	DPort int32
	Tcp   bool
	Udp   bool
}

var otgConfig *OtgConfig

func NewOtgConfig() *OtgConfig {
	if otgConfig == nil {
		otgConfig = &OtgConfig{
			Pps:   100,
			Count: 1000,
			Size:  128,
			SIp:   "1.1.1.1",
			DIp:   "2.2.2.2",
			SPort: 5000,
			DPort: 6000,
		}
	}

	return otgConfig
}

func Generate(oc *OtgConfig) {
	api := otg.NewOtgApi(&testing.T{})
	defer api.CleanupConfig()

	c := getConfig(api, oc)

	api.SetConfig(c)

	api.StartTransmit()

	api.WaitFor(
		func() bool {
			ClearScreen()
			for _, m := range getFlowMetrics(api) {
				if m.Transmit() != gosnappi.FlowMetricTransmit.STOPPED {
					return false
				}
			}
			return true
		},
		&otg.WaitForOpts{
			Interval: 200 * time.Millisecond,
			Timeout:  60 * time.Second,
		},
	)
}

func getConfig(api *otg.OtgApi, oc *OtgConfig) gosnappi.Config {
	c := api.Api().NewConfig()
	p1 := c.Ports().Add().SetName("p1").SetLocation(api.TestConfig().OtgPorts[0])
	p2 := c.Ports().Add().SetName("p2").SetLocation(api.TestConfig().OtgPorts[1])

	c.Layer1().Add().
		SetName("ly").
		SetPortNames([]string{p1.Name(), p2.Name()}).
		SetSpeed(gosnappi.Layer1SpeedEnum(api.TestConfig().OtgSpeed))

	f1 := c.Flows().Add().SetName("f1")
	f1.TxRx().Port().
		SetTxName(p1.Name()).
		SetRxName(p2.Name())
	f1.Duration().FixedPackets().SetPackets(oc.Count)
	f1.Rate().SetPps(oc.Pps)
	f1.Size().SetFixed(oc.Size)
	f1.Metrics().SetEnable(true)

	eth := f1.Packet().Add().Ethernet()
	eth.Src().SetValue("00:00:00:00:00:AA")
	eth.Dst().SetValue("00:00:00:00:00:BB")

	ip := f1.Packet().Add().Ipv4()
	ip.Src().SetValue(oc.SIp)
	ip.Dst().SetValue(oc.DIp)

	if oc.Tcp {
		tcp := f1.Packet().Add().Tcp()
		tcp.SrcPort().SetValue(oc.SPort)
		tcp.DstPort().SetValue(oc.DPort)
	} else if oc.Udp {
		udp := f1.Packet().Add().Udp()
		udp.SrcPort().SetValue(oc.SPort)
		udp.DstPort().SetValue(oc.DPort)
	}

	return c
}

func getFlowMetrics(api *otg.OtgApi) []gosnappi.FlowMetric {

	mr := api.Api().NewMetricsRequest()
	mr.Flow()
	res, err := api.Api().GetMetrics(mr)
	api.LogWrnErr(nil, err, true)

	tb := table.NewTable(
		"Flow Metrics",
		[]string{
			"Name",
			"State",
			"Frames Tx",
			"Frames Rx",
			"FPS Tx",
			"FPS Rx",
			"Bytes Tx",
			"Bytes Rx",
		},
		15,
	)
	for _, v := range res.FlowMetrics().Items() {
		if v != nil {
			tb.AppendRow([]interface{}{
				v.Name(),
				v.Transmit(),
				v.FramesTx(),
				v.FramesRx(),
				v.FramesTxRate(),
				v.FramesRxRate(),
				v.BytesTx(),
				v.BytesRx(),
			})
		}
	}

	fmt.Println(tb)
	return res.FlowMetrics().Items()
}

func ClearScreen() {
	switch runtime.GOOS {
	case "darwin":
		fallthrough
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	default:
	}
}
