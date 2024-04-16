package user // "github.com/amieldelatorre/notifi/internal/db/user"

import (
	"context"

	userModel "github.com/amieldelatorre/notifi/internal/model/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Postgres {
	return Postgres{Db: db}
}

func (pg *Postgres) GetUserById(userId int) (userModel.User, error) {
	var user userModel.User

	err := pg.Db.QueryRow(context.Background(), "SELECT * FROM users").Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.DatetimeCreated, &user.DatetimeUpdated)
	if err != nil {
		return user, err
	}

	return user, nil
}
