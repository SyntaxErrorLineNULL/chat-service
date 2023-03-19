package dto

type UserDTO struct {
	ID        string `bson:"id"`
	FirstName string `bson:"firstname"`
	LastName  string `bson:"lastname"`
	UserName  string `bson:"username"`
	Email     string `bson:"email"`
}
