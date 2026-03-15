package prometheus

import (
	"strconv"

	"github.com/otto-wagner/github-observer/internal/core"
	e "github.com/otto-wagner/github-observer/internal/executor"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v73/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type executor struct {
	registry         *prometheus.Registry
	eventPullRequest *prometheus.CounterVec
	workflowRun      *prometheus.GaugeVec
	pullRequest      *prometheus.GaugeVec
}

func NewExecutor() e.IExecutor {
	registry := prometheus.NewRegistry()
	eventPullRequest := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "observer_event_pull_request", Help: "Number of pull requests processed"},
		[]string{"action", "pull_request_title", "pull_request_state", "repository_full_name"})

	workflowRun := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "observer_latest_workflows", Help: "Failed workflow runs"},
		[]string{"repository_full_name", "workflow_name", "workflow_run_id", "run_number", "state", "conclusion"})

	pullRequest := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "observer_pull_request_sum", Help: "Number of pull requests"},
		[]string{"repository_full_name"})

	registry.MustRegister(eventPullRequest, workflowRun, pullRequest)
	return &executor{registry, eventPullRequest, workflowRun, pullRequest}
}

func (e *executor) Handler() gin.HandlerFunc {
	return gin.WrapH(promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{}))
}

func (e *executor) EventPullRequest(event github.PullRequestEvent) {
	// todo: delete old pull requests
	pr := core.ConvertPREToGitPullRequest(event)
	e.eventPullRequest.With(map[string]string{
		"action":               pr.Action,
		"pull_request_title":   pr.PullRequest.Title,
		"pull_request_state":   pr.PullRequest.State,
		"repository_full_name": pr.Repository.FullName,
	}).Inc()
	return
}

// EventPullRequestReview currently ignored
func (e *executor) EventPullRequestReview(github.PullRequestReviewEvent) {
	// ignored
	return
}

func (e *executor) EventWorkflowRun(run github.WorkflowRunEvent) {
	flow := core.ConvertToWorkflowRun(run).WorkflowRun
	e.workflowRun.DeletePartialMatch(map[string]string{
		"repository_full_name": flow.Repository.FullName,
		"workflow_run_id":      strconv.FormatInt(flow.WorkflowRunId, 10),
		"run_number":           strconv.Itoa(flow.RunNumber),
	})

	if flow.Conclusion != "success" {
		e.workflowRun.With(map[string]string{
			"repository_full_name": flow.Repository.FullName,
			"workflow_name":        flow.Name,
			"workflow_run_id":      strconv.FormatInt(flow.WorkflowRunId, 10),
			"run_number":           strconv.Itoa(flow.RunNumber),
			"state":                flow.Status,
			"conclusion":           flow.Conclusion,
		}).Set(1)
	}
}

func (e *executor) LastWorkflows(repository core.Repository, workflows []*github.WorkflowRun) {
	e.workflowRun.DeletePartialMatch(map[string]string{"repository_full_name": repository.Owner + "/" + repository.Name})
	for _, run := range workflows {
		flow := core.ConvertToWorkflow(*run)
		if flow.Conclusion != "success" {
			e.workflowRun.With(map[string]string{
				"repository_full_name": flow.Repository.FullName,
				"workflow_name":        flow.Name,
				"workflow_run_id":      strconv.FormatInt(flow.WorkflowRunId, 10),
				"run_number":           strconv.Itoa(flow.RunNumber),
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
