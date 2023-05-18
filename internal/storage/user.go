package storage

import (
	"context"
	"database/sql"
	"github.com/altiby/son/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	postgres *Postgres
}

type User struct {
	ID         string `db:"id"`
	Role       string `db:"role"`
	FirstName  string `db:"first_name"`
	SecondName string `db:"second_name"`
	Birthdate  string `db:"birthdate"`
	Biography  string `db:"biography"`
	City       string `db:"city"`
	Password   string `db:"password"`
}

func (u *User) FromDomainUser(user domain.User, password string) {
	u.ID = user.ID
	u.Role = user.Role
	u.FirstName = user.FirstName
	u.SecondName = user.SecondName
	u.Birthdate = user.Birthdate
	u.Biography = user.Biography
	u.City = user.City
	u.Password = password
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		ID:         u.ID,
		Role:       u.Role,
		FirstName:  u.FirstName,
		SecondName: u.SecondName,
		Birthdate:  u.Birthdate,
		Biography:  u.Biography,
		City:       u.City,
	}
}

func (u UserStorage) RegisterUser(ctx context.Context, user domain.User, password string) error {
	dbUser := &User{}
	dbUser.FromDomainUser(user, password)
	_, err := sqlx.NamedExecContext(
		ctx,
		u.postgres.db,
		`INSERT INTO 
    				users(id, role, first_name, second_name, birthdate, biography, city, password) 
				VALUES 
				    (:id, :role, :first_name, :second_name, :birthdate, :biography, :city, :password)`,
		dbUser)
	return err
}

func (u UserStorage) AuthorizeUser(ctx context.Context, id string, password string) (domain.User, error) {
	dbUser := &User{}
	err := sqlx.GetContext(
		ctx,
		u.postgres.db,
		dbUser,
		"SELECT id, role, first_name, second_name, birthdate, biography, city FROM users WHERE id = $1 AND password = $2", id, password)

	if err == sql.ErrNoRows {
		return domain.User{}, domain.ErrUserNotFound
	}

	return dbUser.ToDomain(), err
}

func (u UserStorage) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	dbUser := &User{}
	err := sqlx.GetContext(
		ctx,
		u.postgres.db,
		dbUser,
		"SELECT id, role, first_name, second_name, birthdate, biography, city FROM users WHERE id = $1", id)

	if err == sql.ErrNoRows {
		return domain.User{}, domain.ErrUserNotFound
	}

	return dbUser.ToDomain(), err
}

func NewUserStorage(p *Postgres) *UserStorage {
	return &UserStorage{postgres: p}
}
