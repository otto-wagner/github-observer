package core

import (
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Token struct {
	Key    string
	Expiry time.Duration
}

func (t *Token) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: viper.GetString(t.Key), Expiry: time.Now().Add(t.Expiry)}, nil
}
