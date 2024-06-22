# configurations
This directory contains the configuration files for the project

## app

### watcher
Enable the watcher to get the status of Actions, Pull Requests and Pull Request reviews.

### executors
Enable the executors to execute actions based on the status of the observed events.
- **logging**
- **prometheus**l

### repositories
The repositories to observe.
- **owner:** The owner of the repository.
- **name:** The name of the repository.
- **branch:** The observed branch of the repository.
