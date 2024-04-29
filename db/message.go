package db

func (d *Database) DeleteMessage(messageID int) error {
	var message Message
	message.ID = messageID
	result := d.db.Where("ID=?", message.ID).Delete(&Message{})
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
