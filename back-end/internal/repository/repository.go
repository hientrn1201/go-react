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

	// get the existing movie by id just for display
	OneMovie(id int) (*models.Movie, error)

	//get the existing movie by id to edit (required authorization)
	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)

	//query user by email
	GetUserByEmail(email string) (*models.User, error)

	//query user by id
	GetUserByID(id int) (*models.User, error)

	// get all genres
	AllGenres() ([]*models.Genre, error)

	// insert one movie
	InsertMovie(movie models.Movie) (int, error)

	//update movie
	UpdateMovie(movie models.Movie) error

	//update movie genres id list
	UpdateMovieGenres(id int, genreIDs []int) error

	//delete one movie
	DeleteMovie(id int) error
}
