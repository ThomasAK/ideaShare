package routes

import (
	"github.com/gofiber/fiber/v2"
	"ideashare/config"
	"ideashare/models"
)

func emptyIdea() *models.Idea               { return &models.Idea{} }
func emptyIdeaComment() *models.IdeaComment { return &models.IdeaComment{} }
func emptyIdeaLike() *models.IdeaLike       { return &models.IdeaLike{} }
func emptyUser() *models.User               { return &models.User{} }
func emptyUserSetting() *models.UserSetting { return &models.UserSetting{} }
func emptySiteSetting() *models.SiteSetting { return &models.SiteSetting{} }

func ConfigureRoutes(app *fiber.App, container *config.AppContainer) {
	app.Static("/", "./public")
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	api := app.Group("/api")
	api.Get("/user/current", AppRoute(container, func(container *config.AppContainer, c *fiber.Ctx) (interface{}, error) {
		return c.Context().UserValue("user"), nil
	}))
	var noPreloads []string
	RegisterCrudder(api, "/user",
		container, emptyUser, NewOwnerOrAdminAuthorizer(&AuthorizeOverrides[*models.User]{}), 100, noPreloads)

	RegisterCrudder(api, "/user/setting",
		container, emptyUserSetting, NewOwnerOrAdminAuthorizer(&AuthorizeOverrides[*models.UserSetting]{}), 100, noPreloads)

	RegisterCrudder(api, "/setting",
		container, emptySiteSetting, NewSiteAdminAuthorizer[*models.SiteSetting](), 100, noPreloads)

	RegisterCrudder(api, "/idea",
		container, emptyIdea, NewOwnerOrAdminAuthorizer(&AuthorizeOverrides[*models.Idea]{
			ReadAll: &AllowAllForUser,
			ReadOne: &AllowAllForModel,
		}), 50, []string{"comments"}, &CrudderEventHandler[*models.Idea]{
			Handles: EventType{ReadOne, AfterLoad},
			Handle: func(event *CrudderEvent[*models.Idea]) error {
				model := event.Model
				user := event.User
				var count int64
				container.Db.Model(&models.IdeaLike{}).Where("idea_id = ?", model.ID).Count(&count)
				model.Likes = int(count)
				container.Db.Model(&models.IdeaLike{}).Where("created_by = ?", user.ID).Count(&count)
				return nil
			},
		})

	RegisterCrudder(api, "/idea/comment",
		container, emptyIdeaComment, NewOwnerOrAdminAuthorizer(&AuthorizeOverrides[*models.IdeaComment]{
			ReadAll: &AllowAllForUser,
			ReadOne: &AllowAllForModel,
		}), 100, noPreloads)

	RegisterCrudder(api, "/idea/like",
		container, emptyIdeaLike, NewOwnerOrAdminAuthorizer(&AuthorizeOverrides[*models.IdeaLike]{}), 100, noPreloads)
}
