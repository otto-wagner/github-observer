package cmd

import (
	"github-listener/internal/listener"
	"github-listener/internal/router"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Start server",
	Long: `You can configure the HTTP Server via the command line flags, define a configuration file (eg. JSON, TOML, YAML)
			JSON Config File Example: conf/common.json`,
	Run: startServer,
}

func init() {
	// config
	serverCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is conf/common.json)")
	// app
	serverCmd.PersistentFlags().StringP("app.listenAddress", "l", "0.0.0.0:8443", "Action address of the HTTP API endpoint")
	//serverCmd.PersistentFlags().StringSliceP("app.trustedProxies", "t", []string{""}, "Action address of the HTTP API endpoint")
	serverCmd.PersistentFlags().StringP("app.mode", "g", "release", "Gin mode")
	serverCmd.PersistentFlags().StringP("app.executors", "e", "logging", "Executors to execute on event")
	// ssl
	serverCmd.PersistentFlags().StringP("ssl.cert", "s", "conf/ssl.cert", "Path to SSL certificate")
	serverCmd.PersistentFlags().StringP("ssl.key", "k", "conf/ssl.key", "Path to SSL key")

	// bind flags
	err := viper.BindPFlags(serverCmd.PersistentFlags())
	if err != nil {
		zap.S().Errorw("failed to bind flags", "error", err)
	}
}

func startServer(_ *cobra.Command, _ []string) {
	err := configuration.Validate()
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			if e.ActualTag() == "required" {
				zap.S().Errorw("Missing required configuration value", "field", e.StructNamespace())
			} else if e.Param() != "" {
				zap.S().Errorw("Validation failed", "field", e.StructNamespace(), "tag", e.ActualTag(), "param", e.Param())
			} else {
				zap.S().Errorw("Validation failed", "field", e.StructNamespace(), "tag", e.ActualTag())
			}
		}
		zap.S().Fatal("configuration validation failed")
	}

	router.InitializeRoutes(engine, listener.NewListener(executors))

	err = engine.Run(configuration.App.ListenAddress)
	if err != nil {
		zap.S().Fatalw("failed to start httpServer", "error", err)
	}
}
