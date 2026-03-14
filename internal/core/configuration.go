package core

import (
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
)

const ExecutorMetrics Executor = "metrics"

type Executor string

type Config struct {
	Mode         string   `mapstructure:"mode" validate:"required,oneof=release debug test"`
	Address      string   `mapstructure:"address" validate:"required"`
	Repositories []string `mapstructure:"repositories" validate:"required,repositories"`
	Executors    []string `mapstructure:"executors" validate:"required,executors"`
	Watcher      Watcher  `mapstructure:"watcher" validate:"watcher"`
	Listener     Listener `mapstructure:"listener" validate:"listener"`
}

type Watcher struct {
	Enabled     bool   `mapstructure:"enabled"`
	GithubToken string `mapstructure:"github_token"`
}

type Listener struct {
	Enabled    bool   `mapstructure:"enabled"`
	HmacSecret string `mapstructure:"hmac_secret"`
}

type Repo struct {
	Name   string `mapstructure:"name"`
	Owner  string `mapstructure:"owner"`
	Branch string `mapstructure:"branch"`
}

type IServerConfig interface {
	Validate() error
	Config() Config
}

func (c Config) Config() Config {
	return c
}

func (c Config) Validate() error {
	err := c.validation()
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
	}
	return err
}

func (c Config) validation() error {
	validate := validator.New()
	err := validate.RegisterValidation("repositories", c.validateRepositories)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("executors", c.validateExecutors)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("listener", c.validateListener)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("watcher", c.validateWatcher)
	if err != nil {
		return err
	}
	return validate.Struct(c)
}

func (c Config) validateRepositories(fl validator.FieldLevel) bool {
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

func (c Config) validateExecutors(fl validator.FieldLevel) bool {
	executors := fl.Field().Interface().([]string)
	allowed := map[string]bool{
		string(ExecutorMetrics): true,
	}
	for _, executor := range executors {
		if !allowed[executor] {
			return false
		}
	}
	return true
}

func (c Config) validateListener(fl validator.FieldLevel) bool {
	listener := fl.Field().Interface().(Listener)
	if listener.Enabled && len(listener.HmacSecret) == 0 {
		return false
	}
	return true
}

func (c Config) validateWatcher(fl validator.FieldLevel) bool {
	watcher := fl.Field().Interface().(Watcher)
	if watcher.Enabled && len(watcher.GithubToken) == 0 {
		return false
	}
	return true
}
