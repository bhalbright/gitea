package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"code.gitea.io/gitea/graph/generated"
	"code.gitea.io/gitea/graph/model"
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/convert"
)

func (r *queryResolver) Repository(ctx context.Context, owner string, name string) (*model.Repository, error) {
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

	err = AuthorizeRepository(r.giteaApiContext)
	if err != nil {
		return nil, err
	}

	gqlRepo := convert.ToGraphRepository(repo, models.AccessModeRead)
	return gqlRepo, nil

	return nil, errors.New("both owner and repository name must be provided")
}

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *repositoryResolver) Collaborators(ctx context.Context, obj *model.Repository, first *int, after *string, last *int, before *string) (*model.UserConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Repository returns generated.RepositoryResolver implementation.
func (r *Resolver) Repository() generated.RepositoryResolver { return &repositoryResolver{r} }

type queryResolver struct{ *Resolver }
type repositoryResolver struct{ *Resolver }
