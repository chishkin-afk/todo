package dtos

type Token struct {
	Token string `json:"token"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Group struct {
	ID        string `json:"id"`
	OwnerID   string `json:"owner_id"`
	Title     string `json:"title"`
	Tasks     []Task `json:"tasks"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Task struct {
	ID         string `json:"id"`
	OwnerID    string `json:"owner_id"`
	GroupID    string `json:"group_id"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Priority   string `json:"priority"`
	PriorityID int64  `json:"priority_id"`
}

type Groups struct {
	Groups []Group `json:"groups"`
}
