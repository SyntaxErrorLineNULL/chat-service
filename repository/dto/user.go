package dto

import "github.com/SyntaxErrorLineNULL/chat-service/domain"

type repositoryUser struct {
	ID        string `bson:"id"`
	FirstName string `bson:"firstname"`
	LastName  string `bson:"lastname"`
	UserName  string `bson:"username"`
	Email     string `bson:"email"`
}

func (ru *repositoryUser) toDomain() *domain.User {
	return &domain.User{
		ID:        ru.ID,
		FirstName: ru.FirstName,
		LastName:  ru.LastName,
		UserName:  ru.UserName,
		Email:     ru.Email,
	}
}

func (ru *repositoryUser) fromDomain(user *domain.User) {
	ru.ID = user.ID
	ru.FirstName = user.FirstName
	ru.LastName = user.LastName
	ru.UserName = user.UserName
	ru.Email = user.Email
}
