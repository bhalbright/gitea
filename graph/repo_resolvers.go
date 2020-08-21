package graph

import (
	"code.gitea.io/gitea/modules/setting"
	"context"
	"errors"
	"strconv"
	"strings"

	"code.gitea.io/gitea/graph/model"
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/convert"
)

func (r *queryResolver) resolveRepository(ctx context.Context, owner string, name string) (*model.Repository, error) {
	var (
		repoOwner *models.User
		err       error
	)

	// Check if the user is the same as the repository owner.
	if r.giteaApiContext.IsSigned && r.giteaApiContext.User.LowerName == strings.ToLower(owner) {
		repoOwner = r.giteaApiContext.User
	} else {
		repoOwner, err = models.GetUserByName(owner)
		if err != nil {
			return nil, err
		}
	}
	r.giteaApiContext.Repo.Owner = repoOwner

	// Get repository.
	repo, err := models.GetRepositoryByName(repoOwner.ID, name)
	if err != nil {
		return nil, err
	}

	repo.Owner = repoOwner
	r.giteaApiContext.Repo.Repository = repo

	r.giteaApiContext.Repo.Permission, err = models.GetUserRepoPermission(repo, r.giteaApiContext.User)
	if err != nil {
		return nil, err
	}

	if !r.giteaApiContext.Repo.HasAccess() {
		return nil, errors.New("repo not found")
	}

	err = authorizeRepository(r.giteaApiContext)
	if err != nil {
		return nil, err
	}

	gqlRepo := convert.ToGraphRepository(repo, models.AccessModeRead)
	return gqlRepo, nil
}

func (r *queryResolver) resolveRepositoryById(ctx context.Context, id string) (*model.Repository, error) {
	var err error

	internalID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("Unable to parse id")
	}

	// Get repository.
	repo, err := models.GetRepositoryByID(internalID)
	if err != nil {
		return nil, err
	}
	repo.MustOwner()

	r.giteaApiContext.Repo.Owner = repo.Owner
	r.giteaApiContext.Repo.Repository = repo
	r.giteaApiContext.Repo.Permission, err = models.GetUserRepoPermission(repo, r.giteaApiContext.User)
	if err != nil {
		return nil, err
	}

	if !r.giteaApiContext.Repo.HasAccess() {
		return nil, errors.New("repo not found")
	}

	err = authorizeRepository(r.giteaApiContext)
	if err != nil {
		return nil, err
	}

	gqlRepo := convert.ToGraphRepository(repo, models.AccessModeRead)
	return gqlRepo, nil
}

func (r *repositoryResolver) resolveCollaborators(ctx context.Context, obj *model.Repository, first *int, after *string, last *int, before *string) (*model.UserConnection, error) {
	err := authorizeCollaborators(r.giteaApiContext)
	if err != nil {
		return nil, err
	}
	totalSize, err := r.giteaApiContext.Repo.Repository.CountCollaborators()
	if err != nil {
		return nil, err
	}

	listOptions := GetListOptions(totalSize, first, after, last, before, setting.API.MaxResponseItems)
	collaborators, err := r.giteaApiContext.Repo.Repository.GetCollaborators(listOptions)
	if err != nil {
		return nil, err
	}
	users := []*model.User{}
	for _, collaborator := range collaborators {
		user := convert.ToGraphUser(collaborator.User, r.giteaApiContext.IsSigned,
			r.giteaApiContext.User != nil && r.giteaApiContext.User.IsAdmin)
		users = append(users, user)
	}

	startPosition := listOptions.Offset + 1
	cursorPosition := startPosition
	edges := []*model.UserEdge{}
	nodes := []*model.User{}
	for _, user := range users {
		edges = append(edges, &model.UserEdge{
			Cursor: offsetToCursor(cursorPosition),
			Node:   user,
		})
		cursorPosition++
		nodes = append(nodes, user)
	}

	var firstEdgeCursor, lastEdgeCursor string
	if len(edges) > 0 {
		firstEdgeCursor = edges[0].Cursor
		lastEdgeCursor = edges[len(edges)-1:][0].Cursor
	}

	conn := &model.UserConnection{
		TotalCount: &totalSize,
		Edges:      edges,
	}

	conn.PageInfo = &model.PageInfo{
		StartCursor:     &firstEdgeCursor,
		EndCursor:       &lastEdgeCursor,
		HasPreviousPage: startPosition > 1,
		HasNextPage:     (startPosition-1)+len(users) < totalSize,
	}
	return conn, nil
}

func (r *queryResolver) resolveUserById(goCtx context.Context, id string) (*model.User, error) {
	internalID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("Unable to parse id")
	}
	user, err := models.GetUserByID(internalID)
	if err != nil {
		return nil, errors.New("Unable to find user")
	}
	return convert.ToGraphUser(user, r.giteaApiContext.IsSigned,
		r.giteaApiContext.User != nil && r.giteaApiContext.User.IsAdmin), nil
}

/*
// BranchesResolver resolves the branches of a repository
func BranchesResolver(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context.Value(contextKeyType("giteaApiContext")).(*giteaCtx.APIContext)
	err := authorizeBranches(ctx)
	if err != nil {
		return nil, err
	}

	branches, err := repo_module.GetBranches(ctx.Repo.Repository)
	if err != nil {
		return nil, err
	}

	apiBranches := []interface{}{}
	for i := range branches {
		c, err := branches[i].GetCommit()
		if err != nil {
			return nil, err
		}
		branchProtection, err := ctx.Repo.Repository.GetBranchProtection(branches[i].Name)
		log.Info("branch name %d: %s", i, branches[i].Name)
		if err != nil {
			return nil, err
		}
		apiBranch, err := convert.ToBranch(ctx.Repo.Repository, branches[i], c, branchProtection, ctx.User, ctx.Repo.IsAdmin())
		if err != nil {
			return nil, err
		}
		apiBranches = append(apiBranches, apiBranch)
	}

	return apiBranches, nil
}



 */
