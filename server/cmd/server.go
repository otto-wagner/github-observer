package cmd

import (
	"github-observer/conf"
	"github-observer/server"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func Server() (cmdServer *cobra.Command) {
	cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Start server",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Usage()
			if err != nil {
				return
			}
		},
	}

	configuration, err := conf.InitCommon()
	if err != nil {
		slog.Error("failed to init config", "error", err)
		os.Exit(1)
	}
	validationConfiguration(configuration)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "run server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Run(configuration)
		},
	}
	cmdServer.AddCommand(cmd)

	return
}

func validationConfiguration(configuration conf.Config) {
	err := configuration.Validate()
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			if e.ActualTag() == "required" {
				slog.Error("Missing required configuration value", "field", e.StructNamespace())
			} else if e.Param() != "" {
				slog.Error("Validation failed", "field", e.StructNamespace(), "tag", e.ActualTag(), "param", e.Param())
			} else {
				slog.Error("Validation failed", "field", e.StructNamespace(), "tag", e.ActualTag())
			}
		}
		slog.Error("configuration validation failed")
	}
}
