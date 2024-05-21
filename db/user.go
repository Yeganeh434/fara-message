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

func (d *Database) AddContact(id int, userID int, contactID int) error {
	var user User
	result := d.db.First(&user, "ID=?", contactID)
	if result.Error != nil {
		return result.Error
	}
	contact := Contacts{ID: id, UserID: userID, ContactID: contactID}
	result = d.db.Create(&contact)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) IsContactExist(userID int, contactID int)(bool ,error){
	var contact Contacts
	result:=d.db.Table("contacts").Select("contacts.*").Where("user_id=? AND contact_id=?",userID,contactID).Find(&contact)
	if result.Error != nil {
		return true,result.Error
	}
	if contact.ID==0 {
		return false,nil
	}
	return true,nil
}

func (d *Database) DeleteContact(userID int, contactID int) error {
	result := d.db.Where("user_id=? AND contact_id=?", userID, contactID).Delete(&Contacts{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) GetContact(userID int) ([]User, error) {
	var contacts []User
	// result := d.db.Table("contacts").Select("User.*").Joins("JOIN users ON contacts.contact_id=users.id").Where("contacts.user_id=?", userID).Find(&contacts)
	result := d.db.Table("contacts").Select("users.*").Joins("JOIN users ON contacts.contact_id = users.ID").Where("contacts.user_id = ?", userID).Find(&contacts)
	if result.Error != nil {
		return nil, result.Error
	}
	return contacts, nil
}

func (d *Database) ChangePassword(ID string, newPassword string) error {
	result := d.db.Model(&User{}).Where("ID=?", ID).Update("Password", newPassword)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) IsUsernameAvailable(username string) (bool,error){
	var dbUsername string
	result:=d.db.Table("users").Select("users.username").Where("username=?",username).Find(&dbUsername)
	if result.Error!=nil {
		return false,result.Error
	}
	if len(dbUsername)==0 {
		return true,nil
	}
	return false,nil
}