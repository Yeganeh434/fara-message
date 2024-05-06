package db

func (d *Database) CreateUser(user User) error {
	result := d.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) ReadUserByUsername(username string) (User, error) {
	var user User
	result := d.db.First(&user, "username=?", username)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (d *Database) ReadUser(ID string) (UserInfo, error) {
	var user UserInfo
	result := d.db.First(&user, "ID=?", ID)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (d *Database) UpdateUser(ID string, newInfo UserInfo) error {
	result := d.db.Model(&User{}).Where("ID=?", ID).Updates(User{FirstName: newInfo.FirstName, LastName: newInfo.LastName, Gender: newInfo.Gender, DateOfBirth: newInfo.DateOfBirth})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) DeleteUser(ID string) error {
	var user User
	result := d.db.First(&user, "ID=?", ID)
	if result.Error != nil {
		return result.Error
	}
	result = d.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d* Database) AddContact(userID int,contactID int) error{
	contact:=Contacts{UserID:userID , ContactID:contactID}
	result:=d.db.Create(&contact)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
