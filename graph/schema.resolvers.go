package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/graph/generated"
	"code.gitea.io/gitea/graph/model"
)

func (r *queryResolver) Repository(ctx context.Context, owner string, name string) (*model.Repository, error) {
	return r.resolveRepository(ctx, owner, name)
}

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *repositoryResolver) Collaborators(ctx context.Context, obj *model.Repository, first *int, after *string, last *int, before *string) (*model.UserConnection, error) {
	return r.resolveCollaborators(ctx, obj, first, after, last, before)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Repository returns generated.RepositoryResolver implementation.
func (r *Resolver) Repository() generated.RepositoryResolver { return &repositoryResolver{r} }

type queryResolver struct{ *Resolver }
type repositoryResolver struct{ *Resolver }
