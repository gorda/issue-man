package comm

import (
	"github.com/google/uuid"
	"gopkg.in/go-playground/webhooks.v5/github"
	"issue-man/global"
)

// 存储 IssueCommentPayload 里的一些信息
// 基本是目前进行各种操作需要用到的信息
type Info struct {
	// 仓库信息
	Owner      string
	Repository string

	// 评论人信息
	Login string
	// 评论提及到的人
	Mention []string

	// Issue 目前的信息
	IssueURL    string
	IssueNumber int
	Title       string
	Body        string
	Milestone   int
	State       string
	Assignees   []string
	Labels      []string

	// 一个指令的 UUID
	ReqID string
}

// Info
// 从 IssueCommentPayload 里的一些信息
// 避免多次书写出现错误
func (p Info) Parse(payload github.IssueCommentPayload) (info Info) {
	defer func() {
		if p := recover(); p != nil {
			global.Sugar.Errorw("Info panic",
				"req_id", info.ReqID,
				"panic", p)
		}
	}()
	info.ReqID = uuid.New().String()
	info.Owner = payload.Repository.Owner.Login
	info.Repository = payload.Repository.Name

	info.Login = payload.Sender.Login

	info.IssueURL = payload.Issue.URL
	info.IssueNumber = int(payload.Issue.Number)
	info.Title = payload.Issue.Title
	info.Body = payload.Issue.Body
	info.Milestone = int(payload.Issue.Milestone.Number)
	info.State = payload.Issue.State

	info.Assignees = make([]string, len(payload.Issue.Assignees))
	info.Labels = make([]string, len(payload.Issue.Labels))
	for i := 0; i < len(payload.Issue.Assignees) || i < len(payload.Issue.Labels); i++ {
		if i < len(info.Assignees) {
			info.Assignees[i] = payload.Issue.Assignees[i].Login
		}
		if i < len(info.Labels) {
			info.Labels[i] = payload.Issue.Labels[i].Name
		}
	}
	return
}