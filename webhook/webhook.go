package webhook

import (
	"context"
	"github-observer/conf"
	"github-observer/internal/core"
	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"os"
)

func Create(configuration conf.WebHookConfig) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		slog.Error("missing GITHUB_TOKEN")
		os.Exit(1)
	}

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	for _, hook := range configuration.Webhooks {
		for _, repo := range configuration.Repositories {
			repository := core.ToRepository(repo)
			hookResponse, response, err := client.Repositories.CreateHook(context.Background(), repository.Owner, repository.Name, &github.Hook{
				Config: &github.HookConfig{
					URL:         github.String(hook.PayloadUrl),
					ContentType: github.String(hook.ContentType),
					Secret:      github.String(configuration.HmacSecret),
					InsecureSSL: github.String(hook.InsecureSsl),
				},
				Events: hook.Events,
				Active: github.Bool(true),
			})
			if response.Response.StatusCode == http.StatusUnprocessableEntity {
				slog.Info("hook already exists", "url", hook.PayloadUrl, "message", err)
				continue
			}
			if err != nil {
				slog.Error("failed to create webhook", "error", err)
				continue
			}
			slog.Info("webhook created", "url", hookResponse.GetURL(), "events", hookResponse.Events, "active", hookResponse.GetActive())
		}
	}

	return
}

func List(configuration conf.WebHookConfig) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		slog.Error("missing GITHUB_TOKEN")
		os.Exit(1)
	}

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	for _, repo := range configuration.Repositories {
		repository := core.ToRepository(repo)
		hooks, _, err := client.Repositories.ListHooks(context.Background(), repository.Owner, repository.Name, nil)
		if err != nil {
			slog.Error("failed to list webhooks", "error", err)
			os.Exit(1)
		}
		for _, hook := range hooks {
			slog.Info("webhook", "url", hook.GetURL(), "events", hook.Events, "active", hook.GetActive())
		}
	}
	slog.Info("webhooks listed")
	return
}

func Delete(configuration conf.WebHookConfig) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		slog.Error("missing GITHUB_TOKEN")
		os.Exit(1)
	}

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	for _, repo := range configuration.Repositories {
		repository := core.ToRepository(repo)
		hooks, _, err := client.Repositories.ListHooks(context.Background(), repository.Owner, repository.Name, nil)
		if err != nil {
			slog.Error("failed to list webhooks", "error", err)
			os.Exit(1)
		}
		for _, hook := range hooks {
			_, err := client.Repositories.DeleteHook(context.Background(), repository.Owner, repository.Name, *hook.ID)
			if err != nil {
				slog.Error("failed to delete webhook", "error", err)
				continue
			}
			slog.Info("webhook deleted", "url", hook.GetURL(), "events", hook.Events, "active", hook.GetActive())
		}
	}
	slog.Info("webhooks deleted")
	return
}
