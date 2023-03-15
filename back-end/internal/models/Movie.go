package models

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	RunTime     int       `json:"runtime"`
	MPAARating  string    `json:"mpaa_rating"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreateAt    time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Genres      []*Genre  `json:"genres,omitempty"`
	GenresArray []int     `json:"genres_array,omitempty"`
}

//"-" means that dont include in JSON

type Genre struct {
	ID      int    `json:"id"`
	Genre   string `json:"genre"`
	Checked bool   `json:"checked"`
}
