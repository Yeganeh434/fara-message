package db

import "errors"

func (d *Database) DeleteMessage(message Message) error {
	result:=d.db.Delete(&message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) SendMessage(message Message) error {
	result := d.db.Create(&message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) FindMessageInfo(messageID int) (Message, error) {
	var message Message
	result := d.db.Table("messages").Select("messages.*").Where("id=?", messageID).Find(&message)
	if result.Error != nil {
		return message, result.Error
	}
	if result.RowsAffected == 0 {
		return message, errors.New("no record found with this id")
	}
	return message, nil
}

func (d *Database) IsChatContact(userID int, chatID int) (bool, error) {
	var dbChatID int
	result := d.db.Table("chat_members").Select("chat_id").Where("chat_id=? AND user_id=?", chatID, userID).Find(&dbChatID)
	if result.Error != nil {
		return false, result.Error
	}
	if dbChatID == 0 {
		return false, nil
	}
	return true, nil
}
