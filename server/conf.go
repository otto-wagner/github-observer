package server

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const (
	ExecutorLogger     Executor = "logging"
	ExecutorPrometheus Executor = "prometheus"
	ObserverFile       string   = "./var/log/observer.log"
	WatcherFile        string   = "./var/log/watcher.log"
	ListenerFile       string   = "./var/log/listener.log"
	ExecutorFile       string   = "./var/log/executor.log"
)

type Executor string

type Config struct {
	Mode           string   `json:"mode" validate:"required,oneof=release debug test"`
	Address        string   `json:"address" validate:"required"`
	TrustedProxies []string `json:"trustedProxies" validate:"omitempty,dive,ip"`
	Ssl            Ssl      `json:"ssl"`
	App            App      `json:"app" validate:"required"`
}

type Ssl struct {
	Cert string `json:"cert" validate:"required_with=Key"`
	Key  string `json:"key" validate:"required_with=Cert"`
}

type App struct {
	Repositories []string `json:"repositories" validate:"required,repositories"`
	Executors    []string `json:"executors" validate:"required,executors"`
	Watcher      Watcher  `json:"watcher" validate:"watcher"`
	Listener     Listener `json:"listener" validate:"listener"`
	Logger       []string `json:"logger" validate:"logger"`
}

type Watcher struct {
	Enabled     bool   `json:"enabled"`
	GithubToken string `json:"githubToken"`
}

type Listener struct {
	Enabled    bool   `json:"enabled"`
	HmacSecret string `json:"hmacSecret"`
}

type Repository struct {
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Branch string `json:"branch"`
}

type IServerConfig interface {
	Validate() error
	Config() Config
}

func InitConfig() (IServerConfig, error) {
	config := Config{}
	configPath, err := filepath.Abs("./conf")
	if err != nil {
		slog.Error("failed to get absolute path", "error", err)
		os.Exit(1)
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("server")
	viper.SetConfigType("json")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
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
	err = validate.RegisterValidation("logger", c.validateLogger)
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
		"logging":    true,
		"prometheus": true,
	}
	for _, executor := range executors {
		if !allowed[executor] {
			return false
		}
		if executor == "logging" {
			if len(c.App.Logger) == 0 {
				return false
			}
			for _, logger := range c.App.Logger {
				if logger != "listener" && logger != "executor" && logger != "watcher" {
					return false
				}
			}
		}
	}
	return true
}

func (c Config) validateLogger(fl validator.FieldLevel) bool {
	loggers := fl.Field().Interface().([]string)
	for _, logger := range loggers {
		if logger != "listener" && logger != "executor" && logger != "watcher" {
			return false
		}
	}
	executors := c.App.Executors
	for _, executor := range executors {
		if executor == "logging" {
			return true
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
