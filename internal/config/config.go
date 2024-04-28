package config

import (
	"errors"
	"go.uber.org/zap"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type AppConfig struct {
	ListenAddress  string             `json:"listenAddress" validate:"hostname_port"`
	TrustedProxies []string           `json:"trustedProxies"`
	Mode           string             `json:"mode" validate:"omitempty,oneof=release debug test"`
	Executors      []string           `json:"executors" validate:"omitempty"`
	Watcher        bool               `json:"watcher"`
	Repositories   []RepositoryConfig `json:"repositories"`
}

type RepositoryConfig struct {
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	HtmlUrl string `json:"htmlUrl"`
}

type SslConfig struct {
	Cert string `json:"cert" validate:"required,file"`
	Key  string `json:"key" validate:"required,file"`
}

type Config struct {
	App AppConfig `json:"app"`
	//Ssl SslConfig `json:"ssl"`
}

func InitConfig(cfgFile string) (Config, error) {
	var c Config

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in default location conf/common.json
		configPath, err := filepath.Abs("../conf")
		if err != nil {
			zap.S().Fatalw("failed to load config - filepath", "error", err)
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

	err := viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
