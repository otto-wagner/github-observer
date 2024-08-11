# configurations
This directory contains the configuration files for the project

## secret
The secret configuration for hmac authentication.

## ssl
The SSL configuration in the `common.json` file.
- `ssl.activate`: Activate SSL (1 = activate, 0 = deactivate).
- `ssl.certFile`: The path to the certificate file.
- `ssl.keyFile`: The path to the key file.

## app
The app configuration in the `common.json` file.

### logs
Enable log endpoints.
- `executor`
- `watcher`
- `listener`

### watcher
Enable the watcher to get the status of Actions, Pull Requests and Pull Request reviews.

### executors
Enable the executors to execute actions based on the status of the observed events.
- `logging`
- `prometheus`

### repositories
The repositories to observe.
- `owner:` The owner of the repository.
- `name:` The name of the repository.
- `branch:` The observed branch of the repository.

## webhook
The webhook configuration in the `webhook.json` file.
- `webhooks.payloadUrl`: The endpoint of the webhook. Currently, only `listen/workflow` and `listen/pullrequest`, `listen/pullrequest/review` are supported.
- `webhooks.contentType`: The content type of the webhook. Currently, only `json` is supported.
- `webhooks.secret`: The secret of the webhook.
- `webhooks.insecureSsl`: The status of the SSL (1 = insecure, 0 = secure).
- `webhooks.events`: The events that should be listened to. Currently, only the following events are supported:
    - `workflow_run`
    - `pull_request`
    - `pull_request_review`
- `repositories.owner`: The owner of the repository.
- `repositories.name`: The name of the repository.
