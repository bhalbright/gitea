// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package convert

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"code.gitea.io/gitea/graph/model"
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/markup"
	"code.gitea.io/gitea/modules/structs"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/util"
	"code.gitea.io/gitea/modules/webhook"

	"github.com/unknwon/com"
)

// ToEmail convert models.EmailAddress to api.Email
func ToEmail(email *models.EmailAddress) *api.Email {
	return &api.Email{
		Email:    email.Email,
		Verified: email.IsActivated,
		Primary:  email.IsPrimary,
	}
}

// ToBranch convert a git.Commit and git.Branch to an api.Branch
func ToBranch(repo *models.Repository, b *git.Branch, c *git.Commit, bp *models.ProtectedBranch, user *models.User, isRepoAdmin bool) (*api.Branch, error) {
	if bp == nil {
		var hasPerm bool
		var err error
		if user != nil {
			hasPerm, err = models.HasAccessUnit(user, repo, models.UnitTypeCode, models.AccessModeWrite)
			if err != nil {
				return nil, err
			}
		}

		return &api.Branch{
			Name:                b.Name,
			Commit:              ToCommit(repo, c),
			Protected:           false,
			RequiredApprovals:   0,
			EnableStatusCheck:   false,
			StatusCheckContexts: []string{},
			UserCanPush:         hasPerm,
			UserCanMerge:        hasPerm,
		}, nil
	}

	branch := &api.Branch{
		Name:                b.Name,
		Commit:              ToCommit(repo, c),
		Protected:           true,
		RequiredApprovals:   bp.RequiredApprovals,
		EnableStatusCheck:   bp.EnableStatusCheck,
		StatusCheckContexts: bp.StatusCheckContexts,
	}

	if isRepoAdmin {
		branch.EffectiveBranchProtectionName = bp.BranchName
	}

	if user != nil {
		branch.UserCanPush = bp.CanUserPush(user.ID)
		branch.UserCanMerge = bp.IsUserMergeWhitelisted(user.ID)
	}

	return branch, nil
}

// ToBranchProtection convert a ProtectedBranch to api.BranchProtection
func ToBranchProtection(bp *models.ProtectedBranch) *api.BranchProtection {
	pushWhitelistUsernames, err := models.GetUserNamesByIDs(bp.WhitelistUserIDs)
	if err != nil {
		log.Error("GetUserNamesByIDs (WhitelistUserIDs): %v", err)
	}
	mergeWhitelistUsernames, err := models.GetUserNamesByIDs(bp.MergeWhitelistUserIDs)
	if err != nil {
		log.Error("GetUserNamesByIDs (MergeWhitelistUserIDs): %v", err)
	}
	approvalsWhitelistUsernames, err := models.GetUserNamesByIDs(bp.ApprovalsWhitelistUserIDs)
	if err != nil {
		log.Error("GetUserNamesByIDs (ApprovalsWhitelistUserIDs): %v", err)
	}
	pushWhitelistTeams, err := models.GetTeamNamesByID(bp.WhitelistTeamIDs)
	if err != nil {
		log.Error("GetTeamNamesByID (WhitelistTeamIDs): %v", err)
	}
	mergeWhitelistTeams, err := models.GetTeamNamesByID(bp.MergeWhitelistTeamIDs)
	if err != nil {
		log.Error("GetTeamNamesByID (MergeWhitelistTeamIDs): %v", err)
	}
	approvalsWhitelistTeams, err := models.GetTeamNamesByID(bp.ApprovalsWhitelistTeamIDs)
	if err != nil {
		log.Error("GetTeamNamesByID (ApprovalsWhitelistTeamIDs): %v", err)
	}

	return &api.BranchProtection{
		BranchName:                  bp.BranchName,
		EnablePush:                  bp.CanPush,
		EnablePushWhitelist:         bp.EnableWhitelist,
		PushWhitelistUsernames:      pushWhitelistUsernames,
		PushWhitelistTeams:          pushWhitelistTeams,
		PushWhitelistDeployKeys:     bp.WhitelistDeployKeys,
		EnableMergeWhitelist:        bp.EnableMergeWhitelist,
		MergeWhitelistUsernames:     mergeWhitelistUsernames,
		MergeWhitelistTeams:         mergeWhitelistTeams,
		EnableStatusCheck:           bp.EnableStatusCheck,
		StatusCheckContexts:         bp.StatusCheckContexts,
		RequiredApprovals:           bp.RequiredApprovals,
		EnableApprovalsWhitelist:    bp.EnableApprovalsWhitelist,
		ApprovalsWhitelistUsernames: approvalsWhitelistUsernames,
		ApprovalsWhitelistTeams:     approvalsWhitelistTeams,
		BlockOnRejectedReviews:      bp.BlockOnRejectedReviews,
		BlockOnOutdatedBranch:       bp.BlockOnOutdatedBranch,
		DismissStaleApprovals:       bp.DismissStaleApprovals,
		RequireSignedCommits:        bp.RequireSignedCommits,
		ProtectedFilePatterns:       bp.ProtectedFilePatterns,
		Created:                     bp.CreatedUnix.AsTime(),
		Updated:                     bp.UpdatedUnix.AsTime(),
	}
}

