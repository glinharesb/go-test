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

func (d *Database) LoadAccountByName(accountName string) (Account, error) {
	var account Account

	query := "SELECT id, name, password, premdays FROM `accounts` WHERE `name` = ? LIMIT 1"
	res, err := d.Connection.Query(query, accountName)
	if err != nil {
		return account, err
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&account.Id, &account.Name, &account.Password, &account.Premdays)
		if err != nil {
			return account, err
		}
	}

	return account, nil
}

func (d *Database) LoadCharactersById(accountId int) ([]string, error) {
	var characters []string

	res, err := d.Connection.Query("SELECT name FROM `players` WHERE `account_id` = ?", accountId)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var character Character
		err = res.Scan(&character.Name)
		if err != nil {
			return nil, err
		}

		characters = append(characters, character.Name)
	}

	return characters, nil
}
