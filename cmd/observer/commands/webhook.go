package commands

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/otto-wagner/github-observer/internal/webhook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	webhookConf    webhook.Config
	webhookCfgFile string
	cmdWebhook     = &cobra.Command{
		Use:   "webhook",
		Short: "api to manage webhooks",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Usage()
		},
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := viper.Unmarshal(&webhookConf); err != nil {
				slog.Error("failed to unmarshal configuration from viper", "error", err)
				os.Exit(1)
			}

			if webhookCfgFile != "" {
				data, err := os.ReadFile(webhookCfgFile)
				if err != nil {
					slog.Error("failed to read config file", "file", webhookCfgFile, "error", err)
					os.Exit(1)
				}
				if err = json.Unmarshal(data, &webhookConf); err != nil {
					slog.Error("failed to parse config file", "file", webhookCfgFile, "error", err)
					os.Exit(1)
				}
			}

			if err := webhookConf.Validate(); err != nil {
				os.Exit(1)
			}
			return nil
		},
	}
)

func init() {
	_ = godotenv.Load()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cmdWebhook.PersistentFlags().StringSliceVarP(&webhookConf.Repositories, "repositories", "r", webhookConf.Repositories, "owner/repo@branch")
	cmdWebhook.PersistentFlags().StringVarP(&webhookConf.GithubToken, "github_token", "g", webhookConf.GithubToken, "github token")
	cmdWebhook.PersistentFlags().StringVarP(&webhookConf.HmacSecret, "hmac_secret", "s", webhookConf.HmacSecret, "hmac secret")
	cmdWebhook.PersistentFlags().StringVarP(&webhookCfgFile, "file", "f", "", "path to a JSON config file (e.g. example.json)")

	err := viper.BindPFlags(cmdWebhook.PersistentFlags())
	if err != nil {
		slog.Error("failed to bind flags", "error", err)
		os.Exit(1)
	}

	cmdWebhook.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "create webhooks",
		Run: func(_ *cobra.Command, _ []string) {
			webhook.Create(webhookConf)
		},
	})

	cmdWebhook.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list webhooks",
		Run: func(_ *cobra.Command, _ []string) {
			webhook.List(webhookConf)
		},
	})

	cmdWebhook.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "delete webhooks",
		Run: func(_ *cobra.Command, _ []string) {
			webhook.Delete(webhookConf)
		},
	})

	rootCmd.AddCommand(cmdWebhook)
}
