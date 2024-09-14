package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"ideashare/config"
	"ideashare/models"
	"time"
)

var dummyUser = &models.User{
	Base: models.Base{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	},
	ExternalID: "1",
	FirstName:  "Admin",
	LastName:   "Adminovich",
	Roles:      []*models.UserRole{{UserID: 1, Role: models.SiteAdmin}},
}

func AppRouteWithBody[T any](container *config.AppContainer, newBody func() T, handler func(container *config.AppContainer, b T, c *fiber.Ctx) (interface{}, error)) func(c *fiber.Ctx) error {
	return AppRoute(container, func(container *config.AppContainer, c *fiber.Ctx) (interface{}, error) {
		reqBody := newBody()
		if err := c.BodyParser(reqBody); err != nil {
			return nil, c.SendStatus(400)
		}
		return handler(container, reqBody, c)
	})

}

func AppRoute(container *config.AppContainer, handler func(container *config.AppContainer, c *fiber.Ctx) (interface{}, error)) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		//TODO authenticate and lookup user
		c.Locals("user", dummyUser)
		res, err := handler(container, c)
		if err != nil || res == nil {
			return err
		}
		c.Append("Content-Type", "application/json")
		if err = c.JSON(res); err != nil {
			return c.SendStatus(500)
		}
		return nil
	}
}

type Crudder[T models.BaseModel] struct {
	container     *config.AppContainer
	newEmptyModel func() T
	authorizer    Authorizer[T]
	maxPageSize   int
	preloads      []string
	eventHandlers map[EventType][]*CrudderEventHandler[T]
}

type EventAction string

const (
	AfterLoad EventAction = "afterLoad"
)

type CrudMethod string

const (
	Create  CrudMethod = "create"
	ReadOne CrudMethod = "readOne"
	ReadAll CrudMethod = "readAll"
	Update  CrudMethod = "update"
	Delete  CrudMethod = "delete"
)

type EventType struct {
	Method CrudMethod
	Action EventAction
}

type CrudderEventHandler[T any] struct {
	Handles EventType
	Handle  func(event CrudderEvent[T]) error
}

type CrudderEvent[T any] struct {
	Type         EventType
	CanTransform bool
	Context      context.Context
	Transformed  *T
}

func RegisterCrudder[T models.BaseModel](
	router fiber.Router,
	basePath string,
	container *config.AppContainer,
	newEmptyModel func() T,
	authorizer Authorizer[T],
	maxPageSize int,
	preloads []string,
	eventHandlers ...*CrudderEventHandler[T],
) *Crudder[T] {
	crudder := &Crudder[T]{
		container:     container,
		newEmptyModel: newEmptyModel,
		authorizer:    authorizer,
		maxPageSize:   maxPageSize,
		preloads:      preloads,
		eventHandlers: make(map[EventType][]*CrudderEventHandler[T]),
	}
	for _, handler := range eventHandlers {
		if _, ok := crudder.eventHandlers[handler.Handles]; !ok {
			crudder.eventHandlers[handler.Handles] = []*CrudderEventHandler[T]{handler}
		} else {
			crudder.eventHandlers[handler.Handles] = append(crudder.eventHandlers[handler.Handles], handler)
		}
	}
	crudder.registerRoutes(basePath, router)
	return crudder
}

func (c *Crudder[T]) registerRoutes(basePath string, router fiber.Router) {
	router.Get(basePath, c.ReadAll())
	router.Post(basePath, c.Create())
	router.Get(basePath+"/:id", c.ReadOneById())
	router.Put(basePath+"/:id", c.UpdateById())
	router.Delete(basePath+"/:id", c.DeleteById())
}

type ContextValue struct {
	Key string
	Val interface{}
}

const CurrentRequestContextKey = "requestContext"
const CurrentModelContextKey = "modelContext"
const CurrentUserContextKey = "userContext"

func eventContext(ctx *fiber.Ctx, values ...ContextValue) context.Context {
	c := context.WithValue(context.Background(), CurrentRequestContextKey, ctx)
	for _, v := range values {
		c = context.WithValue(c, v.Key, v.Val)
	}
	return c
}

