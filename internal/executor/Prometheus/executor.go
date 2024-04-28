package Prometheus

import (
	"github-observer/internal/core"
	e "github-observer/internal/executor"
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus"
)

type executor struct {
	countAction   *prometheus.CounterVec
	countPR       *prometheus.CounterVec
	countPRR      *prometheus.CounterVec
	gaugeWorkflow *prometheus.GaugeVec
	registerPRSum *prometheus.GaugeVec
}

func NewExecutor() e.IExecutor {
	return &executor{
		registerActionCount(),
		registerPrCount(),
		registerPrReviewCount(),
		registerWorkflowRun(),
		registerPRSum(),
	}
}

func (e *executor) Name() string {
	return "Prometheus"
}

func (e *executor) RunEvent(event github.CheckRunEvent) {
	action := core.ConvertToGitAction(event)
	e.countAction.With(action.ToMap()).Inc()
	return
}

func (e *executor) PullRequestEvent(event github.PullRequestEvent) {
	pr := core.ConvertPREToGitPullRequest(event)
	e.countPR.With(pr.ToMap()).Inc()
	return
}

func (e *executor) PullRequestReviewEvent(event github.PullRequestReviewEvent) {
	prr := core.ConvertToGitPullRequestReview(event)
	e.countPRR.With(prr.ToMap()).Inc()
	return
}

func (e *executor) WorkflowRuns(allRuns []*github.WorkflowRun) {
	repoBranchRuns := make(map[core.Repository]map[string][]core.WorkflowRun)
	for _, run := range allRuns {
		workflow := core.ConvertToWorkflow(*run)
		if _, ok := repoBranchRuns[workflow.Repository]; !ok {
			repoBranchRuns[workflow.Repository] = make(map[string][]core.WorkflowRun)
		}
		repoBranchRuns[workflow.Repository][workflow.HeadBranch] = append(repoBranchRuns[workflow.Repository][workflow.HeadBranch], workflow)
	}

	for _, branches := range repoBranchRuns {
		for _, workflows := range branches {
			var lastWorkflow core.WorkflowRun
			for _, workflow := range workflows {
				if lastWorkflow.RunNumber < workflow.RunNumber {
					lastWorkflow = workflow
					if workflow.Conclusion != "success" {
						e.gaugeWorkflow.With(workflow.ToMap()).Set(1)
					} else {
						e.gaugeWorkflow.With(workflow.ToMap()).Set(0)
					}
				}
			}
		}
	}

}

func (e *executor) PullRequests(repository core.Repository, pullRequests []*github.PullRequest) {
	var count int
	for _, pullRequest := range pullRequests {
		if pullRequest.GetState() != "closed" && pullRequest.GetState() != "merged" {
			count++
		}
	}
	e.registerPRSum.With(map[string]string{"repository_name": repository.Name, "repository_url": repository.HtmlUrl}).Set(float64(count))

	return
}

func registerActionCount() (counter *prometheus.CounterVec) {
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "count_request_processed_action_event",
			Help: "Number of action requests processed"},
		[]string{"action", "check_run_name", "check_run_url", "check_run_status", "check_run_conclusion",
			"check_run_started_at", "check_run_completed", "repository_name", "repository_url"})
	prometheus.MustRegister(counter)
	return
}

func registerPrCount() (counter *prometheus.CounterVec) {
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "count_request_processed_pull_request_event",
			Help: "Number of pull requests processed"},
		[]string{"action", "pull_request_number", "pull_request_title", "pull_request_body", "pull_request_user", "pull_request_url", "pull_request_state",
			"pull_request_created_at", "pull_request_updated_at", "pull_request_closed_at", "pull_request_merged_at", "repository_name", "repository_url"})
	prometheus.MustRegister(counter)
	return
}

func registerPrReviewCount() (counter *prometheus.CounterVec) {
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "count_request_processed_pull_request_review_event",
			Help: "Number of pull request reviews processed"},
		[]string{"action", "pull_request_title", "pull_request_user", "pull_request_url", "pull_request_state",
			"review_body", "review_state", "review_user", "review_url", "repository_name", "repository_url"})
	prometheus.MustRegister(counter)
	return
}

func registerWorkflowRun() (gauge *prometheus.GaugeVec) {
	gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "workflow_run_failed",
			Help: "Failed workflow runs"},
		[]string{"workflow_id", "run_number", "name", "head_branch", "event", "display_title", "conclusion", "html_url",
			"created_at", "updated_at", "user", "commit_message", "repository_name", "repository_html_url"})
	prometheus.MustRegister(gauge)
	return
}

func registerPRSum() (gauge *prometheus.GaugeVec) {
	gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pull_request_sum",
			Help: "Number of pull requests"},
		[]string{"repository_name", "repository_url"})
	prometheus.MustRegister(gauge)
	return
}
