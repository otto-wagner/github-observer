package core

import (
	"github.com/google/go-github/v61/github"
	"strconv"
)

type WorkflowRun struct {
	WorkflowId    int64      `json:"workflow_id"`
	RunNumber     int        `json:"run_number"`
	Name          string     `json:"name"`
	HeadBranch    string     `json:"head_branch"`
	CommitMessage string     `json:"commit_message"`
	Event         string     `json:"event"`
	DisplayTitle  string     `json:"display_title"`
	Conclusion    string     `json:"conclusion"`
	HtmlUrl       string     `json:"html_url"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
	User          User       `json:"user"`
	Repository    Repository `json:"repository"`
}

func ConvertToWorkflow(flow github.WorkflowRun) WorkflowRun {
	return WorkflowRun{
		WorkflowId:   flow.GetWorkflowID(),
		RunNumber:    flow.GetRunNumber(),
		Name:         flow.GetName(),
		HeadBranch:   flow.GetHeadBranch(),
		Event:        flow.GetEvent(),
		DisplayTitle: flow.GetDisplayTitle(),
		Conclusion:   flow.GetConclusion(),
		HtmlUrl:      flow.GetHTMLURL(),
		CreatedAt:    flow.GetCreatedAt().String(),
		UpdatedAt:    flow.GetUpdatedAt().String(),
		User: User{
			Login: flow.GetActor().GetLogin(),
		},
		CommitMessage: flow.GetHeadCommit().GetMessage(),
		Repository: Repository{
			Name:    flow.GetRepository().GetName(),
			HtmlUrl: flow.GetRepository().GetHTMLURL(),
		},
	}
}

func (w *WorkflowRun) ToMap() map[string]string {
	return map[string]string{
		"workflow_id":         strconv.FormatInt(w.WorkflowId, 10),
		"run_number":          strconv.Itoa(w.RunNumber),
		"name":                w.Name,
		"head_branch":         w.HeadBranch,
		"event":               w.Event,
		"display_title":       w.DisplayTitle,
		"conclusion":          w.Conclusion,
		"html_url":            w.HtmlUrl,
		"created_at":          w.CreatedAt,
		"updated_at":          w.UpdatedAt,
		"user":                w.User.Login,
		"commit_message":      w.CommitMessage,
		"repository_name":     w.Repository.Name,
		"repository_html_url": w.Repository.HtmlUrl,
	}

}