func (c *Crudder[T]) fireEvent(method CrudMethod, action EventAction, ctx context.Context, canTransform bool) (interface{}, error) {
	eventType := EventType{
		Method: method,
		Action: action,
	}
	event := CrudderEvent[T]{
		Type:         eventType,
		CanTransform: canTransform,
		Context:      ctx,
	}
	for _, handler := range c.eventHandlers[eventType] {
		if err := handler.Handle(event); err != nil {
			return nil, err
		}
	}
	if event.CanTransform && event.Transformed != nil {
		return event.Transformed, nil
	}
	return nil, nil
}

func (c *Crudder[T]) Create() func(c *fiber.Ctx) error {
	return AppRouteWithBody(c.container, c.newEmptyModel, func(container *config.AppContainer, incoming T, ctx *fiber.Ctx) (interface{}, error) {
		user := ctx.Locals("user").(*models.User)
		if !c.authorizer.CanCreate(user, incoming) {
			return nil, ctx.SendStatus(403)
		}
		incoming.SetCreatedBy(user.ID)
		result := container.Db.Create(incoming)
		if result.Error != nil {
			return nil, result.Error
		}
		ctx.Status(201)
		return incoming, nil
	})
}

func (c *Crudder[T]) ReadAll() func(ctx *fiber.Ctx) error {
	return AppRoute(c.container, func(container *config.AppContainer, ctx *fiber.Ctx) (interface{}, error) {
		user := ctx.Locals("user").(*models.User)
		if !c.authorizer.CanReadAll(user) {
			return nil, ctx.SendStatus(403)
		}
		size := ctx.QueryInt("size", 10)
		page := ctx.QueryInt("page", 1)
		var results []T
		result := container.Db.Offset((page - 1) * size).Limit(size).Find(&results)
		if result.Error != nil {
			return nil, result.Error
		}
		return results, nil
	})
}

func (c *Crudder[T]) ReadOneById() func(c *fiber.Ctx) error {
	return AppRoute(c.container, func(container *config.AppContainer, ctx *fiber.Ctx) (interface{}, error) {
		user := ctx.Locals("user").(*models.User)
		found := c.newEmptyModel()
		tx := container.Db
		for _, preload := range c.preloads {
			tx = tx.Preload(preload)
		}
		tx.Find(found, ctx.Params("id"))
		if !c.authorizer.CanReadOne(user, found) {
			return nil, ctx.SendStatus(403)
		}
		if found.GetID() == 0 {
			return nil, ctx.SendStatus(404)
		}
		handled, err := c.fireEvent(ReadOne, AfterLoad, eventContext(ctx,
			ContextValue{CurrentModelContextKey, found},
			ContextValue{CurrentUserContextKey, user},
		), true)
		if err != nil {
			return nil, err
		}
		if handled != nil {
			return handled, nil
		}
		return found, nil
	})
}

func (c *Crudder[T]) UpdateById() func(c *fiber.Ctx) error {
	return AppRouteWithBody(c.container, c.newEmptyModel, func(container *config.AppContainer, incoming T, ctx *fiber.Ctx) (interface{}, error) {
		found := c.newEmptyModel()
		container.Db.Find(found, ctx.Params("id"))
		if !c.authorizer.CanUpdate(ctx.Locals("user").(*models.User), found) {
			return nil, ctx.SendStatus(403)
		}
		if found.GetID() == 0 {
			return nil, ctx.SendStatus(404)
		}
		result := container.Db.Save(&incoming)
		if result.Error != nil {
			return nil, result.Error
		}
		return incoming, nil
	})

}

func (c *Crudder[T]) DeleteById() func(c *fiber.Ctx) error {
	return AppRouteWithBody(c.container, c.newEmptyModel, func(container *config.AppContainer, incoming T, ctx *fiber.Ctx) (interface{}, error) {
		found := c.newEmptyModel()
		container.Db.Find(found, ctx.Params("id"))
		if !c.authorizer.CanDelete(ctx.Locals("user").(*models.User), found) {
			return nil, ctx.SendStatus(403)
		}
		result := container.Db.Delete(found, ctx.Params("id"))
		if result.Error != nil {
			return nil, result.Error
		}
		ctx.Status(204)
		return nil, nil
	})
}
