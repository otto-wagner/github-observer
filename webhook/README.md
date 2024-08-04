# Webhook
Webhooks send a POST request to a specified URL when something happens in a repository.
You can find webhooks in your repository settings under `Webhooks`.

## Description
This cli tool is used to create, list and delete webhooks for a GitHub repository.
see [GitHub Webhook Events](https://docs.github.com/en/webhooks/webhook-events-and-payloads)

## Configuration
Before you can use the webhook cli, you need to configure the webhook.json file.
see [conf/readme.md](../conf/README.md) 

## Commands
The following commands are available:

- `create`: Create all webhooks for all configured webhooks for all configured repositories.
```shell
go run ./cmd/main.go webhook create
```

- `list`: List all webhooks of all configured repositories.
```shell
go run ./cmd/main.go webhook list
```

- `delete`: Delete all webhooks of all configured repositories.
```shell
go run ./cmd/main.go webhook delete
```