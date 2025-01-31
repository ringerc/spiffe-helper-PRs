package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spiffe/spiffe-helper/cmd/spiffe-helper/config"
	"github.com/spiffe/spiffe-helper/pkg/health"
	"github.com/spiffe/spiffe-helper/pkg/sidecar"
)

func main() {
	cliFlags := config.NewCLIFlags()
	cliFlags.Parse()
	log := logrus.WithField("system", "spiffe-helper")

	if configFile := cliFlags.ConfigFile(); configFile == "" {
		log.Error("No configuration file provided")
		os.Exit(1)
	} else {
		log.Infof("Using configuration file: %q", configFile)
	}

	hclConfig, err := config.ParseConfig(log, cliFlags)
	if err != nil {
		log.WithError(err).Errorf("failed to parse configuration")
		os.Exit(1)
	}

	if err := hclConfig.ValidateConfig(log); err != nil {
		log.WithError(err).Errorf("invalid configuration")
		os.Exit(1)
	}

	if err = startSidecar(hclConfig, log); err != nil {
		log.WithError(err).Errorf("Error starting spiffe-helper")
		os.Exit(1)
	}

	log.Infof("Exiting")
	os.Exit(0)
}

func startSidecar(hclConfig *config.Config, log logrus.FieldLogger) error {
	sidecarConfig := config.NewSidecarConfig(hclConfig, log)
	spiffeSidecar := sidecar.New(sidecarConfig)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if *hclConfig.DaemonMode && hclConfig.HealthCheck.ListenerEnabled {
		log.Info("Starting health server")
		if err := health.StartHealthServer(hclConfig.HealthCheck, log, spiffeSidecar); err != nil {
			return err
		}
	}

	if !*hclConfig.DaemonMode {
		log.Info("Daemon mode disabled")
		return spiffeSidecar.Run(ctx)
	}

	log.Info("Launching daemon")
	return spiffeSidecar.RunDaemon(ctx)
}
