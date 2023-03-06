package repository

import (
	"backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	//return the pointer to the SQL db
	Connection() *sql.DB

	//return a list of pointers that point to every movie queried from the database
	AllMovies() ([]*models.Movie, error)
}
