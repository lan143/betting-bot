package memory

import (
	"main/src/entities"
	"main/src/repositories"
)

type InMemoryUsersRepository struct {
	users []*entities.User
}

func (r *InMemoryUsersRepository) GetUserByExternalId(id int64, chatId int64) (*entities.User, error) {
	for _, user := range r.users {
		if user.ExternalId == id && user.ExternalChatId == chatId {
			return user, nil
		}
	}

	return nil, nil
}

func (r *InMemoryUsersRepository) UpdateBalance(user *entities.User, amount uint) error {
	user.Balance += amount

	return nil
}

func (r *InMemoryUsersRepository) Save(user *entities.User) error {
	for index, existUser := range r.users {
		if existUser.Id == user.Id {
			r.users[index] = user
			return nil
		}
	}

	r.users = append(r.users, user)

	return nil
}

func NewInMemoryUsersRepository() repositories.UsersRepository {
	return &InMemoryUsersRepository{}
}
