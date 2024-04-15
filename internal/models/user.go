package models

type UserRole string

const (
	Client UserRole = "client"
	Admin  UserRole = "admin"
)

type User struct {
	Username string   `db:"username" json:"username" validate:"required"`
	Password string   `db:"password" json:"password" validate:"required"`
	Role     UserRole `db:"role" json:"role"`
}

type Credentials struct {
	HashPassword string   `db:"hash_password" json:"hash_password" validate:"required"`
	Role         UserRole `db:"role" json:"role"`
}

type UserResources struct {
	Role      UserRole `db:"role" json:"role"`
	Resources []string `db:"resources" json:"resources"`
}
