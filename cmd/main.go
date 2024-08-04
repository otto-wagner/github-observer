package main

import (
	server "github-observer/server/cmd"
	webhook "github-observer/webhook/cmd"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func init() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)
}

func main() {
	if err := rootCommand().Execute(); err != nil {
		slog.Error("failed to execute root command", "error", err)
		os.Exit(1)
	}
}

func rootCommand() (rootCmd *cobra.Command) {
	rootCmd = &cobra.Command{
		Use:   "observer",
		Short: "GitHub Observer",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Usage()
			if err != nil {
				return
			}
		},
	}
	rootCmd.AddCommand(server.Server())
	rootCmd.AddCommand(webhook.WebHook())

	return
}
