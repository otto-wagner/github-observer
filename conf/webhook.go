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

type WebHookConfig struct {
	Webhook Webhook `json:"webhook" validate:"required"`
}

type Webhook struct {
	Hooks        []HookConfig       `json:"hooks" validate:"required"`
	Repositories []RepositoryConfig `json:"repositories" validate:"required"`
}

type HookConfig struct {
	PayloadUrl  string   `json:"payloadUrl" validate:"required"`
	ContentType string   `json:"contentType" validate:"required"`
	Secret      string   `json:"secret" validate:"required"`
	InsecureSsl string   `json:"insecureSsl" validate:"required"`
	Events      []string `json:"events" validate:"required"`
}

func InitWebhook(cfgFile string) (c WebHookConfig, err error) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		configPath, err := filepath.Abs("../conf")
		if err != nil {
			slog.Error("failed to get absolute path", "error", err)
			os.Exit(1)
		}
		viper.AddConfigPath(configPath)
		viper.SetConfigName("webhook")
		viper.SetConfigType("json")
	}
	// load env variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var err2 viper.ConfigFileNotFoundError
		if !errors.As(err, &err2) {
			return c, err2
		}
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}
	return
}

func (c WebHookConfig) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
