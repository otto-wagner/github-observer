package cmd

import (
	"context"
	"github-observer/internal/config"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	"github-observer/internal/executor/Logging"
	"github-observer/internal/executor/Prometheus"
	l "github-observer/internal/listener"
	w "github-observer/internal/watcher"
	"github-observer/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"os"
)

var (
	cfgFile       string
	configuration config.Config
	engine        *gin.Engine
	executors     []executor.IExecutor
	listener      l.IListener
	rootCmd       = &cobra.Command{
		Use:   "github-observer",
		Short: "github-observer is a simple GitHub observer",
		Long:  "github-observer is a simple GitHub observer.",
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
	cobra.OnInitialize(initConfig, initLogging, initExecutor, initListener, initWatcher, initEngine)
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
			executors = append(executors, Logging.NewExecutor(executor.NewMemory()))
		case "prometheus":
			executors = append(executors, Prometheus.NewExecutor())
		}
	}
	if len(executors) == 0 {
		zap.S().Fatal("no executor")
	}
}

func initListener() {
	listener = l.NewListener(executors)
}

func initWatcher() {
	if viper.GetBool("app.watcher") {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			zap.S().Fatal("no GITHUB_TOKEN")
		}

		var repositoriesConfig []config.RepositoryConfig
		err := viper.UnmarshalKey("app.repositories", &repositoriesConfig)
		if err != nil {
			zap.S().Fatalw("failed to unmarshal app.repositories", "error", err)
		}

		var repositories []core.Repository
		for _, repo := range repositoriesConfig {
			coreRepo := core.Repository{
				Name:    repo.Name,
				Owner:   repo.Owner,
				HtmlUrl: repo.HtmlUrl,
			}
			repositories = append(repositories, coreRepo)
		}

		client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
		w.Watch(client, repositories, executors)
	}
}

func initEngine() {
	gin.SetMode(configuration.App.Mode)
	engine = gin.New()
}
