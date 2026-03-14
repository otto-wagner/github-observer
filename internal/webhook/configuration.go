package webhook

import (
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Webhook struct {
	PayloadUrl  string   `json:"payloadUrl" mapstructure:"payload_url" validate:"required"`
	ContentType string   `json:"contentType" mapstructure:"content_type" validate:"required"`
	InsecureSsl string   `json:"insecureSsl" mapstructure:"insecure_ssl" validate:"required"`
	Events      []string `json:"events" mapstructure:"events" validate:"required"`
}

type Config struct {
	Repositories []string  `json:"repositories" mapstructure:"repositories" validate:"required,repositories"`
	GithubToken  string    `json:"github_token" mapstructure:"github_token"`
	HmacSecret   string    `json:"hmac_secret" mapstructure:"hmac_secret"`
	Webhooks     []Webhook `json:"webhooks" mapstructure:"webhook" validate:"required,webhooks"`
}

type IWebhookConfig interface {
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
	err = validate.RegisterValidation("webhooks", c.validateWebhooks)
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

func (c Config) validateWebhooks(fl validator.FieldLevel) bool {
	webhooks := fl.Field().Interface().([]Webhook)
	for _, w := range webhooks {
		if w.PayloadUrl == "" || w.ContentType == "" || w.InsecureSsl == "" || len(w.Events) == 0 {
			return false
		}
	}
	return true
}
