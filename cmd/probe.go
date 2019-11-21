package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Readiness and Liveness probe",
	Long:  "Probe to ensure liveness and readiness for app.",
	Run:   probe,
}

func init() {
	RootCmd.AddCommand(probeCmd)
	probeCmd.PersistentFlags().String("type", "live", "Type of probe: ready or liveness")
}

func probe(cmd *cobra.Command, args []string) {
	//Add tests here
	ty := cmd.Flag("type").Value.String()
	log.Info("Type: ", ty)
	if ty == "ready" {
		time.Sleep(500 * time.Millisecond)
		log.Info("I'm ready to rumble.")
	} else if ty == "live" {
		log.Info("Ah ah ah ah stayin alive.")
	}
}
