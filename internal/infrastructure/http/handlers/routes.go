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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.RegisterRequest true "Registration details"
// @Success 201 {object} dtos.Token "User created successfully"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 409 {object} dtos.ErrMsg "User already exists"
// @Router /api/v1/register [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} dtos.Token "Login successful"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Invalid credentials"
// @Router /api/v1/login [post]
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

// GetSelf godoc
// @Summary Get current user info
// @Description Retrieve information about the currently authenticated user.
// @Tags User
// @Produce json
// @Security jwt
// @Success 200 {object} dtos.User "User information"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/user [get]
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

// UpdateUser godoc
// @Summary Update user info
// @Description Update information of the currently authenticated user.
// @Tags User
// @Accept json
// @Produce json
// @Security jwt
// @Param request body dtos.UpdateUserRequest true "Update details"
// @Success 200 {object} dtos.User "Updated user information"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/user [patch]
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

// DeleteUser godoc
// @Summary Delete user account
// @Description Delete the currently authenticated user account.
// @Tags User
// @Security jwt
// @Success 204 "Account deleted"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/user [delete]
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

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new task group for the user.
// @Tags Groups
// @Accept json
// @Produce json
// @Security jwt
// @Param request body dtos.CreateGroupRequest true "Group details"
// @Success 201 {object} dtos.Group "Group created"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/group [post]
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

// UpdateGroup godoc
// @Summary Update a group
// @Description Update an existing task group.
// @Tags Groups
// @Accept json
// @Produce json
// @Security jwt
// @Param request body dtos.UpdateGroupRequest true "Update details"
// @Success 200 {object} dtos.Group "Group updated"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Failure 404 {object} dtos.ErrMsg "Group not found"
// @Router /api/v1/group [patch]
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

// DeleteGroup godoc
// @Summary Delete a group
// @Description Delete an existing task group by ID.
// @Tags Groups
// @Security jwt
// @Param id path string true "Group ID" format(uuid)
// @Success 204 "Group deleted"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Failure 404 {object} dtos.ErrMsg "Group not found"
// @Router /api/v1/group/{id} [delete]
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

// GetListGroupsByUserID godoc
// @Summary Get all user groups
// @Description Retrieve a list of all task groups for the current user.
// @Tags Groups
// @Produce json
// @Security jwt
// @Success 200 {array} dtos.Group "List of groups"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Router /api/v1/groups [get]
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

// GetGroupByID godoc
// @Summary Get group by ID
// @Description Retrieve a specific task group by ID.
// @Tags Groups
// @Produce json
// @Security jwt
// @Param id path string true "Group ID" format(uuid)
// @Success 200 {object} dtos.Group "Group details"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Failure 404 {object} dtos.ErrMsg "Group not found"
// @Router /api/v1/group/{id} [get]
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

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task within a group.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security jwt
// @Param request body dtos.CreateTaskRequest true "Task details"
// @Success 201 {object} dtos.Task "Task created"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Failure 404 {object} dtos.ErrMsg "Group not found"
// @Router /api/v1/task [post]
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

// UpdateTask godoc
// @Summary Update a task
// @Description Update an existing task.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security jwt
// @Param request body dtos.UpdateTaskRequest true "Update details"
// @Success 200 {object} dtos.Task "Task updated"
// @Failure 400 {object} dtos.ErrMsg "Invalid request"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Failure 404 {object} dtos.ErrMsg "Task not found"
// @Router /api/v1/task [patch]
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

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete an existing task by ID.
// @Tags Tasks
// @Security jwt
// @Param id path string true "Task ID" format(uuid)
// @Success 204 "Task deleted"
// @Failure 400 {object} dtos.ErrMsg "Invalid ID"
// @Failure 401 {object} dtos.ErrMsg "Unauthorized"
// @Failure 404 {object} dtos.ErrMsg "Task not found"
// @Router /api/v1/task/{id} [delete]
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