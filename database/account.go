package database

func GetAccountByName(accountName string) (Account, error) {
	var account Account

	if err := GetDb().Where("name = ?", accountName).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func GetCharactersById(accountID int) ([]string, error) {
	var names []string

	result := GetDb().Model(&Player{}).Where("account_id = ?", accountID).Pluck("name", &names)
	if result.Error != nil {
		return nil, result.Error
	}

	return names, nil
}
