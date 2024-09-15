package routes

import (
	"github.com/gofiber/fiber/v2"
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

func SelfOrSiteAdminAuthorizer[T models.BaseModel](ctx *CrudderCtx[T]) bool {
	var zero T
	if ctx.Model != zero && ctx.Model.GetID() == ctx.User.ID {
		return true
	}
	return SiteAdminAuthorizer(ctx)
}

func ConfigureRoutes(app *fiber.App, container *config.AppContainer) {
	app.Static("/", "./public")
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	api := app.Group("/api")
	api.Get("/user/current", AppRoute(container, func(container *config.AppContainer, c *fiber.Ctx) (interface{}, error) {
		return c.Context().UserValue("user"), nil
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
			return query.Where("created_by = ?", ctx.User.ID).Where("user_id = ?", userID), nil
		},
	})

	RegisterCrudder(&CrudderConfig[*models.User]{
		router:        api,
		basePath:      "/user",
		container:     container,
		newEmptyModel: func() *models.User { return &models.User{} },
		authorizer:    SelfOrSiteAdminAuthorizer[*models.User],
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
			if !isSiteAdmin(ctx.User) && (ctx.Method == ReadAll || ctx.Method == ReadOne || ctx.Method == Update) {
				return false
			}
			if ctx.RequestBody != nil && ctx.RequestBody.UserID != ctx.User.ID {
				return false
			}
			paramsInt, err := ctx.ReqCtx.ParamsInt("ideaID")
			if err != nil {
				return false
			}
			if (ctx.Method == Create || ctx.Method == Update) && ctx.RequestBody.IdeaID != paramsInt {
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
