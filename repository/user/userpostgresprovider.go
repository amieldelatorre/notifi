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

func (provider *UserPostgresProvider) GetUserById(ctx context.Context, id int) (model.User, error) {
	var user model.User

	// User Id's should be unique when querying the user table
	err := provider.DbPool.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1", id).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.DatetimeCreated, &user.DatetimeUpdated)
	if err != nil {
		return user, err
	}

	return user, nil
}
