package routes

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"ideashare/config"
	"ideashare/models"
	"log/slog"
	"time"
)

var dummyUser = &models.User{
	SoftDeleteModel: models.SoftDeleteModel{
		HardDeleteModel: models.HardDeleteModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: 1,
		},
		DeletedAt: gorm.DeletedAt{},
	},
	ExternalID: "1",
	FirstName:  "Admin",
	LastName:   "Admin",
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
		//container.AuthProvider.authenticate(c)
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
	config        *CrudderConfig[T]
	eventHandlers map[EventType][]*CrudderEventHandler[T]
}

type EventType string

const (
	AfterLoad EventType = "afterLoad"
)

type CrudMethod string

const (
	Create  CrudMethod = "create"
	ReadOne CrudMethod = "readOne"
	ReadAll CrudMethod = "readAll"
	Update  CrudMethod = "update"
	Delete  CrudMethod = "delete"
)

type CrudderEventHandler[T models.BaseModel] struct {
	Handles EventType
	Handle  func(event *CrudderEvent[T]) error
}

type CrudderCtx[T models.BaseModel] struct {
	Method      CrudMethod
	User        *models.User
	RequestBody T
	Model       T
	Rows        []T
	ReqCtx      *fiber.Ctx
}

type CrudderEvent[T models.BaseModel] struct {
	Type EventType
	Ctx  *CrudderCtx[T]
}

type CrudderConfig[T models.BaseModel] struct {
	router          fiber.Router
	basePath        string
	container       *config.AppContainer
	newEmptyModel   func() T
	authorizer      func(ctx *CrudderCtx[T]) bool
	maxPageSize     int
	readOnePreloads []string
	eventHandlers   []*CrudderEventHandler[T]
	applyFilter     func(ctx *CrudderCtx[T], query *gorm.DB) (*gorm.DB, error)
}

func RegisterCrudder[T models.BaseModel](
	config *CrudderConfig[T],
) *Crudder[T] {
	if config.maxPageSize == 0 {
		config.maxPageSize = 50
	}
	if config.applyFilter == nil {
		config.applyFilter = func(ctx *CrudderCtx[T], query *gorm.DB) (*gorm.DB, error) {
			return query, nil
		}
	}
	crudder := &Crudder[T]{
		config:        config,
		eventHandlers: make(map[EventType][]*CrudderEventHandler[T]),
	}
	if config.eventHandlers != nil {
		for _, handler := range config.eventHandlers {
			if _, ok := crudder.eventHandlers[handler.Handles]; !ok {
				crudder.eventHandlers[handler.Handles] = []*CrudderEventHandler[T]{handler}
			} else {
				crudder.eventHandlers[handler.Handles] = append(crudder.eventHandlers[handler.Handles], handler)
			}
		}
	}
	crudder.registerRoutes()
	return crudder
}

func (c *Crudder[T]) registerRoutes() {
	c.config.router.Get(c.config.basePath, c.ReadAll())
	c.config.router.Post(c.config.basePath, c.Create())
	c.config.router.Get(c.config.basePath+"/:id", c.ReadOneById())
	c.config.router.Put(c.config.basePath+"/:id", c.UpdateById())
	c.config.router.Delete(c.config.basePath+"/:id", c.DeleteById())
}

type ContextValue struct {
	Key string
	Val interface{}
}

func (c *Crudder[T]) fireEvent(event *CrudderEvent[T]) error {
	for _, handler := range c.eventHandlers[event.Type] {
		if err := handler.Handle(event); err != nil {
			return err
		}
	}
	return nil
}

func (c *Crudder[T]) Create() func(c *fiber.Ctx) error {
	return AppRouteWithBody(c.config.container, c.config.newEmptyModel, func(container *config.AppContainer, incoming T, ctx *fiber.Ctx) (interface{}, error) {
		user := ctx.Locals("user").(*models.User)
		incoming.SetCreatedBy(user.ID)
		crudderCtx := &CrudderCtx[T]{Method: Create, User: user, Model: incoming, ReqCtx: ctx, RequestBody: incoming}
		if !c.config.authorizer(crudderCtx) {
			return nil, ctx.SendStatus(403)
		}
		result := container.Db.Create(incoming)
		if result.Error != nil {
			return nil, handleDbError(result.Error)
		}
		ctx.Status(201)
		return incoming, nil
	})
}

