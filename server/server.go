package server

import (
	"context"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	log "github-observer/internal/executor/Logger"
	prometheus "github-observer/internal/executor/Prometheus"
	l "github-observer/internal/listener"
	w "github-observer/internal/watcher"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
)

func Run(configuration Config) {
	repositories := GetRepositories(configuration)
	executors := initExecutors(configuration)
	listener := initListener(repositories, executors)

	if len(configuration.App.Watcher.GithubToken) > 0 {
		initWatcher(repositories, executors, configuration)
	}

	gin.SetMode(configuration.Mode)
	engine := gin.New()
	InitializeRoutes(engine, listener, configuration)

	var err error
	if len(configuration.Ssl.Key) > 0 && len(configuration.Ssl.Cert) > 0 {
		err = engine.RunTLS(configuration.Address, configuration.Ssl.Cert, configuration.Ssl.Key)
	} else {
		err = engine.Run(configuration.Address)
	}
	if err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
	return
}

func initListener(repositories []core.Repository, executors []executor.IExecutor) l.IListener {
	fileLogger := slog.New(slog.NewJSONHandler(&lumberjack.Logger{Filename: ListenerFile, MaxSize: 10, MaxAge: 1}, nil))
	return l.NewListener(repositories, executors, fileLogger)
}

func initExecutors(configuration Config) (executors []executor.IExecutor) {
	for _, e := range configuration.App.Executors {
		switch e {
		case string(ExecutorLogger):
			fileLogger := slog.New(slog.NewJSONHandler(&lumberjack.Logger{Filename: ExecutorFile, MaxSize: 10, MaxAge: 1}, nil))
			executors = append(executors, log.NewExecutor(log.NewMemory(), fileLogger))
		case string(ExecutorPrometheus):
			executors = append(executors, prometheus.NewExecutor())
		}
	}
	return
}

func initWatcher(repositories []core.Repository, executors []executor.IExecutor, configuration Config) {
	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: configuration.App.Watcher.GithubToken})))
	fileLogger := slog.New(slog.NewJSONHandler(&lumberjack.Logger{Filename: WatcherFile, MaxSize: 10, MaxAge: 1}, nil))
	w.NewWatcher(configuration.App.Watcher.GithubToken, client, repositories, executors, fileLogger).Watch()
	return
}

func GetRepositories(configuration Config) []core.Repository {
	var repositories []core.Repository
	for _, repo := range configuration.App.Repositories {
		repositories = append(repositories, core.ToRepository(repo))
	}
	return repositories
}
