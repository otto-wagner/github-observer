package webhook

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strings"
)

func Cmd() (cmdWebhook *cobra.Command) {
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

	configuration, err := InitConfig()
	if err != nil {
		slog.Error("failed to init config", "error", err)
		os.Exit(1)
	}
	err = configuration.Validate()
	if err != nil {
		slog.Error("configuration validation failed")
		os.Exit(1)
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			Create(configuration.Config())
		},
	}
	cmdWebhook.AddCommand(create)

	list := &cobra.Command{
		Use:   "list",
		Short: "list webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			List(configuration.Config())
		},
	}
	cmdWebhook.AddCommand(list)

	deleteWebhook := &cobra.Command{
		Use:   "delete",
		Short: "delete webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			Delete(configuration.Config())
		},
	}
	cmdWebhook.AddCommand(deleteWebhook)

	return
}

type Repository struct {
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Branch string `json:"branch"`
}

func validateRepositories(fl validator.FieldLevel) bool {
	repositories := fl.Field().Interface().([]string)
	for _, repository := range repositories {
		split := strings.Split(repository, "/")
		if len(split) != 2 {
			return false
		}

		i := strings.Split(split[1], "@")
		if len(i) != 2 {
			return false
		}
	}
	return true
}
