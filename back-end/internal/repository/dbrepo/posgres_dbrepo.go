package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

// declare a type PostgresDBRepo that inherits the interface DatabaseRepo
type PostgresDBRepo struct {
	DB *sql.DB
}

// if users interact with the DB more than 3 seconds, time out
const dbTimeOut = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {
	//you have a limited time with the context before time out
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	//coalesce(image, '') = return image if exists otherwise return ''
	//we need to do this since Go really does not like nil value
	query := `
		select
			id, title, release_date, runtime,
			mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		from
			movies
		order by
			title
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	//you have to close rows
	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreateAt,
			&movie.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	//you have a limited time with the context before time out
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
						created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.Lastname,
		&user.Password,
		&user.CreatedAt,
		&user.UpdateAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
