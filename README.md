# github-listener
A simple GitHub webhook listener that listen for Actions, Pull Requests and Pull Request review events and log them.

## local development
For local development you can use ngrok to expose your local server to the internet.

**Ngrok is only free for non-commercial uses**
* https://ngrok.com/docs/getting-started/
* https://ngrok.com/docs/http/webhook-verification/

```shell
ngrok http --domain={domain} 8443 --verify-webhook=github --verify-webhook-secret=mySecret
```

## Endpoints
- 127.0.0.1:8443/listen/action
- 127.0.0.1:8443/listen/pullrequest
- 127.0.0.1:8443/listen/pullrequest/review
