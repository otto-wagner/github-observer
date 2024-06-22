# Listener

## Github Webhook
To use this webhook you need to create a webhook in your GitHub repository.
- Go to your repository
- Go to Settings
- Go to Webhooks
- Click on Add webhook
    - listen action
        - Set the Payload URL to your endpoint (e.g. https://{domain}/listen/action)
        - content type: application/json
        - Let me select individual events
        - Select Check runs
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

## ngrok domain
For local development you can use ngrok to expose your local Endpoints to the internet.
**Ngrok is only free for non-commercial and pre-release versions (https://ngrok.com/pricing)**

+ Register at https://ngrok.com/
+ Installation: https://ngrok.com/docs/getting-started/
+ Set environment variable NGROK_AUTHTOKEN= (https://dashboard.ngrok.com/get-started/your-authtoken)
+ Start ngrok with docker-compose
    * comment ngrok in docker-compose
    * you can use ngrok.yml to set your domain (https://dashboard.ngrok.com/cloud-edge/domains)
* Start ngrok with shell command
  ```shell
  ngrok http --domain={domain} 8443 --verify-webhook=github --verify-webhook-secret=mySecret
  ```
