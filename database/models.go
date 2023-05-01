package database

type Account struct {
	Id       int
	Name     string
	Password string
	Premdays int
}

type Character struct {
	Name string
}

type Player struct {
	Name      string
	AccountId int
}
