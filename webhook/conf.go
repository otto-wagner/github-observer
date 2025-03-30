package webhook

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Webhook struct {
	PayloadUrl  string   `json:"payloadUrl" validate:"required"`
	ContentType string   `json:"contentType" validate:"required"`
	InsecureSsl string   `json:"insecureSsl" validate:"required"`
	Events      []string `json:"events" validate:"required"`
}

type Config struct {
	Repositories []string  `json:"repositories" validate:"required,repositories"`
	HmacSecret   string    `json:"hmacSecret" validate:"required"`
	GithubToken  string    `json:"githubToken" validate:"required"`
	Webhooks     []Webhook `json:"webhook" validate:"required,dive"`
}

type IWebhookConfig interface {
	Validate() error
	Config() Config
}

func InitConfig() (IWebhookConfig, error) {
	config := Config{}
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
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c Config) Config() Config {
	return c
}

func (c Config) Validate() (err error) {
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

func (c Config) validation() error {
	validate := validator.New()
	err := validate.RegisterValidation("repositories", validateRepositories)
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
