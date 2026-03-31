package taskpg

import (
	"database/sql"
)

type taskPersistenceRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *taskPersistenceRepository {
	return &taskPersistenceRepository{
		db: db,
	}
}
