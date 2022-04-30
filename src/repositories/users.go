package repositories

import "main/src/entities"

type UsersRepository interface {
	GetUserByExternalId(id int64, chatId int64) (*entities.User, error)
	Save(user *entities.User) error
	UpdateBalance(user *entities.User, amount uint) error
}
