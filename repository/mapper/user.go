package mapper

import (
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository/dto"
)

type UserMapper struct{}

func (m *UserMapper) ToDTO(usr *domain.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:        usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		UserName:  usr.UserName,
		Email:     usr.Email,
	}
}

func (m *UserMapper) ToModel(dto *dto.UserDTO) *domain.User {
	return &domain.User{
		ID:        dto.ID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		UserName:  dto.UserName,
		Email:     dto.Email,
	}
}
