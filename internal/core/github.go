package core

import (
	"strings"

	"github.com/google/go-github/v73/github"
)

type User struct {
	Login string `json:"login"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Branch   string `json:"branch"`
	Owner    string `json:"owner"`
	HtmlUrl  string `json:"html_url"`
}

func ToRepository(repo string) Repository {
	first := strings.Split(repo, "/")
	if len(first) != 2 {
		return Repository{}
	}
	second := strings.Split(first[1], "@")
	if len(second) != 2 {
		return Repository{}
	}
	return Repository{
		Owner:  first[0],
		Name:   second[0],
		Branch: second[1],
	}
}

type GitPullRequest struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

type PullRequest struct {
	Number    int    `json:"number"`
	HtmlUrl   string `json:"html_url"`
	IssueUrl  string `json:"issue_url"`
	State     string `json:"state"`
	Title     string `json:"title"`
	User      User   `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  string `json:"closed_at"`
	MergedAt  string `json:"merged_at"`
}

func ConvertPREToGitPullRequest(event github.PullRequestEvent) GitPullRequest {
	return GitPullRequest{
		Action: event.GetAction(),
		PullRequest: PullRequest{
			Number:    event.GetPullRequest().GetNumber(),
			HtmlUrl:   event.GetPullRequest().GetHTMLURL(),
			State:     event.GetPullRequest().GetState(),
			Title:     event.GetPullRequest().GetTitle(),
			User:      User{Login: event.GetPullRequest().GetUser().GetLogin()},
			CreatedAt: event.GetPullRequest().GetCreatedAt().String(),
			UpdatedAt: event.GetPullRequest().GetUpdatedAt().String(),
			ClosedAt:  event.GetPullRequest().GetClosedAt().String(),
			MergedAt:  event.GetPullRequest().GetMergedAt().String(),
		},
		Repository: Repository{
			FullName: event.GetRepo().GetFullName(),
			HtmlUrl:  event.GetRepo().GetHTMLURL(),
		},
	}
}

func ConvertPRToGitPullRequest(p github.PullRequest) GitPullRequest {
	return GitPullRequest{
		PullRequest: PullRequest{
			Number:    p.GetNumber(),
			HtmlUrl:   p.GetHTMLURL(),
			State:     p.GetState(),
			Title:     p.GetTitle(),
			User:      User{Login: p.GetUser().GetLogin()},
			CreatedAt: p.GetCreatedAt().String(),
			UpdatedAt: p.GetUpdatedAt().String(),
			ClosedAt:  p.GetClosedAt().String(),
			MergedAt:  p.GetMergedAt().String(),
		},
		Repository: Repository{
			FullName: p.GetBase().GetRepo().GetFullName(),
			HtmlUrl:  p.GetBase().GetRepo().GetHTMLURL(),
		},
	}
}

type WorkflowRun struct {
	WorkflowId    int64      `json:"workflow_id"`
	WorkflowRunId int64      `json:"workflow_run_id"`
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

type WorkflowRunEvent struct {
	Action      string      `json:"action"`
	WorkflowRun WorkflowRun `json:"workflow_run"`
}

func ConvertToWorkflow(w github.WorkflowRun) WorkflowRun {
	return WorkflowRun{
		WorkflowId:    w.GetWorkflowID(),
		WorkflowRunId: w.GetID(),
		RunNumber:     w.GetRunNumber(),
		Name:          w.GetName(),
		HeadBranch:    w.GetHeadBranch(),
		Event:         w.GetEvent(),
		DisplayTitle:  w.GetDisplayTitle(),
		Status:        w.GetStatus(),
		Conclusion:    w.GetConclusion(),
		HtmlUrl:       w.GetHTMLURL(),
		CreatedAt:     w.GetCreatedAt().String(),
		UpdatedAt:     w.GetUpdatedAt().String(),
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

func ConvertToWorkflowRun(event github.WorkflowRunEvent) WorkflowRunEvent {
	return WorkflowRunEvent{
		Action:      event.GetAction(),
		WorkflowRun: ConvertToWorkflow(*event.GetWorkflowRun()),
	}
}
