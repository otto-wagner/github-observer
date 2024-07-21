package server

import (
	"context"
	"github-observer/conf"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	"github-observer/internal/executor/Logging"
	"github-observer/internal/executor/Prometheus"
	l "github-observer/internal/listener"
	"github-observer/internal/router"
	w "github-observer/internal/watcher"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
)

func Run(configuration conf.Config) {
	validationConfiguration(configuration)

	repositories := GetRepositories(configuration)
	executors := initExecutors(configuration)
	listener := initListener(repositories, executors)

	if configuration.App.Watcher {
		initWatcher(repositories, executors)
	}

	gin.SetMode(configuration.App.Mode)
	engine := gin.New()
	router.InitializeRoutes(engine, listener, configuration.App.Watcher)

	err := engine.Run(configuration.App.ListenAddress)
	if err != nil {
		slog.Error("failed to start httpServer", "error", err)
		os.Exit(1)
	}
	return
}

func initListener(repositories []core.Repository, executors []executor.IExecutor) l.IListener {
	return l.NewListener(repositories, executors)
}

func initExecutors(configuration conf.Config) (executors []executor.IExecutor) {
	for _, e := range configuration.App.Executors {
		switch e {
		case conf.Logging:
			fileLogger := slog.New(slog.NewJSONHandler(&lumberjack.Logger{Filename: "executor.log", MaxSize: 10, MaxAge: 1}, nil))
			executors = append(executors, Logging.NewExecutor(Logging.NewMemory(), fileLogger))
		case conf.Prometheus:
			executors = append(executors, Prometheus.NewExecutor())
		}
	}
	return
}

func initWatcher(repositories []core.Repository, executors []executor.IExecutor) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		slog.Error("missing GITHUB_TOKEN")
		os.Exit(1)
	}

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	fileLogger := slog.New(slog.NewJSONHandler(&lumberjack.Logger{Filename: conf.WatcherFile, MaxSize: 10, MaxAge: 1}, nil))
	w.NewWatcher(token, client, repositories, executors, fileLogger).Watch()
	return
}

func GetRepositories(configuration conf.Config) []core.Repository {
	var repositories []core.Repository
	for _, repo := range configuration.App.Repositories {
		coreRepo := core.Repository{
			Name:   repo.Name,
			Owner:  repo.Owner,
			Branch: repo.Branch,
		}
		repositories = append(repositories, coreRepo)
	}
	return repositories
}

func validationConfiguration(configuration conf.Config) {
	err := configuration.Validate()
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
	}

	slog.Info("start server", "mode", configuration.App.Mode, "listenAddress", configuration.App.ListenAddress,
		"trustedProxies", configuration.App.TrustedProxies, "executors", configuration.App.Executors,
		"watcher", configuration.App.Watcher, "repositories", configuration.App.Repositories)
}
