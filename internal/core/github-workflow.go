package core

import (
	"github.com/google/go-github/v61/github"
)

type WorkflowRun struct {
	WorkflowId    int64      `json:"workflow_id"`
	RunNumber     int        `json:"run_number"`
	Name          string     `json:"name"`
	HeadBranch    string     `json:"head_branch"`
	CommitMessage string     `json:"commit_message"`
	Event         string     `json:"event"`
	DisplayTitle  string     `json:"display_title"`
	Status        string     `json:"status"`
	Conclusion    string     `json:"conclusion"`
	HtmlUrl       string     `json:"html_url"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
	User          User       `json:"user"`
	Repository    Repository `json:"repository"`
}

func ConvertToWorkflow(w github.WorkflowRun) WorkflowRun {
	return WorkflowRun{
		WorkflowId:   w.GetWorkflowID(),
		RunNumber:    w.GetRunNumber(),
		Name:         w.GetName(),
		HeadBranch:   w.GetHeadBranch(),
		Event:        w.GetEvent(),
		DisplayTitle: w.GetDisplayTitle(),
		Status:       w.GetStatus(),
		Conclusion:   w.GetConclusion(),
		HtmlUrl:      w.GetHTMLURL(),
		CreatedAt:    w.GetCreatedAt().String(),
		UpdatedAt:    w.GetUpdatedAt().String(),
		User: User{
			Login: w.GetActor().GetLogin(),
		},
		CommitMessage: w.GetHeadCommit().GetMessage(),
		Repository: Repository{
			FullName: w.GetRepository().GetFullName(),
			HtmlUrl:  w.GetRepository().GetHTMLURL(),
		},
	}
}
