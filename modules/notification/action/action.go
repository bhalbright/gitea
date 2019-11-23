// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package action

import (
	"fmt"
	"path"
	"strings"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/notification/base"
)

type actionNotifier struct {
	base.NullNotifier
}

var (
	_ base.Notifier = &actionNotifier{}
)

// NewNotifier create a new actionNotifier notifier
func NewNotifier() base.Notifier {
	return &actionNotifier{}
}

func (a *actionNotifier) NotifyNewIssue(issue *models.Issue) {
	if err := issue.LoadPoster(); err != nil {
		log.Error("issue.LoadPoster: %v", err)
		return
	}
	if err := issue.LoadRepo(); err != nil {
		log.Error("issue.LoadRepo: %v", err)
		return
	}
	repo := issue.Repo

	if err := models.NotifyWatchers(&models.Action{
		ActUserID: issue.Poster.ID,
		ActUser:   issue.Poster,
		OpType:    models.ActionCreateIssue,
		Content:   fmt.Sprintf("%d|%s", issue.Index, issue.Title),
		RepoID:    repo.ID,
		Repo:      repo,
		IsPrivate: repo.IsPrivate,
	}); err != nil {
		log.Error("NotifyWatchers: %v", err)
	}
}

func (a *actionNotifier) NotifyNewPullRequest(pull *models.PullRequest) {
	if err := pull.LoadIssue(); err != nil {
		log.Error("pull.LoadIssue: %v", err)
		return
	}
	if err := pull.Issue.LoadRepo(); err != nil {
		log.Error("pull.Issue.LoadRepo: %v", err)
		return
	}
	if err := pull.Issue.LoadPoster(); err != nil {
		log.Error("pull.Issue.LoadPoster: %v", err)
		return
	}

	if err := models.NotifyWatchers(&models.Action{
		ActUserID: pull.Issue.Poster.ID,
		ActUser:   pull.Issue.Poster,
		OpType:    models.ActionCreatePullRequest,
		Content:   fmt.Sprintf("%d|%s", pull.Issue.Index, pull.Issue.Title),
		RepoID:    pull.Issue.Repo.ID,
		Repo:      pull.Issue.Repo,
		IsPrivate: pull.Issue.Repo.IsPrivate,
	}); err != nil {
		log.Error("NotifyWatchers: %v", err)
	}
}

func (a *actionNotifier) NotifyRenameRepository(doer *models.User, repo *models.Repository, oldRepoName string) {
	log.Trace("action.ChangeRepositoryName: %s/%s", doer.Name, repo.Name)

	if err := models.NotifyWatchers(&models.Action{
		ActUserID: doer.ID,
		ActUser:   doer,
		OpType:    models.ActionRenameRepo,
		RepoID:    repo.ID,
		Repo:      repo,
		IsPrivate: repo.IsPrivate,
		Content:   oldRepoName,
	}); err != nil {
		log.Error("NotifyWatchers: %v", err)
	}
}

func (a *actionNotifier) NotifyTransferRepository(doer *models.User, repo *models.Repository, oldOwnerName string) {
	if err := models.NotifyWatchers(&models.Action{
		ActUserID: doer.ID,
		ActUser:   doer,
		OpType:    models.ActionTransferRepo,
		RepoID:    repo.ID,
		Repo:      repo,
		IsPrivate: repo.IsPrivate,
		Content:   path.Join(oldOwnerName, repo.Name),
	}); err != nil {
		log.Error("NotifyWatchers: %v", err)
	}
}

func (a *actionNotifier) NotifyCreateRepository(doer *models.User, u *models.User, repo *models.Repository) {
	if err := models.NotifyWatchers(&models.Action{
		ActUserID: doer.ID,
		ActUser:   doer,
		OpType:    models.ActionCreateRepo,
		RepoID:    repo.ID,
		Repo:      repo,
		IsPrivate: repo.IsPrivate,
	}); err != nil {
		log.Error("notify watchers '%d/%d': %v", doer.ID, repo.ID, err)
	}
}

func (a *actionNotifier) NotifyForkRepository(doer *models.User, oldRepo, repo *models.Repository) {
	if err := models.NotifyWatchers(&models.Action{
		ActUserID: doer.ID,
		ActUser:   doer,
		OpType:    models.ActionCreateRepo,
		RepoID:    repo.ID,
		Repo:      repo,
		IsPrivate: repo.IsPrivate,
	}); err != nil {
		log.Error("notify watchers '%d/%d': %v", doer.ID, repo.ID, err)
	}
}

func (a *actionNotifier) NotifyPullRequestReview(pr *models.PullRequest, review *models.Review, comment *models.Comment) {
	if err := review.LoadReviewer(); err != nil {
		log.Error("LoadReviewer '%d/%d': %v", review.ID, review.ReviewerID, err)
		return
	}
	if err := review.LoadCodeComments(); err != nil {
		log.Error("LoadCodeComments '%d/%d': %v", review.Reviewer.ID, review.ID, err)
		return
	}

	var actions = make([]*models.Action, 0, 10)
	for _, lines := range review.CodeComments {
		for _, comments := range lines {
			for _, comm := range comments {
				actions = append(actions, &models.Action{
					ActUserID: review.Reviewer.ID,
					ActUser:   review.Reviewer,
					Content:   fmt.Sprintf("%d|%s", review.Issue.Index, strings.Split(comm.Content, "\n")[0]),
					OpType:    models.ActionCommentIssue,
					RepoID:    review.Issue.RepoID,
					Repo:      review.Issue.Repo,
					IsPrivate: review.Issue.Repo.IsPrivate,
					Comment:   comm,
					CommentID: comm.ID,
				})
			}
		}
	}

	if review.Type != models.ReviewTypeComment || strings.TrimSpace(comment.Content) != "" {
		action := &models.Action{
			ActUserID: review.Reviewer.ID,
			ActUser:   review.Reviewer,
			Content:   fmt.Sprintf("%d|%s", review.Issue.Index, strings.Split(comment.Content, "\n")[0]),
			RepoID:    review.Issue.RepoID,
			Repo:      review.Issue.Repo,
			IsPrivate: review.Issue.Repo.IsPrivate,
			Comment:   comment,
			CommentID: comment.ID,
		}

		switch review.Type {
		case models.ReviewTypeApprove:
			action.OpType = models.ActionApprovePullRequest
		case models.ReviewTypeReject:
			action.OpType = models.ActionRejectPullRequest
		default:
			action.OpType = models.ActionCommentIssue
		}

		actions = append(actions, action)
	}

	if err := models.NotifyWatchersActions(actions); err != nil {
		log.Error("notify watchers '%d/%d': %v", review.Reviewer.ID, review.Issue.RepoID, err)
	}
}

func (*actionNotifier) NotifyMergePullRequest(pr *models.PullRequest, doer *models.User, baseRepo *git.Repository) {
	if err := models.NotifyWatchers(&models.Action{
		ActUserID: doer.ID,
		ActUser:   doer,
		OpType:    models.ActionMergePullRequest,
		Content:   fmt.Sprintf("%d|%s", pr.Issue.Index, pr.Issue.Title),
		RepoID:    pr.Issue.Repo.ID,
		Repo:      pr.Issue.Repo,
		IsPrivate: pr.Issue.Repo.IsPrivate,
	}); err != nil {
		log.Error("NotifyWatchers [%d]: %v", pr.ID, err)
	}
}
