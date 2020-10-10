package entity

type UserInterface interface {
	SetID()
	SetHashedPassword()
	ClearPassword()
}

type User struct {
	ID        string
	Email     string
	Password  string
	FirstName string
	LastName  string
}
