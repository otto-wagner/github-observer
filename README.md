# github-observer
Github observer would like to observe github projects and forward the status.

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
