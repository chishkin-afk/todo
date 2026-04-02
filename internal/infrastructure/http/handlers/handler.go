package handlers

import (
	"context"
	"errors"

	"github.com/chishkin-afk/todo/internal/application/dtos"
	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type authService interface {
	Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.Token, error)
	Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.Token, error)
	GetSelf(ctx context.Context) (*dtos.User, error)
	Update(ctx context.Context, req *dtos.UpdateUserRequest) (*dtos.User, error)
	Delete(ctx context.Context) error
}

type taskService interface {
	CreateGroup(ctx context.Context, req *dtos.CreateGroupRequest) (*dtos.Group, error)
	UpdateGroup(ctx context.Context, req *dtos.UpdateGroupRequest) (*dtos.Group, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) error
	GetListGroupsByUserID(ctx context.Context) (*dtos.Groups, error)
	GetGroupByID(ctx context.Context, id uuid.UUID) (*dtos.Group, error)
	CreateTask(ctx context.Context, req *dtos.CreateTaskRequest) (*dtos.Task, error)
	UpdateTask(ctx context.Context, req *dtos.UpdateTaskRequest) (*dtos.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

func New(cfg *config.Config, as authService, ts taskService, mls []gin.HandlerFunc) (*gin.Engine, error) {
	var router *gin.Engine
	switch cfg.App.Env {
	case "prod":
		router = gin.New()
		router.Use(gin.Recovery())
	case "dev", "local":
		router = gin.Default()
	default:
		return nil, errors.New("invalid environment")
	}

	router.Use(mls...)

	routes := Routes{
		as: as,
		ts: ts,
	}

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/register", routes.Register())
			v1.POST("/login", routes.Login())

			v1.GET("/user", routes.GetSelf())
			v1.PATCH("/user", routes.UpdateUser())
			v1.DELETE("/user", routes.DeleteUser())

			v1.POST("/group", routes.CreateGroup())
			v1.PATCH("/group", routes.UpdateGroup())
			v1.DELETE("/group/:id", routes.DeleteGroup())
			v1.GET("/groups", routes.GetListGroupsByUserID())
			v1.GET("/group/:id", routes.GetGroupByID())

			v1.POST("/task", routes.CreateTask())
			v1.PATCH("/task", routes.UpdateTask())
			v1.DELETE("/task/:id", routes.DeleteTask())

		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router, nil
}
