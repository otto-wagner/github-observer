# Listener
The listener listens for Actions, Pull Requests and Pull Request review webhooks from GitHub.

## Github Webhook
A webhook is a mechanism that allows an external service to be notified when a specific event occurs.

### Manually
You can also create a webhook manually in the GitHub repository.
- Go to your repository
- Go to Settings
- Go to Webhooks
- Click on Add webhook
    - listen workflow_run
        - Set the Payload URL to your endpoint (e.g. https://{domain}/listen/workflow)
        - content type: application/json
        - Let me select individual events
        - Select Workflow runs
    - listen pullrequest
        - Set the Payload URL to your endpoint (e.g. https://{domain}/listen/pullrequest)
        - content type: application/json
        - Let me select individual events
        - Select Pull requests
    - listen pullrequest review
        - Set the Payload URL to your endpoint (e.g. https://{domain}/listen/pullrequest/review)
        - content type: application/json
        - Let me select individual events
        - Select Pull request reviews

## devtunnel

For local development you can use devtunnel to expose your local Endpoints to the internet.
```shell
devtunnel user login -g # login via github
```
```shell
devtunnel host -p 8443 --allow-anonymous --expiration 4h --protocol https
# open a tunnel to your localhost port 8443. you will get a public URL
```
