package Prometheus

import (
	"github-observer/internal/core"
	e "github-observer/internal/executor"
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type executor struct {
	eventRun         *prometheus.CounterVec
	eventPullRequest *prometheus.CounterVec
	workflowRun      *prometheus.GaugeVec
	pullRequest      *prometheus.GaugeVec
}

func NewExecutor() e.IExecutor {
	eventRun := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_run", Help: "Number of action requests processed"},
		[]string{"action", "check_run_id", "check_run_name", "check_run_status", "check_run_conclusion", "repository_full_name"})

	eventPullRequest := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_pull_request", Help: "Number of pull requests processed"},
		[]string{"action", "pull_request_title", "pull_request_state", "repository_full_name"})

	workflowRun := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "latest_workflows", Help: "Failed workflow runs"},
		[]string{"repository_full_name", "workflow_name", "state", "conclusion"})

	pullRequest := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pull_request_sum", Help: "Number of pull requests"},
		[]string{"repository_full_name"})

	prometheus.MustRegister(eventRun, eventPullRequest, workflowRun, pullRequest)
	return &executor{eventRun, eventPullRequest, workflowRun, pullRequest}
}

func (e *executor) EventRun(event github.CheckRunEvent) {
	action := core.ConvertToGitAction(event)
	e.eventRun.With(map[string]string{
		"action":               action.Action,
		"check_run_id":         strconv.FormatInt(action.CheckRun.Id, 10),
		"check_run_name":       action.CheckRun.Name,
		"check_run_status":     action.CheckRun.Status,
		"check_run_conclusion": action.CheckRun.Conclusion,
		"repository_full_name": action.Repository.FullName,
	}).Inc()
	return
}

func (e *executor) EventPullRequest(event github.PullRequestEvent) {
	pr := core.ConvertPREToGitPullRequest(event)
	e.eventPullRequest.With(map[string]string{
		"action":               pr.Action,
		"pull_request_title":   pr.PullRequest.Title,
		"pull_request_state":   pr.PullRequest.State,
		"repository_full_name": pr.Repository.FullName,
	}).Inc()
	return
}

func (e *executor) EventPullRequestReview(github.PullRequestReviewEvent) {
	// ignored
	return
}

func (e *executor) LastWorkflows(repository core.Repository, workflows []*github.WorkflowRun) {
	e.workflowRun.DeletePartialMatch(map[string]string{"repository_full_name": repository.Owner + "/" + repository.Name})
	for _, run := range workflows {
		flow := core.ConvertToWorkflow(*run)
		if flow.Conclusion != "success" {
			e.workflowRun.With(map[string]string{
				"repository_full_name": flow.Repository.FullName,
				"workflow_name":        flow.Name,
				"state":                flow.Status,
				"conclusion":           flow.Conclusion,
			}).Set(1)
		}
	}
}

func (e *executor) PullRequests(repository core.Repository, pullRequests []*github.PullRequest) {
	e.pullRequest.DeletePartialMatch(map[string]string{"repository_full_name": repository.Owner + "/" + repository.Name})
	if len(pullRequests) == 0 {
		return
	}
	var count int
	for _, pullRequest := range pullRequests {
		if pullRequest.GetState() != "closed" && pullRequest.GetState() != "merged" {
			count++
		}
	}
	e.pullRequest.With(map[string]string{"repository_full_name": core.ConvertPRToGitPullRequest(*pullRequests[0]).Repository.FullName}).Set(float64(count))
	return
}
