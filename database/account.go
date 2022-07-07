package database

import "log"

type Account struct {
	Id       int
	Name     string
	Password string
	Premdays int
}

type Character struct {
	Name string
}

func (d *Database) LoadAccountByName(accountName string) Account {
	res, err := d.Connection.Query("SELECT id, name, password, premdays FROM `accounts` WHERE `name` = ?", accountName)
	if err != nil {
		log.Fatal(err)
	}

	var account Account
	for res.Next() {
		err = res.Scan(&account.Id, &account.Name, &account.Password, &account.Premdays)
		if err != nil {
			log.Fatal(err)
		}
	}

	return account
}

func (d *Database) LoadCharactersById(accountId int) []string {
	var characters []string

	res, err := d.Connection.Query("SELECT name FROM `players` WHERE `account_id` = ?", accountId)
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {
		var character Character
		err = res.Scan(&character.Name)
		if err != nil {
			log.Fatal(err)
		}

		characters = append(characters, character.Name)
	}

	return characters
}
