package models

type User struct {
	TgID      int64
	ChatID    int64
	Role      string
	FirstName string
	LastName  string
	UserName  string
}
