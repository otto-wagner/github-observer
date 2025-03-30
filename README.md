# github-observer
Github observer would like to observe github projects and forward the status.

# Setup Application
- You need to create a GitHub token and add it to the environment variable.
- Set necessary configurations in the `conf` directory (see [conf/readme.md](conf/README.md)).

# Start Application
You can start the application with the necessary configurations in the `conf` directory with the following command:
```bash
    docker run -d --name github-observer \
      -v "$(pwd)/conf/server.json:/conf/server.json" \
      -p 8443:8443 \
      ghcr.io/otto-wagner/github-observer:latest
```

# Listener
A simple GitHub webhook listener that listen for Actions, Pull Requests and Pull Request review events (see [internal/listener/readme.md](internal/listener/README.md)).

# Watcher
The watcher use the GitHub API to get the status of Actions, Pull Requests and Pull Request reviews (see [internal/watcher/readme.md](internal/watcher/README.md)).
To activate the watcher you need to add a GitHub token to the environment variable and set it in the config (see [conf/readme.md](conf/README.md)).

# Executors
We want to execute actions based on the status of the observed events. (see [internal/executor/readme.md](internal/executor/README.md))

# Grafana
Grafana is used to visualize the metrics of the observer with Prometheus.
- 127.0.0.1:3000