// ToTag convert a git.Tag to an api.Tag
func ToTag(repo *models.Repository, t *git.Tag) *api.Tag {
	return &api.Tag{
		Name:       t.Name,
		ID:         t.ID.String(),
		Commit:     ToCommitMeta(repo, t),
		ZipballURL: util.URLJoin(repo.HTMLURL(), "archive", t.Name+".zip"),
		TarballURL: util.URLJoin(repo.HTMLURL(), "archive", t.Name+".tar.gz"),
	}
}

// ToCommit convert a git.Commit to api.PayloadCommit
func ToCommit(repo *models.Repository, c *git.Commit) *api.PayloadCommit {
	authorUsername := ""
	if author, err := models.GetUserByEmail(c.Author.Email); err == nil {
		authorUsername = author.Name
	} else if !models.IsErrUserNotExist(err) {
		log.Error("GetUserByEmail: %v", err)
	}

	committerUsername := ""
	if committer, err := models.GetUserByEmail(c.Committer.Email); err == nil {
		committerUsername = committer.Name
	} else if !models.IsErrUserNotExist(err) {
		log.Error("GetUserByEmail: %v", err)
	}

	return &api.PayloadCommit{
		ID:      c.ID.String(),
		Message: c.Message(),
		URL:     util.URLJoin(repo.HTMLURL(), "commit", c.ID.String()),
		Author: &api.PayloadUser{
			Name:     c.Author.Name,
			Email:    c.Author.Email,
			UserName: authorUsername,
		},
		Committer: &api.PayloadUser{
			Name:     c.Committer.Name,
			Email:    c.Committer.Email,
			UserName: committerUsername,
		},
		Timestamp:    c.Author.When,
		Verification: ToVerification(c),
	}
}

// ToVerification convert a git.Commit.Signature to an api.PayloadCommitVerification
func ToVerification(c *git.Commit) *api.PayloadCommitVerification {
	verif := models.ParseCommitWithSignature(c)
	commitVerification := &api.PayloadCommitVerification{
		Verified: verif.Verified,
		Reason:   verif.Reason,
	}
	if c.Signature != nil {
		commitVerification.Signature = c.Signature.Signature
		commitVerification.Payload = c.Signature.Payload
	}
	if verif.SigningUser != nil {
		commitVerification.Signer = &structs.PayloadUser{
			Name:  verif.SigningUser.Name,
			Email: verif.SigningUser.Email,
		}
	}
	return commitVerification
}

// ToPublicKey convert models.PublicKey to api.PublicKey
func ToPublicKey(apiLink string, key *models.PublicKey) *api.PublicKey {
	return &api.PublicKey{
		ID:          key.ID,
		Key:         key.Content,
		URL:         apiLink + com.ToStr(key.ID),
		Title:       key.Name,
		Fingerprint: key.Fingerprint,
		Created:     key.CreatedUnix.AsTime(),
	}
}

