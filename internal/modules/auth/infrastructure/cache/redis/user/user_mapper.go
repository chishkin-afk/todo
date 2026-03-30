package userredis

import (
	"encoding/json"

	"github.com/chishkin-afk/todo/internal/modules/auth/domain/user"
)

func ToBytes(domain *user.User) ([]byte, error) {
	bytes, err := json.Marshal(UserModel{
		ID:        domain.ID(),
		Email:     domain.Email().String(),
		Username:  domain.Username(),
		CreatedAt: domain.CreatedAt(),
		UpdatedAt: domain.UpdatedAt(),
	})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ToDomain(bytes []byte) (*user.User, error) {
	var model UserModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	return user.From(
		model.ID,
		user.Email(model.Email),
		"",
		model.Username,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
