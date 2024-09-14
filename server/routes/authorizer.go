package routes

import "ideashare/models"

type Authorizer[T models.BaseModel] interface {
	CanCreate(user *models.User, model T) bool
	CanUpdate(user *models.User, model T) bool
	CanDelete(user *models.User, model T) bool
	CanReadOne(user *models.User, model T) bool
	CanReadAll(user *models.User) bool
}

type ModelAuthorizer = func(user *models.User, model models.BaseModel) bool
type UserAuthorizer = func(user *models.User) bool

var AllowAllForUser = func(user *models.User) bool {
	return true
}

var AllowAllForModel = func(user *models.User, model models.BaseModel) bool {
	return true
}

type AuthorizeOverrides[T models.BaseModel] struct {
	Create  *ModelAuthorizer
	Update  *ModelAuthorizer
	Delete  *ModelAuthorizer
	ReadOne *ModelAuthorizer
	ReadAll *UserAuthorizer
}

func hasAdminRole(user *models.User) bool {
	for _, role := range user.Roles {
		if role.Role == models.SiteAdmin || role.Role == models.IdeaAdmin {
			return true
		}
	}
	return false
}

func isSiteAdmin(user *models.User) bool {
	for _, role := range user.Roles {
		if role.Role == models.SiteAdmin {
			return true
		}
	}
	return false
}

func isAdminOrOwner(user *models.User, model models.BaseModel) bool {
	if hasAdminRole(user) {
		return true
	}
	if model.GetCreatedBy() == user.GetID() {
		return true
	}
	return false
}

type SiteAdminAuthorizer[T models.BaseModel] struct{}

func NewSiteAdminAuthorizer[T models.BaseModel]() *SiteAdminAuthorizer[T] {
	return &SiteAdminAuthorizer[T]{}
}

func (s *SiteAdminAuthorizer[T]) CanCreate(user *models.User, _ T) bool {
	return isSiteAdmin(user)
}

func (s *SiteAdminAuthorizer[T]) CanUpdate(user *models.User, _ T) bool {
	return isSiteAdmin(user)
}

func (s *SiteAdminAuthorizer[T]) CanDelete(user *models.User, _ T) bool {
	return isSiteAdmin(user)
}

func (s *SiteAdminAuthorizer[T]) CanReadOne(user *models.User, _ T) bool {
	return isSiteAdmin(user)
}

func (s *SiteAdminAuthorizer[T]) CanReadAll(user *models.User) bool {
	return isSiteAdmin(user)
}

type OwnerOrAdminAuthorizer[T models.BaseModel] struct {
	overrides *AuthorizeOverrides[T]
}

func NewOwnerOrAdminAuthorizer[T models.BaseModel](overrides *AuthorizeOverrides[T]) *OwnerOrAdminAuthorizer[T] {
	return &OwnerOrAdminAuthorizer[T]{overrides: overrides}
}

func (o *OwnerOrAdminAuthorizer[T]) CanCreate(user *models.User, model T) bool {
	if o.overrides != nil && o.overrides.Create != nil {
		return (*o.overrides.Create)(user, model)
	}
	return isAdminOrOwner(user, model)
}
func (o *OwnerOrAdminAuthorizer[T]) CanUpdate(user *models.User, model T) bool {
	if o.overrides != nil && o.overrides.Update != nil {
		return (*o.overrides.Update)(user, model)
	}
	return isAdminOrOwner(user, model)
}
func (o *OwnerOrAdminAuthorizer[T]) CanDelete(user *models.User, model T) bool {
	if o.overrides != nil && o.overrides.Delete != nil {
		return (*o.overrides.Delete)(user, model)
	}
	return isAdminOrOwner(user, model)
}
func (o *OwnerOrAdminAuthorizer[T]) CanReadOne(user *models.User, model T) bool {
	if o.overrides != nil && o.overrides.ReadOne != nil {
		return (*o.overrides.ReadOne)(user, model)
	}
	return isAdminOrOwner(user, model)
}
func (o *OwnerOrAdminAuthorizer[T]) CanReadAll(user *models.User) bool {
	if o.overrides != nil && o.overrides.ReadAll != nil {
		return (*o.overrides.ReadAll)(user)
	}
	return hasAdminRole(user)
}
