package conf

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Executor string

const (
	Logging      Executor = "logging"
	Prometheus   Executor = "prometheus"
	WatcherFile  string   = "watcher.log"
	ListenerFile string   = "listener.log"
	ExecutorFile string   = "executor.log"
)

type AppConfig struct {
	ListenAddress  string   `json:"listenAddress"`
	TrustedProxies []string `json:"trustedProxies"`
	Mode           string   `json:"mode" validate:"omitempty,oneof=release debug test"`
	Executors      []string `json:"executors" validate:"required,executors"`
	Watcher        bool     `json:"watcher"`
	Logs           []string `json:"logs"`
}

type SslConfig struct {
	Activate bool   `json:"activate"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

type CommonConfig struct {
	App          AppConfig `json:"app"`
	Repositories []string  `json:"repositories" validate:"required,repositories"`
	Ssl          SslConfig `json:"ssl"`
	HmacSecret   string    `json:"hmacSecret" validate:"required"`
}

type ICommonConfig interface {
	Validate() error
	Config() CommonConfig
}

type commonConfig struct {
	commonConfig CommonConfig
}

func InitCommon() (ICommonConfig, error) {
	commonConfig := commonConfig{}
	configPath, err := filepath.Abs("./conf")
	if err != nil {
		slog.Error("failed to get absolute path", "error", err)
		os.Exit(1)
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("common")
	viper.SetConfigType("json")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	err = viper.Unmarshal(&commonConfig.commonConfig)
	if err != nil {
		return nil, err
	}
	return &commonConfig, nil
}

func (c commonConfig) Config() CommonConfig {
	return c.commonConfig
}

func (c commonConfig) Validate() error {
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
		slog.Error("configuration validation failed")
	}
	return err
}

func (c commonConfig) validation() error {
	validate := validator.New()
	err := validate.RegisterValidation("executors", c.validateExecutors)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("repositories", validateRepositories)
	if err != nil {
		return err
	}
	return validate.Struct(c)
}

func (c commonConfig) validateExecutors(fl validator.FieldLevel) bool {
	executors := fl.Field().Interface().([]string)
	allowed := map[string]bool{
		"logging":    true,
		"prometheus": true,
	}
	for _, executor := range executors {
		if !allowed[executor] {
			return false
		}
	}
	return true
}
