# Watcher
Ein Watcher-Prozess wird implementiert, der regelmäßig Daten von GitHub abruft. 
Die abgerufenen Daten werden dann an einen Executor weitergeleitet, der diese verarbeitet und die entsprechenden Aktionen ausführt.

## Github API
GitHub stellt eine REST-API zur Verfügung, die es ermöglicht, auf verschiedene Daten zuzugreifen.
Die API kann verwendet werden, um Repositories, Pull Requests, Issues, Workflow-Runs und vieles mehr abzurufen.
Link zur [GitHub API](https://docs.github.com/en/rest/reference).

## Pull Requests
Es werden alle noch offenen Pull Requests abgefragt und an die Executor weitergeleitet.

## Workflows Runs
Ein wichtiger Teil dieses Prozesses ist die Abfrage und Verarbeitung von Workflow-Runs eines bestimmten Repositories.
1. Zunächst werden alle Workflows der Repositories ermittelt. 
   - Dazu wird die GitHub API verwendet, um alle Workflows eines Repositories abzurufen.
     - `GET /repos/{owner}/{repo}/actions/workflows`: Listet alle Workflows eines Repositories auf.
2. Nun wird der letzte abgeschlossene ("completed") Run des Workflows vom `main` oder `master` Branch abgefragt.
   - Dazu wird die GitHub API verwendet, um alle Workflow-Runs eines Repositories abzurufen.
     - `GET /repos/{owner}/{repo}/actions/runs`: Listet alle Workflow-Runs eines Repositories auf.
3. Abschließend werden die ermittelten Workflow-Runs an die Executor weitergeleitet.
