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

	var cfgFile string
	// config
	cmdServer.PersistentFlags().StringVarP(&cfgFile, "config", "c", "conf/common.json", "config file (default is conf/common.json)")

	// app
	//cmdServer.PersistentFlags().StringP("app.listenAddress", "l", "", "Listen address of the HTTP endpoint")
	//cmdServer.PersistentFlags().StringSliceP("app.trustedProxies", "t", []string{""}, "Action address of the HTTP API endpoint")
	//cmdServer.PersistentFlags().StringP("app.mode", "g", "debug", "Gin mode")
	//cmdServer.PersistentFlags().StringSliceP("app.executors", "e", []string{""}, "Executors to execute on event")

	//cmdServer.PersistentFlags().BoolP("app.watcher", "w", configuration.App.Watcher, "Enable watcher")

	// ssl
	//cmdServer.PersistentFlags().StringP("ssl.cert", "s", "conf/ssl.cert", "Path to SSL certificate")
	//cmdServer.PersistentFlags().StringP("ssl.key", "k", "conf/ssl.key", "Path to SSL key")

	configuration, err := conf.InitCommon(cfgFile)
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
