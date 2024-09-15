package routes

import (
	"github.com/google/uuid"
	"ideashare/models"
	"ideashare/tests/testutil"
	"testing"
)

func TestUsers(t *testing.T) {
	TestCrudEndpoints(t, "/user", func() *models.User {
		return &models.User{}
	}, func(user *models.User) {
		user.FirstName = "Jerry"
		user.LastName = "Seinfeld"
	})
	testutil.PrintTableContents("users")
}

func TestUserSettings(t *testing.T) {
	TestCrudEndpoints(t, "/user/1/setting", func() *models.UserSetting {
		return &models.UserSetting{UserId: 1}
	}, func(setting *models.UserSetting) {
		setting.Key = "test"
		setting.Value = "test"
	})
	testutil.PrintTableContents("user_settings")
}

func TestSiteSettings(t *testing.T) {
	TestCrudEndpoints(t, "/setting", func() *models.SiteSetting {
		return &models.SiteSetting{}
	}, func(setting *models.SiteSetting) {
		setting.Key = "test"
		setting.Value = "test"
	})
	testutil.PrintTableContents("site_settings")
}

func TestIdeas(t *testing.T) {
	TestCrudEndpoints(t, "/idea", func() *models.Idea {
		return &models.Idea{Comments: []*models.IdeaComment{}}
	}, func(idea *models.Idea) {
		idea.Title = "test"
		idea.Description = "test"
	})
	testutil.PrintTableContents("ideas")
}

func TestComments(t *testing.T) {
	//ensure at least 1 idea exists
	testutil.Container.Db.Create(&models.Idea{
		Title:       "test idea",
		Description: "test idea",
		Status:      "none",
	})
	TestCrudEndpoints(t, "/idea/1/comment", func() *models.IdeaComment {
		return &models.IdeaComment{IdeaID: 1}
	}, func(comment *models.IdeaComment) {
		comment.Comment = "test"
	})
	testutil.PrintTableContents("idea_comments")
}

func TestLikes(t *testing.T) {

	TestCrudEndpointsWithoutAll(t, "/idea/1/like", func() *models.IdeaLike {
		return &models.IdeaLike{IdeaID: 1, UserID: 1}
	}, func(like *models.IdeaLike) {
		like.UserID = 1
		like.IdeaID = 1
	})
	testReadAll[*models.IdeaLike](t, "/idea/1/like", func(count int) {
		for i := 1; i <= 20; i++ {
			newUUID, _ := uuid.NewUUID()
			testutil.Container.Db.Create(&models.User{
				ExternalID: newUUID.String(),
				FirstName:  "test",
			})
			testutil.Container.Db.Create(&models.IdeaLike{
				HardDeleteModel: models.HardDeleteModel{
					CreatedBy: 1,
				},
				IdeaID: 1,
				UserID: i,
			})
		}
	})
	testutil.PrintTableContents("idea_likes")

}
