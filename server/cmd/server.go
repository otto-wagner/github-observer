package cmd

import (
	"github-observer/conf"
	"github-observer/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	config, err := conf.InitCommon()
	if err != nil {
		slog.Error("failed to init config", "error", err)
		os.Exit(1)
	}
	configuration := config.Config()

	// app
	cmdServer.PersistentFlags().StringVarP(&configuration.App.ListenAddress, "app.listenAddress", "l", configuration.App.ListenAddress, "Listen address of the HTTP API endpoint")
	cmdServer.PersistentFlags().StringSliceVarP(&configuration.App.TrustedProxies, "app.trustedProxies", "t", configuration.App.TrustedProxies, "Listen address of the HTTP API endpoint")
	cmdServer.PersistentFlags().StringVarP(&configuration.App.Mode, "app.mode", "m", configuration.App.Mode, "Gin mode")
	cmdServer.PersistentFlags().BoolVarP(&configuration.App.Watcher, "app.watcher", "w", configuration.App.Watcher, "Watcher")
	cmdServer.PersistentFlags().StringSliceVarP(&configuration.App.Executors, "app.executors", "e", configuration.App.Executors, "Executors")
	cmdServer.PersistentFlags().StringSliceVarP(&configuration.App.Logs, "app.logs", "o", configuration.App.Logs, "Logs")

	// ssl
	cmdServer.PersistentFlags().BoolVarP(&configuration.Ssl.Activate, "ssl.activate", "a", configuration.Ssl.Activate, "Activate SSL")
	cmdServer.PersistentFlags().StringVarP(&configuration.Ssl.Cert, "ssl.cert", "c", configuration.Ssl.Cert, "Cert")
	cmdServer.PersistentFlags().StringVarP(&configuration.Ssl.Key, "ssl.key", "k", configuration.Ssl.Key, "Key")

	// secret
	cmdServer.PersistentFlags().StringVarP(&configuration.HmacSecret, "secret", "s", configuration.HmacSecret, "HmacSecret")

	// repositories
	cmdServer.PersistentFlags().StringSliceVarP(&configuration.Repositories, "repositories", "r", configuration.Repositories, "owner/repo@branch")

	// validation
	err = config.Validate()
	if err != nil {
		slog.Error("configuration validation failed", "error", err)
		os.Exit(1)
	}

	err = viper.BindPFlags(cmdServer.PersistentFlags())
	if err != nil {
		slog.Error("failed to bind flags", "error", err)
		os.Exit(1)
	}

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
