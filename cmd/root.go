package cmd

import (
	"github-listener/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var (
	cfgFile       string
	configuration config.Config
	engine        *gin.Engine
	rootCmd       = &cobra.Command{
		Use:   "github-listener",
		Short: "github-listener is a simple GitHub webhook listener",
		Long:  "github-listener is a simple GitHub webhook listener.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initEngine)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(docCmd)
}

func initEngine() {
	gin.SetMode(configuration.App.Mode)
	engine = gin.New()
}

func initConfig() {
	var err error
	configuration, err = config.InitConfig(cfgFile)
	if err != nil {
		zap.S().Fatalw("failed to init config", "error", err)
	}
}
