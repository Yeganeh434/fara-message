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
		return message, errors.New("no record found with this message ID")
	}
	return message, nil
}


