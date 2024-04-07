package models

type UserRole string

type User struct {
	Username string `db:"username" json:"username" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
	IsAdmin  bool   `db:"is_admin" json:"is_admin"`
}

type SignInInput struct {
	Username string `db:"username" json:"username" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
}
