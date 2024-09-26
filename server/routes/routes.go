package routes

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ideashare/config"
	"ideashare/models"
)

func AllowAllReadAuthorizer[T models.BaseModel](ctx *CrudderCtx[T]) bool {
	if ctx.Method == ReadAll || ctx.Method == ReadOne {
		return true
	}
	return OwnerOrAdminAuthorizer(ctx)
}

func SelfOrSiteAdminAuthorizer(ctx *CrudderCtx[*models.User]) bool {
	if ctx.Model != nil && ctx.Model.ID == ctx.User.ID {
		return true
	}
	return SiteAdminAuthorizer(ctx)
}

func ConfigureRoutes(app *fiber.App, container *config.AppContainer) {
	app.Static("/", "./public")
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Get("/api/auth/login", func(c *fiber.Ctx) error {
		return c.Redirect(container.OAuth2Config.AuthCodeURL(uuid.New().String()))
	})
	var verifier = container.OIDCProvider.Verifier(&oidc.Config{ClientID: container.OAuth2Config.ClientID})
	app.Get("/api/auth/authorize", func(c *fiber.Ctx) error {
		oauth2Token, err := container.OAuth2Config.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			// handle error
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			// handle missing token
		}

		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(c.Context(), rawIDToken)
		if err != nil {
			// handle error
		}

		// Extract custom claims
		var claims struct {
			Email    string `json:"email"`
			Verified bool   `json:"email_verified"`
		}
		if err := idToken.Claims(&claims); err != nil {
			// handle error
		}
		return c.SendString("OK")
	})
	api := app.Group("/api")
	api.Get("/user/current", AppRoute(container, func(container *config.AppContainer, c *fiber.Ctx) (interface{}, error) {
		return c.Locals("user"), nil
	}))

	RegisterCrudder(&CrudderConfig[*models.UserSetting]{
		router:        api,
		basePath:      "/user/:userID/setting",
		container:     container,
		newEmptyModel: func() *models.UserSetting { return &models.UserSetting{} },
		authorizer: func(ctx *CrudderCtx[*models.UserSetting]) bool {
			userID, err := ctx.ReqCtx.ParamsInt("userID")
			if err != nil {
				return false
			}
			if ctx.RequestBody != nil && ctx.RequestBody.UserId != userID {
				return false
			}
			return isSiteAdmin(ctx.User) || ctx.User.ID == userID
		},
		applyFilter: func(ctx *CrudderCtx[*models.UserSetting], query *gorm.DB) (*gorm.DB, error) {
			userID, err := ctx.ReqCtx.ParamsInt("userID")
			if err != nil {
				return nil, fiber.ErrBadRequest
			}
			return query.Where("user_id = ?", userID), nil
		},
	})

	RegisterCrudder(&CrudderConfig[*models.User]{
		router:        api,
		basePath:      "/user",
		container:     container,
		newEmptyModel: func() *models.User { return &models.User{} },
		authorizer:    SelfOrSiteAdminAuthorizer,
	})

	RegisterCrudder(&CrudderConfig[*models.SiteSetting]{
		router:        api,
		basePath:      "/setting",
		container:     container,
		newEmptyModel: func() *models.SiteSetting { return &models.SiteSetting{} },
		authorizer:    SiteAdminAuthorizer[*models.SiteSetting],
	})

	RegisterCrudder(&CrudderConfig[*models.IdeaLike]{
		router:        api,
		basePath:      "/idea/:ideaID/like",
		container:     container,
		newEmptyModel: func() *models.IdeaLike { return &models.IdeaLike{} },
		authorizer: func(ctx *CrudderCtx[*models.IdeaLike]) bool {
			if isSiteAdmin(ctx.User) {
				return true
			}
			if ctx.Method == ReadAll || ctx.Method == ReadOne || ctx.Method == Update {
				return false
			}
			if ctx.RequestBody != nil && ctx.RequestBody.UserID != ctx.User.ID {
				return false
			}
			paramsInt, err := ctx.ReqCtx.ParamsInt("ideaID")
			if err != nil {
				return false
			}
			if ctx.RequestBody.IdeaID != paramsInt {
				return false
			}
			return OwnerOrAdminAuthorizer(ctx)
		},
		applyFilter: func(ctx *CrudderCtx[*models.IdeaLike], query *gorm.DB) (*gorm.DB, error) {
			ideaId, err := ctx.ReqCtx.ParamsInt("ideaID")
			if err != nil {
				return nil, fiber.ErrBadRequest
			}
			return query.Where("idea_id = ?", ideaId), nil
		},
	})

	RegisterCrudder(&CrudderConfig[*models.IdeaComment]{
		router:        api,
		basePath:      "/idea/:ideaID/comment",
		container:     container,
		newEmptyModel: func() *models.IdeaComment { return &models.IdeaComment{} },
		authorizer: func(ctx *CrudderCtx[*models.IdeaComment]) bool {
			paramsInt, err := ctx.ReqCtx.ParamsInt("ideaID")
			if err != nil {
				return false
			}
			if (ctx.Method == Create || ctx.Method == Update) && ctx.RequestBody.IdeaID != paramsInt {
				return false
			}
			return AllowAllReadAuthorizer(ctx)
		},
		applyFilter: func(ctx *CrudderCtx[*models.IdeaComment], query *gorm.DB) (*gorm.DB, error) {
			paramsInt, err := ctx.ReqCtx.ParamsInt("ideaID")
			if err != nil {
				return nil, fiber.ErrBadRequest
			}
			return query.Where("idea_id = ?", paramsInt), nil
		},
	})

	RegisterCrudder(&CrudderConfig[*models.Idea]{
		router:          api,
		basePath:        "/idea",
		container:       container,
		newEmptyModel:   func() *models.Idea { return &models.Idea{} },
		authorizer:      AllowAllReadAuthorizer[*models.Idea],
		readOnePreloads: []string{"Comments"},
		eventHandlers: []*CrudderEventHandler[*models.Idea]{{
			Handles: AfterLoad,
			Handle: func(event *CrudderEvent[*models.Idea]) error {
				if event.Ctx.Method != ReadOne {
					return nil
				}
				model := event.Ctx.Model
				user := event.Ctx.User
				var count int64
				container.Db.Model(&models.IdeaLike{}).Where("idea_id = ?", model.ID).Count(&count)
				model.Likes = int(count)
				container.Db.Model(&models.IdeaLike{}).Where("created_by = ?", user.ID).Count(&count)
				return nil
			},
		}},
	})

}
