package user // import "github.com/amieldelatorre/notifi/repository/user"

import (
	"context"

	"github.com/amieldelatorre/notifi/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgresProvider struct {
	DbPool *pgxpool.Pool
}

func NewUserPostgresProvider(dbPool *pgxpool.Pool) *UserPostgresProvider {
	return &UserPostgresProvider{DbPool: dbPool}
}

func (provider *UserPostgresProvider) CreateUser(ctx context.Context, input model.UserInput) (model.User, error) {
	var newUser model.User
	tx, err := provider.DbPool.Begin(ctx)
	if err != nil {
		return newUser, err
	}
	defer tx.Rollback(ctx)

	err = provider.DbPool.QueryRow(ctx,
		`INSERT INTO Users (email, firstName, lastName, password, datetimeCreated, datetimeUpdated) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id, email, firstName, lastName, password, datetimeCreated, datetimeUpdated`,
		input.Email, input.FirstName, input.LastName, input.Password).Scan(
		&newUser.Id, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.Password,
		&newUser.DatetimeCreated, &newUser.DatetimeUpdated)
	if err != nil {
		return newUser, nil
	}

	err = tx.Commit(ctx)
	return newUser, err
}

func (provider *UserPostgresProvider) GetUserById(ctx context.Context, id int) (model.User, error) {
	var user model.User

	// User Id's should be unique when querying the user table
	err := provider.DbPool.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.DatetimeCreated, &user.DatetimeUpdated)
	return user, err
}

func (provider *UserPostgresProvider) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	// User email's should be unique when querying the user table
	err := provider.DbPool.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.DatetimeCreated, &user.DatetimeUpdated)
	return user, err
}
