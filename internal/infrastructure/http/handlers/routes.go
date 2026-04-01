package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/chishkin-afk/todo/internal/application/dtos"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Routes struct {
	as authService
	ts taskService
}

func (r *Routes) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.RegisterRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.as.Register(ctx, &req)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

func (r *Routes) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.LoginRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.as.Login(ctx, &req)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) GetSelf() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := r.as.GetSelf(ctx.Request.Context())
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.UpdateUserRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.as.Update(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := r.as.Delete(ctx.Request.Context()); err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (r *Routes) CreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.CreateGroupRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.ts.CreateGroup(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

func (r *Routes) UpdateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.UpdateGroupRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.ts.UpdateGroup(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) DeleteGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		if err := r.ts.DeleteGroup(ctx.Request.Context(), id); err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (r *Routes) GetListGroupsByUserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := r.ts.GetListGroupsByUserID(ctx.Request.Context())
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) GetGroupByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		resp, err := r.ts.GetGroupByID(ctx.Request.Context(), id)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) CreateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.CreateTaskRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.ts.CreateTask(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) UpdateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.UpdateTaskRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid request",
			})
			return
		}

		resp, err := r.ts.UpdateTask(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (r *Routes) DeleteTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrMsg{
				Error: "invalid id",
			})
			return
		}

		if err := r.ts.DeleteTask(ctx.Request.Context(), id); err != nil {
			ctx.JSON(r.getCode(err), dtos.ErrMsg{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (r *Routes) getCode(err error) int {
	switch {
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		return http.StatusRequestTimeout
	case errors.Is(err, errs.ErrInvalidPassword),
		errors.Is(err, errs.ErrInvalidEmail),
		errors.Is(err, errs.ErrInvalidUsername),
		errors.Is(err, errs.ErrInvalidTitle),
		errors.Is(err, errs.ErrInvalidTaskPriority),
		errors.Is(err, errs.ErrInvalidTaskDesc):
		return http.StatusBadRequest
	case errors.Is(err, errs.ErrUserAlreadyExists),
		errors.Is(err, errs.ErrTaskAlreadyDone),
		errors.Is(err, errs.ErrTaskNotDone):
		return http.StatusConflict
	case errors.Is(err, errs.ErrUserNotFound),
		errors.Is(err, errs.ErrGroupNotFound),
		errors.Is(err, errs.ErrTaskNotFound):
		return http.StatusNotFound
	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, errs.ErrNotEnoughRights):
		return http.StatusForbidden
	}

	return http.StatusInternalServerError
}
