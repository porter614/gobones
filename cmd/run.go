package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/porter614/gobones/pkg/config"
	"github.com/porter614/logger"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the app",
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	// Logger instance
	log := logger.Instance()
	log.Info(logger.LogMessage[logger.STARTUP])

	// Exit channel (input) and done channel (output) for our main loop
	exit := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	// Load app config, and start watching it
	cfg := config.Configurator{
		Log: log,
	}
	if err := cfg.LoadConfig(App, "config"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := cfg.WatchConfig(); err != nil {
		log.Errorf("Failed to watch config: %v", err)
	}

	go func() {
		for {
			select {
			case <-exit:
				log.Info("Exiting main loop")
				done <- struct{}{}
				return
			default:
				// Main loop work goes here
				log.Debug("Log level: ", viper.Get("common.log.level"))
				log.Debug("Time: ", time.Now())
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Handle signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Info("Received signal: ", sig)

	// Tell the main loop to quit, and wait for it to do so
	exit <- struct{}{}
	<-done

	log.Info("Exiting app")
}