// ToGPGKey converts models.GPGKey to api.GPGKey
func ToGPGKey(key *models.GPGKey) *api.GPGKey {
	subkeys := make([]*api.GPGKey, len(key.SubsKey))
	for id, k := range key.SubsKey {
		subkeys[id] = &api.GPGKey{
			ID:                k.ID,
			PrimaryKeyID:      k.PrimaryKeyID,
			KeyID:             k.KeyID,
			PublicKey:         k.Content,
			Created:           k.CreatedUnix.AsTime(),
			Expires:           k.ExpiredUnix.AsTime(),
			CanSign:           k.CanSign,
			CanEncryptComms:   k.CanEncryptComms,
			CanEncryptStorage: k.CanEncryptStorage,
			CanCertify:        k.CanSign,
		}
	}
	emails := make([]*api.GPGKeyEmail, len(key.Emails))
	for i, e := range key.Emails {
		emails[i] = ToGPGKeyEmail(e)
	}
	return &api.GPGKey{
		ID:                key.ID,
		PrimaryKeyID:      key.PrimaryKeyID,
		KeyID:             key.KeyID,
		PublicKey:         key.Content,
		Created:           key.CreatedUnix.AsTime(),
		Expires:           key.ExpiredUnix.AsTime(),
		Emails:            emails,
		SubsKey:           subkeys,
		CanSign:           key.CanSign,
		CanEncryptComms:   key.CanEncryptComms,
		CanEncryptStorage: key.CanEncryptStorage,
		CanCertify:        key.CanSign,
	}
}

// ToGPGKeyEmail convert models.EmailAddress to api.GPGKeyEmail
func ToGPGKeyEmail(email *models.EmailAddress) *api.GPGKeyEmail {
	return &api.GPGKeyEmail{
		Email:    email.Email,
		Verified: email.IsActivated,
	}
}

// ToHook convert models.Webhook to api.Hook
func ToHook(repoLink string, w *models.Webhook) *api.Hook {
	config := map[string]string{
		"url":          w.URL,
		"content_type": w.ContentType.Name(),
	}
	if w.HookTaskType == models.SLACK {
		s := webhook.GetSlackHook(w)
		config["channel"] = s.Channel
		config["username"] = s.Username
		config["icon_url"] = s.IconURL
		config["color"] = s.Color
	}

	return &api.Hook{
		ID:      w.ID,
		Type:    w.HookTaskType.Name(),
		URL:     fmt.Sprintf("%s/settings/hooks/%d", repoLink, w.ID),
		Active:  w.IsActive,
		Config:  config,
		Events:  w.EventsArray(),
		Updated: w.UpdatedUnix.AsTime(),
		Created: w.CreatedUnix.AsTime(),
	}
}

// ToGitHook convert git.Hook to api.GitHook
func ToGitHook(h *git.Hook) *api.GitHook {
	return &api.GitHook{
		Name:     h.Name(),
		IsActive: h.IsActive,
		Content:  h.Content,
	}
}

// ToDeployKey convert models.DeployKey to api.DeployKey
func ToDeployKey(apiLink string, key *models.DeployKey) *api.DeployKey {
	return &api.DeployKey{
		ID:          key.ID,
		KeyID:       key.KeyID,
		Key:         key.Content,
		Fingerprint: key.Fingerprint,
		URL:         apiLink + com.ToStr(key.ID),
		Title:       key.Name,
		Created:     key.CreatedUnix.AsTime(),
		ReadOnly:    key.Mode == models.AccessModeRead, // All deploy keys are read-only.
	}
}

// ToOrganization convert models.User to api.Organization
func ToOrganization(org *models.User) *api.Organization {
	return &api.Organization{
		ID:                        org.ID,
		AvatarURL:                 org.AvatarLink(),
		UserName:                  org.Name,
		FullName:                  org.FullName,
		Description:               org.Description,
		Website:                   org.Website,
		Location:                  org.Location,
		Visibility:                org.Visibility.String(),
		RepoAdminChangeTeamAccess: org.RepoAdminChangeTeamAccess,
	}
}

// ToTeam convert models.Team to api.Team
func ToTeam(team *models.Team) *api.Team {
	return &api.Team{
		ID:                      team.ID,
		Name:                    team.Name,
		Description:             team.Description,
		IncludesAllRepositories: team.IncludesAllRepositories,
		CanCreateOrgRepo:        team.CanCreateOrgRepo,
		Permission:              team.Authorize.String(),
		Units:                   team.GetUnitNames(),
	}
}

