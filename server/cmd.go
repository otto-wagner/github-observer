package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

func Cmd() (cmdServer *cobra.Command) {
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

	c, err := InitConfig()
	if err != nil {
		slog.Error("failed to init config", "error", err)
		os.Exit(1)
	}
	configuration := c.Config()

	cmdServer.PersistentFlags().StringVarP(&configuration.Mode, "mode", "m", configuration.Mode, "gin mode")
	cmdServer.PersistentFlags().StringVarP(&configuration.Address, "address", "a", configuration.Address, "Address of the http api endpoint")
	cmdServer.PersistentFlags().StringSliceVarP(&configuration.TrustedProxies, "trustedProxies", "t", configuration.TrustedProxies, "trusted proxies")

	// ssl
	cmdServer.PersistentFlags().StringVarP(&configuration.Ssl.Cert, "ssl.cert", "c", configuration.Ssl.Cert, "server ssl certificate")
	cmdServer.PersistentFlags().StringVarP(&configuration.Ssl.Key, "ssl.key", "k", configuration.Ssl.Key, "server ssl key")

	// app

	// => executors
	cmdServer.PersistentFlags().StringSliceVarP(&configuration.App.Executors, "app.executors", "e", configuration.App.Executors, "Executors to use (logging, prometheus)")

	// => watcher
	cmdServer.PersistentFlags().BoolVarP(&configuration.App.Watcher.Enabled, "app.watcher.enabled", "w", configuration.App.Watcher.Enabled, "enable watcher")
	cmdServer.PersistentFlags().StringVarP(&configuration.App.Watcher.GithubToken, "app.watcher.githubToken", "g", configuration.App.Watcher.GithubToken, "github token")

	// => listener
	cmdServer.PersistentFlags().BoolVarP(&configuration.App.Listener.Enabled, "app.listener.enabled", "l", configuration.App.Listener.Enabled, "enable listener")
	cmdServer.PersistentFlags().StringVarP(&configuration.App.Listener.HmacSecret, "app.listener.hmac_secret", "s", configuration.App.Listener.HmacSecret, "hmac secret")

	// => repositories
	cmdServer.PersistentFlags().StringSliceVarP(&configuration.App.Repositories, "app.repositories", "r", configuration.App.Repositories, "owner/repo@branch")

	// validation
	err = c.Validate()
	if err != nil {
		slog.Error("configuration validation failed")
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
			Run(configuration)
		},
	}
	cmdServer.AddCommand(cmd)

	return
}
