package conf

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
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
	ListenAddress  string     `json:"listenAddress" validate:"hostname_port"`
	TrustedProxies []string   `json:"trustedProxies"`
	Mode           string     `json:"mode" validate:"omitempty,oneof=release debug test"`
	Executors      []Executor `json:"executors"`
	Watcher        bool       `json:"watcher"`
	Logs           []string   `json:"logs"`
}

type SslConfig struct {
	Activate bool   `json:"activate"`
	Cert     string `json:"cert" validate:"required,file"`
	Key      string `json:"key" validate:"required,file"`
}

type Config struct {
	App          AppConfig          `json:"app"`
	Repositories []RepositoryConfig `json:"repositories" validate:"required"`
	Ssl          SslConfig          `json:"ssl"`
	Secret       string             `json:"secret" validate:"required"`
}

func InitCommon(cfgFile string) (c Config, err error) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in default location conf/common.json
		configPath, err := filepath.Abs("../conf")
		if err != nil {
			slog.Error("failed to get absolute path", "error", err)
			os.Exit(1)
		}
		viper.AddConfigPath(configPath)
		viper.SetConfigName("common")
		viper.SetConfigType("json")
	}
	// load env variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var err2 viper.ConfigFileNotFoundError
		if !errors.As(err, &err2) {
			return c, err
		}
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}
	return
}

func (c Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