// ToUser convert models.User to api.User
// signed shall only be set if requester is logged in. authed shall only be set if user is site admin or user himself
func ToUser(user *models.User, signed, authed bool) *api.User {
	result := &api.User{
		UserName:  user.Name,
		AvatarURL: user.AvatarLink(),
		FullName:  markup.Sanitize(user.FullName),
		Created:   user.CreatedUnix.AsTime(),
	}
	// hide primary email if API caller is anonymous or user keep email private
	if signed && (!user.KeepEmailPrivate || authed) {
		result.Email = user.Email
	}
	// only site admin will get these information and possibly user himself
	if authed {
		result.ID = user.ID
		result.IsAdmin = user.IsAdmin
		result.LastLogin = user.LastLoginUnix.AsTime()
		result.Language = user.Language
	}
	return result
}

// ToAnnotatedTag convert git.Tag to api.AnnotatedTag
func ToAnnotatedTag(repo *models.Repository, t *git.Tag, c *git.Commit) *api.AnnotatedTag {
	return &api.AnnotatedTag{
		Tag:          t.Name,
		SHA:          t.ID.String(),
		Object:       ToAnnotatedTagObject(repo, c),
		Message:      t.Message,
		URL:          util.URLJoin(repo.APIURL(), "git/tags", t.ID.String()),
		Tagger:       ToCommitUser(t.Tagger),
		Verification: ToVerification(c),
	}
}

// ToAnnotatedTagObject convert a git.Commit to an api.AnnotatedTagObject
func ToAnnotatedTagObject(repo *models.Repository, commit *git.Commit) *api.AnnotatedTagObject {
	return &api.AnnotatedTagObject{
		SHA:  commit.ID.String(),
		Type: string(git.ObjectCommit),
		URL:  util.URLJoin(repo.APIURL(), "git/commits", commit.ID.String()),
	}
}

// ToCommitUser convert a git.Signature to an api.CommitUser
func ToCommitUser(sig *git.Signature) *api.CommitUser {
	return &api.CommitUser{
		Identity: api.Identity{
			Name:  sig.Name,
			Email: sig.Email,
		},
		Date: sig.When.UTC().Format(time.RFC3339),
	}
}

// ToCommitMeta convert a git.Tag to an api.CommitMeta
func ToCommitMeta(repo *models.Repository, tag *git.Tag) *api.CommitMeta {
	return &api.CommitMeta{
		SHA: tag.Object.String(),
		URL: util.URLJoin(repo.APIURL(), "git/commits", tag.ID.String()),
	}
}

// ToTopicResponse convert from models.Topic to api.TopicResponse
func ToTopicResponse(topic *models.Topic) *api.TopicResponse {
	return &api.TopicResponse{
		ID:        topic.ID,
		Name:      topic.Name,
		RepoCount: topic.RepoCount,
		Created:   topic.CreatedUnix.AsTime(),
		Updated:   topic.UpdatedUnix.AsTime(),
	}
}

// ToOAuth2Application convert from models.OAuth2Application to api.OAuth2Application
func ToOAuth2Application(app *models.OAuth2Application) *api.OAuth2Application {
	return &api.OAuth2Application{
		ID:           app.ID,
		Name:         app.Name,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		RedirectURIs: app.RedirectURIs,
		Created:      app.CreatedUnix.AsTime(),
	}
}

// ToGraphRepository convert from models.Repository to graphmodel.Repository
func ToGraphRepository(repo *models.Repository, mode models.AccessMode, signed, authed bool) *model.Repository {
	apiRepo := repo.APIFormat(mode)
	graphRepo := innerToGraphRepository(apiRepo)
	if apiRepo.Parent != nil {
		graphRepo.Parent = innerToGraphRepository(apiRepo.Parent)
	}
	graphRepo.Owner = innerToGraphUser(repo.Owner.APIFormat())
	return graphRepo
}

