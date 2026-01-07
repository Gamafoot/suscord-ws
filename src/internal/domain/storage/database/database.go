package database

type Storage interface {
	User() UserStorage
	Chat() ChatStorage
	Session() SessionStorage
}
