package db

import (
	"fmt"
	"strconv"
	"time"
)

func (d *Database) NewChat(chatID string, chatName string, chatType int, users []string) error {
	id, _ := strconv.Atoi(chatID)
	chat := Chat{
		ID:          id,
		Name:        chatName,
		Type:        chatType,
		CreatedTime: time.Now(),
	}
	var chatMember ChatMember
	for _, value := range users {
		userID, _ := strconv.Atoi(value)
		chatMember = ChatMember{
			JoinedTime: time.Now(),
			ChatID:     chat.ID,
			UserID:     userID,
		}
		if err := d.db.Create(&chatMember).Error; err != nil {
			// d.db.Delete(&chatMembers)
			return err
		}
	}

	d.db.Create(&chat)
	return nil
}

func (d *Database) GetChatMessages(ChatID int) ([]Message, error) {
	var messages []Message
	if err := d.db.Where("chat_id = ?", ChatID).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("no  message found for chat %w", err)
	}
	return messages, nil
}

func (d *Database) GetChatMembers(chatID int) ([]User, error) {
	var members []User
	result := d.db.Table("chat_members").Select("users.*").Joins("JOIN users ON chat_members.user_id = users.ID").Where("chat_members.chat_id = ?", chatID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}
