package main

import (
	"github-observer/server"
	"github-observer/webhook"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func init() {
	file, err := os.OpenFile(server.ObserverFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		os.Exit(1)
	}

	jsonHandler := slog.NewJSONHandler(file, nil)
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
	rootCmd.AddCommand(server.Cmd())
	rootCmd.AddCommand(webhook.Cmd())

	return
}
