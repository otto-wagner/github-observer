package cmd

import (
	"github-listener/internal/Executor"
	"github-listener/internal/Executor/Logging"
	"github-listener/internal/Executor/Prometheus"
	"github-listener/internal/config"
	"github-listener/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	cfgFile       string
	configuration config.Config
	engine        *gin.Engine
	executors     []Executor.IExecutor
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
	cobra.OnInitialize(initConfig, initLogging, initExecutor, initEngine)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(docCmd)
}

func initConfig() {
	var err error
	configuration, err = config.InitConfig(cfgFile)
	if err != nil {
		zap.S().Fatalw("failed to init config", "error", err)
	}
}

func initLogging() {
	mode := viper.GetString(configuration.App.Mode)
	switch mode {
	case "debug":
		zap.ReplaceGlobals(pkg.NewZapLogger(zapcore.DebugLevel))
	case "production":
		zap.ReplaceGlobals(pkg.NewZapLogger(zapcore.WarnLevel))
	default:
		zap.ReplaceGlobals(pkg.NewZapLogger(zapcore.InfoLevel))
	}
}

func initExecutor() {
	appExecutors := viper.GetStringSlice("app.executors")
	for _, e := range appExecutors {
		switch e {
		case "logging":
			executors = append(executors, Logging.NewExecutor())
		case "prometheus":
			executors = append(executors, Prometheus.NewExecutor())
		}
	}
}

func initEngine() {
	gin.SetMode(configuration.App.Mode)
	engine = gin.New()
}
