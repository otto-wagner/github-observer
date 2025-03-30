# configurations
This directory contains the configuration files for the project

## server
The server configuration in the `server.json` file.

### mode
The mode of the server. The following modes are supported:
- `debug`: Debug mode. This mode is used for development and testing.
- `release`: Release mode. This mode is used for production.
- `test`: Test mode. This mode is used for testing.

### address
The address of the server i.e. `localhost:8443`.

### trustedProxies
The trusted proxies of the server.

### ssl
The SSL configuration
- `ssl.activate`: Activate SSL (1 = activate, 0 = deactivate).
- `ssl.certFile`: The path to the certificate file.
- `ssl.keyFile`: The path to the key file.

### app
The app configurations.

#### repositories
The repositories to observe e.g. `otto-wagner/github-observer@main`.
- `owner:` The owner of the repository.
- `name:` The name of the repository.
- `branch:` The observed branch of the repository.

#### executors
Enable the executors to execute actions based on the status of the observed events.
- `logging`
- `prometheus`

### watcher
The watcher watches the status of Actions, Pull Requests and Pull Request reviews.

#### enabled
Enable the watcher.

#### githubToken
The GitHub token to access the GitHub API. This token is used to authenticate the watcher with the GitHub API.

### listener
The listener listens to the events from the GitHub API.
#### enabled
Enable the listener.

#### hmacSecret
The HMAC secret to protect the server from unauthorized access. This secret is used to verify the authenticity of the events sent by the GitHub webhooks.

### logger
Activate the log endpoints of the server. 
- `executor`
- `watcher`
- `listener`

## webhook
The webhook configuration in the `webhook.json` file.

### hmacSecret
HMAC secret the webhook use to authorize the requests to the server.

### githubToken
The GitHub token to access the GitHub API. This token is used configure the webhooks.

### webhooks
The webhooks to create.

#### payloadUrl
The endpoint of the webhook. Currently, only `listen/workflow` and `listen/pullrequest`, `listen/pullrequest/review` are supported.

#### contentType
The content type of the webhook. Currently, only `json` is supported.

#### secret
The secret of the webhook. This secret is used to verify the authenticity of the events sent by the GitHub webhooks.

#### insecureSsl
The status of the SSL (1 = insecure, 0 = secure). This option is used to disable SSL verification for the webhook.


#### events
The events that should be listened to. Currently, only the following events are supported:
- `workflow_run`
- `pull_request`
- `pull_request_review`

### repositories
The repositories to observe e.g. `otto-wagner/github-observer@main`.
- `owner:` The owner of the repository.
- `name:` The name of the repository.
- `branch:` The observed branch of the repository.