func innerToGraphRepository(apiRepo *structs.Repository) *model.Repository {
	return &model.Repository{
		ID:                        ToGraphId("repository", apiRepo.ID),
		RestAPIID:                 &apiRepo.ID,
		Name:                      &apiRepo.Name,
		FullName:                  &apiRepo.FullName,
		Description:               &apiRepo.Description,
		Private:                   &apiRepo.Private,
		Template:                  &apiRepo.Template,
		Empty:                     &apiRepo.Empty,
		Archived:                  &apiRepo.Archived,
		Size:                      &apiRepo.Size,
		Fork:                      &apiRepo.Fork,
		Mirror:                    &apiRepo.Mirror,
		HTMLURL:                   &apiRepo.HTMLURL,
		SSHURL:                    &apiRepo.SSHURL,
		CloneURL:                  &apiRepo.CloneURL,
		Website:                   &apiRepo.Website,
		Stars:                     &apiRepo.Stars,
		Forks:                     &apiRepo.Forks,
		Watchers:                  &apiRepo.Watchers,
		OpenIssues:                &apiRepo.OpenIssues,
		OpenPulls:                 &apiRepo.OpenPulls,
		Releases:                  &apiRepo.Releases,
		DefaultBranch:             &apiRepo.DefaultBranch,
		Created:                   &apiRepo.Created,
		Updated:                   &apiRepo.Updated,
		Permissions:               ToGraphPermissions(apiRepo.Permissions),
		HasIssues:                 &apiRepo.HasIssues,
		ExternalTracker:           ToGraphExternalTracker(apiRepo.ExternalTracker),
		InternalTracker:           ToGraphInternalTracker(apiRepo.InternalTracker),
		HasWiki:                   &apiRepo.HasWiki,
		ExternalWiki:              ToGraphExternalWiki(apiRepo.ExternalWiki),
		HasPullRequests:           &apiRepo.HasPullRequests,
		IgnoreWhitespaceConflicts: &apiRepo.IgnoreWhitespaceConflicts,
		AllowMerge:                &apiRepo.AllowMerge,
		AllowRebase:               &apiRepo.AllowRebase,
		AllowRebaseMerge:          &apiRepo.AllowRebaseMerge,
		AllowSquash:               &apiRepo.AllowSquash,
		AvatarURL:                 &apiRepo.AvatarURL,
		Internal:                  &apiRepo.Internal,
	}
}

// ToGraphPermissions convert an api.Permission to a graphql Permission
func ToGraphPermissions(apiPermissions *structs.Permission) *model.Permission {
	return &model.Permission{
		Admin: &apiPermissions.Admin,
		Push:  &apiPermissions.Push,
		Pull:  &apiPermissions.Pull,
	}
}

// ToGraphExternalTracker convert an api.ExternalTracker to a graphql ExternalTracker
func ToGraphExternalTracker(apiExternalTracker *structs.ExternalTracker) *model.ExternalTracker {
	return &model.ExternalTracker{
		ExternalTrackerURL:    &apiExternalTracker.ExternalTrackerURL,
		ExternalTrackerFormat: &apiExternalTracker.ExternalTrackerFormat,
		ExternalTrackerStyle:  &apiExternalTracker.ExternalTrackerStyle,
	}
}

// ToGraphInternalTracker convert an api.InternalTracker to a graphql InternalTracker
func ToGraphInternalTracker(apiInternalWiki *structs.InternalTracker) *model.InternalTracker {
	return &model.InternalTracker{
		EnableTimeTracker:                &apiInternalWiki.EnableTimeTracker,
		AllowOnlyContributorsToTrackTime: &apiInternalWiki.AllowOnlyContributorsToTrackTime,
		EnableIssueDependencies:          &apiInternalWiki.EnableIssueDependencies,
	}
}

// ToGraphExternalWiki convert an api.ExternalWiki to a graphql ExternalWiki
func ToGraphExternalWiki(apiExternalWiki *structs.ExternalWiki) *model.ExternalWiki {
	return &model.ExternalWiki{ExternalWikiURL: &apiExternalWiki.ExternalWikiURL}
}

// ToGraphUser convert a models.User to a graphql User
func ToGraphUser(user *models.User, signed, authed bool) *model.User {
	apiUser := ToUser(user, signed, authed)
	return innerToGraphUser(apiUser)
}

