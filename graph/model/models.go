package model

import "time"

type GraphID struct {
	Typename string
	ID       string
}

type Repository struct {
	ID              string          `json:"id"`
	RestAPIID       *int64          `json:"rest_api_id"`
	Name            *string         `json:"name"`
	FullName        *string         `json:"fullName"`
	CollaboratorsID *int            `json:"collaborators"`
	Owner         *User       `json:"owner"`
	Description   *string      `json:"description"`
	Empty         *bool        `json:"empty"`
	Private       *bool        `json:"private"`
	Fork          *bool        `json:"fork"`
	Template      *bool        `json:"template"`
	Parent        *Repository `json:"parent"`
	Mirror        *bool        `json:"mirror"`
	Size          *int         `json:"size"`
	HTMLURL       *string      `json:"html_url"`
	SSHURL        *string      `json:"ssh_url"`
	CloneURL      *string      `json:"clone_url"`
	OriginalURL   *string      `json:"original_url"`
	Website       *string      `json:"website"`
	Stars         *int         `json:"stars_count"`
	Forks         *int         `json:"forks_count"`
	Watchers      *int         `json:"watchers_count"`
	OpenIssues    *int         `json:"open_issues_count"`
	OpenPulls     *int         `json:"open_pr_counter"`
	Releases      *int         `json:"release_counter"`
	DefaultBranch *string      `json:"default_branch"`
	Archived      *bool        `json:"archived"`
	Created *time.Time `json:"created_at"`
	Updated                   *time.Time        `json:"updated_at"`
	Permissions               *Permission      `json:"permissions,omitempty"`
	HasIssues                 *bool             `json:"has_issues"`
	InternalTracker           *InternalTracker `json:"internal_tracker,omitempty"`
	ExternalTracker           *ExternalTracker `json:"external_tracker,omitempty"`
	HasWiki                   *bool             `json:"has_wiki"`
	ExternalWiki              *ExternalWiki    `json:"external_wiki,omitempty"`
	HasPullRequests           *bool             `json:"has_pull_requests"`
	IgnoreWhitespaceConflicts *bool             `json:"ignore_whitespace_conflicts"`
	AllowMerge                *bool             `json:"allow_merge_commits"`
	AllowRebase               *bool             `json:"allow_rebase"`
	AllowRebaseMerge          *bool             `json:"allow_rebase_explicit"`
	AllowSquash               *bool             `json:"allow_squash_merge"`
	AvatarURL                 *string           `json:"avatar_url"`
	Internal                  *bool             `json:"internal"`
	Branches                  []*Branch        `json:"branches"`
}

func (Repository) IsNode() {}

type User struct {
	ID        string  `json:"id"`
	RestAPIID *int64  `json:"rest_api_id"`
	Username  *string `json:"username"`
	FullName *string `json:"full_name"`
	Email *string `json:"email"`
	AvatarURL *string `json:"avatar_url"`
	Language *string `json:"language"`
	IsAdmin *bool `json:"is_admin"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Created *time.Time `json:"created,omitempty"`
}

func (User) IsNode() {}

type Permission struct {
	Admin *bool `json:"admin"`
	Push  *bool `json:"push"`
	Pull  *bool `json:"pull"`
}

type InternalTracker struct {
	EnableTimeTracker *bool `json:"enable_time_tracker"`
	AllowOnlyContributorsToTrackTime *bool `json:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies *bool `json:"enable_issue_dependencies"`
}

type ExternalTracker struct {
	ExternalTrackerURL *string `json:"external_tracker_url"`
	ExternalTrackerFormat *string `json:"external_tracker_format"`
	ExternalTrackerStyle *string `json:"external_tracker_style"`
}

type ExternalWiki struct {
	ExternalWikiURL *string `json:"external_wiki_url"`
}

type Branch struct {
	Name                          *string         `json:"name"`
	Commit                        *PayloadCommit `json:"commit"`
	Protected                     *bool           `json:"protected"`
	RequiredApprovals             *int64          `json:"required_approvals"`
	EnableStatusCheck             *bool           `json:"enable_status_check"`
	StatusCheckContexts           *[]string       `json:"status_check_contexts"`
	UserCanPush                   *bool           `json:"user_can_push"`
	UserCanMerge                  *bool           `json:"user_can_merge"`
	EffectiveBranchProtectionName *string         `json:"effective_branch_protection_name"`
}

type PayloadCommit struct {
	ID           *string                     `json:"id"`
	Message      *string                     `json:"message"`
	URL          *string                     `json:"url"`
	Author       *PayloadUser               `json:"author"`
	Committer    *PayloadUser               `json:"committer"`
	Verification *PayloadCommitVerification `json:"verification"`
	Timestamp *time.Time `json:"timestamp"`
	Added     *[]string  `json:"added"`
	Removed   *[]string  `json:"removed"`
	Modified  *[]string  `json:"modified"`
}

type PayloadCommitVerification struct {
	Verified  *bool         `json:"verified"`
	Reason    *string       `json:"reason"`
	Signature *string       `json:"signature"`
	Signer    *PayloadUser `json:"signer"`
	Payload   *string       `json:"payload"`
}

type PayloadUser struct {
	Name *string `json:"name"`
	Email    *string `json:"email"`
	UserName *string `json:"username"`
}
