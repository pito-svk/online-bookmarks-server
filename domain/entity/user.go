package entity

type Name struct {
	FirstName string
	LastName  string
}

type User struct {
	ID          string
	Name        Name
}