func innerToGraphUser(apiUser *structs.User) *model.User {
	return &model.User{
		ID:        ToGraphId("user", apiUser.ID),
		RestAPIID: &apiUser.ID,
		Username:  &apiUser.UserName,
		FullName:  &apiUser.FullName,
		Email:     &apiUser.Email,
		AvatarURL: &apiUser.AvatarURL,
		Language:  &apiUser.Language,
		IsAdmin:   &apiUser.IsAdmin,
		LastLogin: &apiUser.LastLogin,
		Created:   &apiUser.Created,
	}
}

// ToGraphBranch convert a git.Commit and git.Branch to an graphql Branch
func ToGraphBranch(repo *models.Repository, b *git.Branch, c *git.Commit, bp *models.ProtectedBranch, user *models.User, isRepoAdmin bool) (*model.Branch, error) {
	apiBranch, err := ToBranch(repo, b, c, bp, user, isRepoAdmin)
	if err != nil {
		return nil, err
	}
	return &model.Branch{
		Name:                          &apiBranch.Name,
		Commit:                        ToGraphPayloadCommit(apiBranch.Commit),
		Protected:                     &apiBranch.Protected,
		RequiredApprovals:             &apiBranch.RequiredApprovals,
		EnableStatusCheck:             &apiBranch.EnableStatusCheck,
		StatusCheckContexts:           &apiBranch.StatusCheckContexts,
		UserCanPush:                   &apiBranch.UserCanMerge,
		UserCanMerge:                  &apiBranch.UserCanPush,
		EffectiveBranchProtectionName: &apiBranch.EffectiveBranchProtectionName,
	}, nil
}

// ToGraphPayloadCommit convert a api.PayloadCommit to a graphql PayloadCommit
func ToGraphPayloadCommit(apiPayloadCommit *structs.PayloadCommit) *model.PayloadCommit {
	return &model.PayloadCommit{
		ID:           &apiPayloadCommit.ID,
		Message:      &apiPayloadCommit.Message,
		URL:          &apiPayloadCommit.URL,
		Author:       ToGraphPayloadUser(apiPayloadCommit.Author),
		Committer:    ToGraphPayloadUser(apiPayloadCommit.Committer),
		Verification: ToGraphPayloadCommitVerification(apiPayloadCommit.Verification),
		Timestamp:    &apiPayloadCommit.Timestamp,
		Added:        &apiPayloadCommit.Added,
		Removed:      &apiPayloadCommit.Removed,
		Modified:     &apiPayloadCommit.Modified,
	}
}

// ToGraphPayloadUser convert an api.PayloadUser to a graphql PayloadUser
func ToGraphPayloadUser(apiPayloadUser *structs.PayloadUser) *model.PayloadUser {
	return &model.PayloadUser{
		Name:     &apiPayloadUser.Name,
		Email:    &apiPayloadUser.Email,
		UserName: &apiPayloadUser.UserName,
	}
}

// ToGraphPayloadCommitVerification convert an api.PayloadCommitVerification to a graphql PayloadCommitVerification
func ToGraphPayloadCommitVerification(apiPayloadCommitVerification *structs.PayloadCommitVerification) *model.PayloadCommitVerification {
	return &model.PayloadCommitVerification{
		Verified:  &apiPayloadCommitVerification.Verified,
		Reason:    &apiPayloadCommitVerification.Reason,
		Signature: &apiPayloadCommitVerification.Signature,
		Signer:    ToGraphPayloadUser(apiPayloadCommitVerification.Signer),
		Payload:   &apiPayloadCommitVerification.Payload,
	}
}

// ToGraphId returns an encoded ID from the given typename and ID that is unique.
func ToGraphId(typename string, ID int64) string {
	str := fmt.Sprintf("%v:%v", typename, ID)
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// FromGraphID takes a string created by ToGraphId and returns a struct with the typename and ID
func FromGraphID(graphID string) *model.GraphID {
	//adapted from https://github.com/graphql-go/relay/
	stringID := ""
	bytes, err := base64.StdEncoding.DecodeString(graphID)
	if err == nil {
		stringID = string(bytes)
	}
	idParts := strings.Split(stringID, ":")
	if len(idParts) < 2 {
		return nil
	}
	return &model.GraphID{
		Typename: idParts[0],
		ID:       idParts[1],
	}
}
