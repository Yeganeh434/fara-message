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

// if empty!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func (d *Database) GetChatsList(userID string) ([]string, error) {
	//get groups name
	var chatsName []string
	result := d.db.Table("chats").Select("chats.name").Joins("JOIN chat_members ON chats.id=chat_members.chat_id").Where("chats.type=? AND chat_members.user_id=?", 1, userID).Find(&chatsName)
	if result.Error != nil {
		return nil, result.Error
	}

	//get direct chats name
	var usersID []string
	subQuery := d.db.Table("chats").Select("chats.id").Joins("JOIN chat_members ON chat_members.chat_id=chats.id").Where("chats.type=? AND chat_members.user_id=?", 0, userID)
	if subQuery.Error != nil {
		return nil, subQuery.Error
	}
	result = d.db.Table("chat_members").Select("user_id").Where("chat_id IN (?)", subQuery).Not("user_id=?", userID).Find(&usersID)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, value := range usersID {
		user, err := d.ReadUser(value)
		if err != nil {
			return nil,err
		}
		name := user.FirstName + " " + user.LastName
		chatsName = append(chatsName, name)
	}
	return chatsName, nil
}
