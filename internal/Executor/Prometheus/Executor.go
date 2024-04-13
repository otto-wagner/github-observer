package Prometheus

import (
	"github-listener/internal/Executor"
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus"
)

type executor struct {
	lastRequestReceivedTimeAction *prometheus.GaugeVec
	countAction                   *prometheus.CounterVec

	lastRequestReceivedTimePR *prometheus.GaugeVec
	countPR                   *prometheus.CounterVec

	lastRequestReceivedTimePRR *prometheus.GaugeVec
	countPRR                   *prometheus.CounterVec
}

func NewExecutor() Executor.IExecutor {
	return &executor{
		registerLastRequestReceivedTime("run_event"),
		registerCount("run_event"),
		registerLastRequestReceivedTime("pull_request_event"),
		registerCount("pull_request_event"),
		registerLastRequestReceivedTime("pull_request_review_event"),
		registerCount("pull_request_review_event"),
	}
}

func (e executor) Name() string {
	return "Prometheus"
}

func (e executor) CheckRunEvent(event github.CheckRunEvent) (err error) {
	e.lastRequestReceivedTimeAction.With(prometheus.Labels{
		"repo":       event.GetRepo().GetName(),
		"action":     event.GetAction(),
		"status":     event.GetCheckRun().GetStatus(),
		"conclusion": event.GetCheckRun().GetConclusion(),
		"state":      "",
	}).SetToCurrentTime()
	e.countAction.With(prometheus.Labels{
		"repo":       event.GetRepo().GetName(),
		"action":     event.GetAction(),
		"status":     event.GetCheckRun().GetStatus(),
		"conclusion": event.GetCheckRun().GetConclusion(),
	}).Inc()
	return
}

func (e executor) PullRequestEvent(event github.PullRequestEvent) (err error) {
	e.lastRequestReceivedTimeAction.With(prometheus.Labels{
		"repo":       event.GetRepo().GetName(),
		"action":     event.GetAction(),
		"status":     event.GetPullRequest().GetState(),
		"conclusion": "",
		"state":      "",
	}).SetToCurrentTime()
	e.countAction.With(prometheus.Labels{
		"repo":       event.GetRepo().GetName(),
		"action":     event.GetAction(),
		"status":     event.GetPullRequest().GetState(),
		"conclusion": "",
	}).Inc()
	return
}

func (e executor) PullRequestReviewEvent(event github.PullRequestReviewEvent) (err error) {
	e.lastRequestReceivedTimeAction.With(prometheus.Labels{
		"repo":       event.GetRepo().GetName(),
		"action":     event.GetAction(),
		"status":     event.GetPullRequest().GetState(),
		"conclusion": "",
		"state":      event.GetReview().GetState(),
	}).SetToCurrentTime()
	e.countAction.With(prometheus.Labels{
		"repo":       event.GetRepo().GetName(),
		"action":     event.GetAction(),
		"status":     event.GetPullRequest().GetState(),
		"conclusion": "",
	}).Inc()
	return
}

func registerLastRequestReceivedTime(name string) (gauge *prometheus.GaugeVec) {
	gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "last_request_received_time_" + name,
		Help: "Time when the last request was processed",
	}, []string{"repo", "action", "status", "conclusion", "state"})
	prometheus.MustRegister(gauge)
	return
}

func registerCount(name string) (counter *prometheus.CounterVec) {
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "count_request_processed_" + name,
			Help: "Number of requests processed"},
		[]string{"repo", "action", "status", "conclusion"})
	prometheus.MustRegister(counter)
	return
}
