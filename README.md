# github-observer
GitHub observer would like to observe GitHub projects and forward the status.

# Setup Application
- You need to create a GitHub token and add it to the environment variable.
- Set necessary configurations in the `.env` (see [.env.example](.env.example)).

# Start Application
```bash
    docker run -d --name github-observer ghcr.io/otto-wagner/github-observer:latest
```

# Listener
A simple GitHub webhook listener that listen for Actions, Pull Requests and Pull Request review events (see [internal/listener/readme.md](internal/listener/README.md)).

# Watcher
The watcher use the GitHub API to get the status of Actions, Pull Requests and Pull Request reviews (see [internal/watcher/readme.md](internal/watcher/README.md)).

# Executors
We want to execute actions based on the status of the observed events. (see [internal/executor/readme.md](internal/executor/README.md))
