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

func (d *Database) GetChatMessages(ChatID int64) ([]Message, error) {
	var messages []Message
	if err := d.db.Where("chat_id = ?", ChatID).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("no  message found for chat %w", err)
	}
	return messages, nil
}

func (d *Database) GetUsersChatMembers(userID int) ([]ChatMember, error) {
	var usersChats []ChatMember
	if err := d.db.Where("user_id = ?", userID).Find(&usersChats).Error; err != nil {
		return nil, fmt.Errorf("no  chat found for user %w", err)
	}
	return usersChats, nil
}
