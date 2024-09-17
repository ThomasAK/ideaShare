package routes

import "ideashare/models"

func hasAdminRole(user *models.User) bool {
	if user == nil {
		return false
	}
	for _, role := range user.Roles {
		if role.Role == models.SiteAdmin || role.Role == models.IdeaAdmin {
			return true
		}
	}
	return false
}

func isSiteAdmin(user *models.User) bool {
	if user == nil {
		return false
	}
	for _, role := range user.Roles {
		if role.Role == models.SiteAdmin {
			return true
		}
	}
	return false
}

func isAdminOrOwner[T models.BaseModel](user *models.User, model T) bool {
	if user == nil {
		return false
	}
	if hasAdminRole(user) {
		return true
	}
	if model.GetCreatedBy() == user.GetID() {
		return true
	}
	return false
}

func SiteAdminAuthorizer[T models.BaseModel](ctx *CrudderCtx[T]) bool {
	return isSiteAdmin(ctx.User)
}

func OwnerOrAdminAuthorizer[T models.BaseModel](ctx *CrudderCtx[T]) bool {
	var zero T
	if ctx.Model == zero {
		return hasAdminRole(ctx.User)
	}
	return isAdminOrOwner(ctx.User, ctx.Model)
}
