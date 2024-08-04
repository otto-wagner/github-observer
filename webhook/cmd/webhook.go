package cmd

import (
	"github-observer/conf"
	"github-observer/webhook"
	"github.com/go-playground/validator/v10"
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

	var cfgFile string
	cmdWebhook.PersistentFlags().StringVarP(&cfgFile, "config", "c", "conf/webhook.json", "config file (default is conf/webhook.json)")
	configuration, err := conf.InitWebhook(cfgFile)
	if err != nil {
		slog.Error("failed to init config", "error", err)
		os.Exit(1)
	}
	validationConfiguration(configuration)

	create := &cobra.Command{
		Use:   "create",
		Short: "create webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhook.Create(configuration)
		},
	}
	cmdWebhook.AddCommand(create)

	list := &cobra.Command{
		Use:   "list",
		Short: "list webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhook.List(configuration)
		},
	}
	cmdWebhook.AddCommand(list)

	delete := &cobra.Command{
		Use:   "delete",
		Short: "delete webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhook.Delete(configuration)
		},
	}
	cmdWebhook.AddCommand(delete)

	return
}

func validationConfiguration(configWebhook conf.WebHookConfig) {
	err := configWebhook.Validate()
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
		os.Exit(1)
	}
}
