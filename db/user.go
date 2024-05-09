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

func (d *Database) ReadUser(ID string) (User, error) {
	var user User
	result := d.db.First(&user, "ID=?", ID)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (d *Database) UpdateUser(ID string, newInfo User) error {
	result := d.db.Model(&User{}).Where("ID=?", ID).Updates(User{Username: newInfo.Username, FirstName: newInfo.FirstName, LastName: newInfo.LastName, Gender: newInfo.Gender, DateOfBirth: newInfo.DateOfBirth})
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

func (d *Database) AddContact(userID int, contactID int) error {
	contact := Contacts{UserID: userID, ContactID: contactID}
	result := d.db.Create(&contact)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) DeleteContact(userID int, contactID int) error {
	result := d.db.Where("UserID=? AND ContactID", userID, contactID).Delete(&Contacts{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) GetContact(userID int) ([]User, error) {
	var contacts []User
	result := d.db.Table("Contacts").Select("User.*").Joins("JOIN User ON Contacts.ContactID=User.ID").Where("Contacts.UserID=?", userID).Find(&contacts)
	if result.Error != nil {
		return nil, result.Error
	}
	return contacts, nil
}

func (d *Database) ChangePassword(ID string,newPassword string) error {
	result:=d.db.Model(&User{}).Where("ID=?",ID).Update("Password",newPassword)
	if result.Error !=nil {
		return result.Error
	}
	return nil
} 
