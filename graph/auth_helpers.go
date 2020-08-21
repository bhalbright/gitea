package graph

import (
	"code.gitea.io/gitea/models"
	"errors"

	giteaCtx "code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/routers/api/v1/utils"
)

func authorizeRepository(ctx *giteaCtx.APIContext) error {
	if !utils.IsAnyRepoReader(ctx) {
		return errors.New("Must have permission to read repository")
	}
	return nil
}

func authorizeCollaborators(ctx *giteaCtx.APIContext) error {
	if _, found :=  ctx.Data["IsApiToken"]; !found {
		return errors.New("Api token missing or invalid")
	}
	if !utils.IsAnyRepoReader(ctx) {
		return errors.New("Must have permission to read repository")
	}
	return nil
}

func authorizeBranches(ctx *giteaCtx.APIContext) error {
	if !utils.IsRepoReader(ctx, models.UnitTypeCode) {
		return errors.New("Must have read permission or be a repo or site admin")
	}
	return nil
}
