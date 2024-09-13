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

type IWebhookConfig interface {
	Validate() error
	Config() WebHookConfig
}

type webhookConfig struct {
	webHookConfig WebHookConfig
}

type WebHookConfig struct {
	HmacSecret   string          `json:"hmacSecret" validate:"required"`
	Webhooks     []WebhookConfig `json:"webhooks" validate:"required,dive"`
	Repositories []string        `json:"repositories" validate:"required,repositories"`
}

type WebhookConfig struct {
	PayloadUrl  string   `json:"payloadUrl" validate:"required"`
	ContentType string   `json:"contentType" validate:"required"`
	InsecureSsl string   `json:"insecureSsl" validate:"required"`
	Events      []string `json:"events" validate:"required"`
}

func InitWebhook() (IWebhookConfig, error) {
	webhookConfig := webhookConfig{}
	configPath, err := filepath.Abs("./conf")
	if err != nil {
		slog.Error("failed to get absolute path", "error", err)
		os.Exit(1)
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("webhook")
	viper.SetConfigType("json")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}
	err = viper.Unmarshal(&webhookConfig.webHookConfig)
	if err != nil {
		return nil, err
	}
	return &webhookConfig, nil
}

func (c webhookConfig) Config() WebHookConfig {
	return c.webHookConfig
}

func (c webhookConfig) Validate() (err error) {
	err = c.validation()
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
	return
}

func (c webhookConfig) validation() error {
	validate := validator.New()
	err := validate.RegisterValidation("repositories", validateRepositories)
	if err != nil {
		return err
	}
	return validate.Struct(c.webHookConfig)
}
