package cmd

import (
	"github-observer/conf"
	"github-observer/webhook"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func WebHook() (cmdWebhook *cobra.Command) {
	cmdWebhook = &cobra.Command{
		Use:   "webhook",
		Short: "api to update webhook",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Usage()
			if err != nil {
				return
			}
		},
	}

	configuration, err := conf.InitWebhook()
	if err != nil {
		slog.Error("failed to init config", "error", err)
		os.Exit(1)
	}
	err = configuration.Validate()
	if err != nil {
		slog.Error("configuration validation failed", "error", err)
		os.Exit(1)
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhook.Create(configuration.Config())
		},
	}
	cmdWebhook.AddCommand(create)

	list := &cobra.Command{
		Use:   "list",
		Short: "list webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhook.List(configuration.Config())
		},
	}
	cmdWebhook.AddCommand(list)

	deleteWebhook := &cobra.Command{
		Use:   "delete",
		Short: "delete webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhook.Delete(configuration.Config())
		},
	}
	cmdWebhook.AddCommand(deleteWebhook)

	return
}
