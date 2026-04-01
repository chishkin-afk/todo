package dtos

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Password *string `json:"password"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type CreateGroupRequest struct {
	Title string `json:"title"`
}

type UpdateGroupRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CreateTaskRequest struct {
	GroupID    string `json:"group_id"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	PriorityID int    `json:"priority_id"`
}

type UpdateTaskRequest struct {
	ID         string  `json:"id"`
	IsDone     *bool   `json:"is_done"`
	PriorityID *int    `json:"priority_id"`
	Title      *string `json:"title"`
	Desc       *string `json:"desc"`
}
