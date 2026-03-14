package commands

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v73/github"
	"github.com/joho/godotenv"
	"github.com/otto-wagner/github-observer/internal/core"
	"github.com/otto-wagner/github-observer/internal/executor"
	"github.com/otto-wagner/github-observer/internal/executor/metrics"
	"github.com/otto-wagner/github-observer/internal/listener"
	"github.com/otto-wagner/github-observer/internal/utils"
	"github.com/otto-wagner/github-observer/internal/watcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	conf      core.Config
	cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Starts the HTTP server",
		Run: func(_ *cobra.Command, _ []string) {
			gin.SetMode(conf.Mode)

			root := gin.New()
			root.Use(gin.Recovery())
			root.Use(utils.Logger())
			root.Use(cors.New(cors.Config{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
				AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Authorization"},
			}))
			root.GET("/health")

			var repositories []core.Repository
			for _, repo := range conf.Repositories {
				repositories = append(repositories, core.ToRepository(repo))
			}

			var executors []executor.IExecutor
			for _, e := range conf.Executors {
				switch e {
				case string(core.ExecutorMetrics):
					executor := prometheus.NewExecutor()
					executors = append(executors, executor)
					root.GET("/metrics", executor.Handler())
				}
			}

			if conf.Listener.Enabled {
				listener := listener.NewListener(repositories, executors)
				el := root.Group("/listen")
				el.Use(utils.Hmac([]byte(conf.Listener.HmacSecret)))
				el.POST("/workflow", listener.Workflow)
				el.POST("/pullrequest", listener.PullRequest)
				el.POST("/pullrequest/review", listener.PullRequestReview)
			}

			if conf.Watcher.Enabled {
				client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.ReuseTokenSource(nil, &core.Token{
					Key:    "watcher.github_token",
					Expiry: 5 * time.Minute,
				})))
				watcher.NewWatcher(client, repositories, executors).Start()
			}
			err := root.Run(conf.Address)
			if err != nil {
				slog.Error("failed to start server")
				os.Exit(1)
			}
		},
	}
)

func init() {
	_ = godotenv.Load()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cmdServer.PersistentFlags().StringVarP(&conf.Mode, "mode", "m", conf.Mode, "Run mode: debug, release or test")
	cmdServer.PersistentFlags().StringVarP(&conf.Address, "address", "a", conf.Address, "Address of the http api endpoint")

	// => executors
	cmdServer.PersistentFlags().StringSliceVarP(&conf.Executors, "executors", "e", conf.Executors, "Executors to use (logging, prometheus)")

	// => watcher
	cmdServer.PersistentFlags().BoolVarP(&conf.Watcher.Enabled, "watcher.enabled", "w", conf.Watcher.Enabled, "enable watcher")
	cmdServer.PersistentFlags().StringVarP(&conf.Watcher.GithubToken, "watcher.github_token", "g", conf.Watcher.GithubToken, "github token")

	// => listener
	cmdServer.PersistentFlags().BoolVarP(&conf.Listener.Enabled, "listener.enabled", "l", conf.Listener.Enabled, "enable listener")
	cmdServer.PersistentFlags().StringVarP(&conf.Listener.HmacSecret, "listener.hmac_secret", "s", conf.Listener.HmacSecret, "hmac secret")

	// => repositories
	cmdServer.PersistentFlags().StringSliceVarP(&conf.Repositories, "repositories", "r", conf.Repositories, "owner/repo@branch")

	err := viper.BindPFlags(cmdServer.PersistentFlags())
	if err != nil {
		slog.Error("failed to bind flags", "error", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		slog.Error("failed to unmarshal configuration from viper", "error", err)
		os.Exit(1)
	}

	if err := conf.Validate(); err != nil {
		os.Exit(1)
	}

	rootCmd.AddCommand(cmdServer)
}
