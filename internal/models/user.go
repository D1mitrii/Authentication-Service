package models

import "time"

type User struct {
	Id        int       `db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
