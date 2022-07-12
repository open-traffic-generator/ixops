package otg

import (
	"github.com/open-traffic-generator/ixops/pkg/otg"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:          "otg",
	Short:        "Manage OTG",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	oc := otg.NewOtgConfig()
	configCmd.PersistentFlags().BoolVar(&oc.Udp, "udp", oc.Udp, "Configure UDP")
	configCmd.PersistentFlags().BoolVar(&oc.Tcp, "tcp", oc.Tcp, "Configure TCP")
	configCmd.PersistentFlags().Int64Var(&oc.Pps, "pps", oc.Pps, "Packet rate")
	configCmd.PersistentFlags().Int32Var(&oc.Size, "size", oc.Size, "Packet size in bytes")
	configCmd.PersistentFlags().Int32Var(&oc.Count, "count", oc.Count, "Packet count")
	configCmd.PersistentFlags().StringVar(&oc.SIp, "sip", oc.SIp, "Source IP address")
	configCmd.PersistentFlags().StringVar(&oc.DIp, "dip", oc.DIp, "Destination IP address")
	configCmd.PersistentFlags().Int32Var(&oc.SPort, "sport", oc.SPort, "Source port")
	configCmd.PersistentFlags().Int32Var(&oc.DPort, "dport", oc.DPort, "Destination port")
	configCmd.AddCommand(genCmd)
	return configCmd
}
