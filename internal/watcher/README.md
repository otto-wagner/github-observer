# Watcher
A watcher process is implemented that regularly retrieves data from GitHub.
The retrieved data is then forwarded to an executor, which processes it and performs the corresponding actions.

## Github API
GitHub provides a REST API that allows access to various data.
The API can be used to retrieve repositories, pull requests, issues, workflow runs, and much more.
Link to the [GitHub API](https://docs.github.com/en/rest/reference).

## Pull Requests
All open pull requests are queried and forwarded to the executor.

## Workflow Runs
An important part of this process is the querying and processing of workflow runs of a specific repository.
1. First, all workflows of the repositories are determined.
    - The GitHub API is used to retrieve all workflows of a repository.
        - `GET /repos/{owner}/{repo}/actions/workflows`: Lists all workflows of a repository.
2. Now, the last completed ("completed") run of the workflow from the `main` or `master` branch is queried.
    - The GitHub API is used to retrieve all workflow runs of a repository.
        - `GET /repos/{owner}/{repo}/actions/runs`: Lists all workflow runs of a repository.
3. Finally, the determined workflow runs are forwarded to the executor.

## Github Rate Limit
The GitHub API has a rate limit. This means that only a certain number of requests can be made within a certain period of time.
If the rate limit is exceeded, the API will return an error message.
To avoid this, the watcher process is implemented in such a way that it waits for a certain period of time before making the next request.
