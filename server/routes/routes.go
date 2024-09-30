package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ideashare/auth"
	"ideashare/config"
	"ideashare/models"
	"net/url"
	"strings"
	"time"
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

func setAuthCookie(c *fiber.Ctx, key string, token string, expires time.Time) {
	reqHeaders := c.GetReqHeaders()
	host := "localhost"
	if _, ok := reqHeaders["Host"]; ok {
		host = reqHeaders["Host"][0]
	}
	c.Cookie(&fiber.Cookie{
		Name:     key,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Expires:  expires,
		Secure:   !strings.Contains(host, "localhost") && !strings.Contains(host, "127.0.0.1"),
		SameSite: "strict",
	})
}

func ConfigureRoutes(app *fiber.App, container *config.AppContainer) {
	app.Static("/", "./public")
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Get("/api/auth/login", func(c *fiber.Ctx) error {
		nonce := uuid.New().String()
		loginUrl := strings.Replace(container.OAuth2Config.AuthCodeURL(nonce), "response_type=code", "response_type=id_token", 1) + "&response_mode=form_post&nonce=" + nonce
		return c.Redirect(loginUrl)
	})
	app.Post("/api/auth/authorize", func(c *fiber.Ctx) error {
		params, err := url.ParseQuery(string(c.Body()))
		if err != nil {
			return c.SendStatus(500)
		}
		idToken, ok := params["id_token"]
		if !ok || len(idToken) != 1 {
			return c.SendStatus(500)
		}
		rawIdToken := idToken[0]
		parsedIdToken, err := container.IdTokenVerifier.Verify(c.Context(), rawIdToken)
		if err != nil {
			return c.SendStatus(401)
		}

		setAuthCookie(c, auth.IdeaShareIDToken, rawIdToken, parsedIdToken.Expiry)
		return c.Redirect("/")
	})
	authMiddleware := func(c *fiber.Ctx) error {
		idToken := c.Cookies(auth.IdeaShareIDToken)
		if idToken == "" {
			return c.SendStatus(401)
		}
		parsedToken, err := container.IdTokenVerifier.Verify(c.Context(), idToken)
		if err != nil {
			return c.SendStatus(401)
		}
		user := &models.User{}
		res := container.Db.Model(user).Preload("Roles").Where("external_id = ?", parsedToken.Subject).First(user)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				claims := auth.OidcUserClaims{}
				err = parsedToken.Claims(&claims)
				if err != nil {
					return c.SendStatus(500)
				}
				user.ExternalID = parsedToken.Subject
				user.FirstName = claims.GivenName
				user.LastName = claims.FamilyName
				user.Email = claims.Email
				container.Db.Create(user)
				container.Db.Model(user).Preload("Roles").Where("external_id = ?", parsedToken.Subject).First(user)
			}
		}
		c.Locals("user", user)
		return c.Next()
	}
	api := app.Group("/api").Use(authMiddleware)
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