func (c *Crudder[T]) ReadAll() func(ctx *fiber.Ctx) error {
	return AppRoute(c.config.container, func(container *config.AppContainer, ctx *fiber.Ctx) (interface{}, error) {
		user := ctx.Locals("user").(*models.User)
		size := ctx.QueryInt("size", 10)
		page := ctx.QueryInt("page", 1)
		var results []T
		crudderCtx := &CrudderCtx[T]{Method: ReadAll, User: user, Rows: results, ReqCtx: ctx}
		if !c.config.authorizer(crudderCtx) {
			return nil, ctx.SendStatus(403)
		}
		filter, err := c.config.applyFilter(crudderCtx, container.Db)
		if err != nil {
			return nil, err
		}
		result := filter.Offset((page - 1) * size).Limit(size).Find(&results)
		if result.Error != nil {
			return nil, handleDbError(result.Error)
		}
		err = c.fireEvent(&CrudderEvent[T]{
			AfterLoad,
			crudderCtx,
		})
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

func (c *Crudder[T]) ReadOneById() func(c *fiber.Ctx) error {
	return AppRoute(c.config.container, func(container *config.AppContainer, ctx *fiber.Ctx) (interface{}, error) {
		user := ctx.Locals("user").(*models.User)
		found := c.config.newEmptyModel()
		tx := container.Db
		for _, preload := range c.config.readOnePreloads {
			tx = tx.Preload(preload)
		}
		crudderCtx := &CrudderCtx[T]{Method: ReadOne, User: user, Model: found, ReqCtx: ctx}
		filter, err := c.config.applyFilter(crudderCtx, tx)
		if err != nil {
			return nil, err
		}
		res := filter.Take(found, ctx.Params("id"))
		if res.Error != nil {
			return nil, handleDbError(res.Error)
		}
		if !c.config.authorizer(crudderCtx) {
			return nil, ctx.SendStatus(403)
		}
		err = c.fireEvent(&CrudderEvent[T]{
			AfterLoad,
			crudderCtx,
		})
		if err != nil {
			return nil, err
		}
		return found, nil
	})
}

func (c *Crudder[T]) UpdateById() func(c *fiber.Ctx) error {
	return AppRouteWithBody(c.config.container, c.config.newEmptyModel, func(container *config.AppContainer, incoming T, ctx *fiber.Ctx) (interface{}, error) {
		found := c.config.newEmptyModel()
		user := ctx.Locals("user").(*models.User)
		crudderCtx := &CrudderCtx[T]{Method: Update, User: user, Model: found, ReqCtx: ctx, RequestBody: incoming}
		lookupFilter, err := c.config.applyFilter(crudderCtx, container.Db)
		if err != nil {
			return nil, err
		}
		res := lookupFilter.Take(found, ctx.Params("id"))
		if res.Error != nil {
			return nil, handleDbError(res.Error)
		}
		if !c.config.authorizer(crudderCtx) {
			return nil, ctx.SendStatus(403)
		}
		filter, err := c.config.applyFilter(crudderCtx, container.Db)
		if err != nil {
			return nil, err
		}
		result := filter.Save(&incoming)
		if result.Error != nil {
			return nil, handleDbError(result.Error)
		}
		return incoming, nil
	})

}

func (c *Crudder[T]) DeleteById() func(c *fiber.Ctx) error {
	return AppRoute(c.config.container, func(container *config.AppContainer, ctx *fiber.Ctx) (interface{}, error) {
		found := c.config.newEmptyModel()
		user := ctx.Locals("user").(*models.User)
		crudderCtx := &CrudderCtx[T]{Method: Delete, User: user, Model: found, ReqCtx: ctx}
		filter, err := c.config.applyFilter(crudderCtx, container.Db)
		if err != nil {
			return nil, err
		}
		res := filter.Take(found, ctx.Params("id"))
		if res.Error != nil {
			return nil, handleDbError(res.Error)
		}
		if !c.config.authorizer(crudderCtx) {
			return nil, ctx.SendStatus(403)
		}
		result := container.Db.Delete(found)
		if result.Error != nil {
			return nil, handleDbError(result.Error)
		}
		ctx.Status(204)
		return nil, nil
	})
}

func handleDbError(err error) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.ErrConflict
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrNotFound
	}
	if errors.Is(err, gorm.ErrInvalidData) {
		return fiber.ErrBadRequest
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return fiber.ErrConflict
	}
	slog.Error(err.Error())
	return fiber.ErrInternalServerError
}
