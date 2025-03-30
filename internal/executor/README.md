# Executors
Executors are the core of the system. They are responsible for executing the tasks that are assigned to them.
To activate executors you need to set it in the config. (see [conf/readme.md](../../conf/README.md))

### Logging
The logging executor logs the status of the observed events.
- 127.0.0.1:8443/logs/watcher
- 127.0.0.1:8443/logs/executor
- 127.0.0.1:8443/logs/listener

### Prometheus
Prometheus is used to store the metrics of the observer.
- 127.0.0.1:8443/metrics
- 127.0.0.1:9090
