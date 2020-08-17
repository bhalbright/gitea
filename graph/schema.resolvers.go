package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"code.gitea.io/gitea/graph/generated"
	"code.gitea.io/gitea/graph/model"
	"code.gitea.io/gitea/modules/convert"
)

func (r *queryResolver) Repository(ctx context.Context, owner string, name string) (*model.Repository, error) {
	return r.resolveRepository(ctx, owner, name)
}

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	//TODO invalid id is blowing up
	resolvedID := convert.FromGraphID(id)
	switch resolvedID.Typename {
	case "repository":
		return r.resolveRepositoryById(ctx, resolvedID.ID)
	case "user":
		return r.resolveUserById(ctx, resolvedID.ID)
	default:
		return nil, errors.New("Unknown node type")
	}
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
