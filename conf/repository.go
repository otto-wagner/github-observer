package conf

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type Repository struct {
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Branch string `json:"branch"`
}

func validateRepositories(fl validator.FieldLevel) bool {
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
