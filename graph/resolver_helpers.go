package graph

import (
	"errors"

	giteaCtx "code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/routers/api/v1/utils"
)

// AuthorizeRepository returns error if not authorized to resolve repostiory, nil otherwise
func AuthorizeRepository(ctx *giteaCtx.APIContext) error {
	if !utils.IsAnyRepoReader(ctx) {
		return errors.New("Must have permission to read repository")
	}
	return nil
}
