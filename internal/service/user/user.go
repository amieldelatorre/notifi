package user // import "github.com/amieldelatorre/notifi/internal/service/user"

import (
	userDb "github.com/amieldelatorre/notifi/internal/db/user"
	userModel "github.com/amieldelatorre/notifi/internal/model/user"
)

type UserDbProvider interface {
	GetUserById(userId int) (userModel.User, error)
}

type Service struct {
	Db *userDb.Postgres
}

func New(db *userDb.Postgres) Service {
	return Service{Db: db}
}

func (service *Service) GetUserById(userId int) (userModel.User, error) {
	user, err := service.Db.GetUserById(userId)
	return user, err
}
