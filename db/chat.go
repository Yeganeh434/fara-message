package db

import (
	"errors"
	"fmt"
	"time"
)

func (d *Database) NewChat(chatID uint64, chatName string, chatType int, users []uint64) error {
	var hashID uint64
	if chatType == 0 {
		hashID = users[0] + users[1]
	}
	chat := Chat{
		ID:          chatID,
		HashID:      hashID,
		Name:        chatName,
		Type:        chatType,
		CreatedTime: time.Now(),
	}
	var chatMember ChatMember
	for _, value := range users {
		chatMember = ChatMember{
			JoinedTime: time.Now(),
			ChatID:     chat.ID,
			UserID:     value,
		}
		if err := d.db.Create(&chatMember).Error; err != nil {
			return err
		}
	}
	//handle error!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	d.db.Create(&chat)
	return nil
}

func (d *Database) AddMemberToGroup(chatID uint64, newMemberID uint64) error {
	chatMember := ChatMember{
		ChatID:     chatID,
		UserID:     newMemberID,
		JoinedTime: time.Now(),
	}
	result := d.db.Create(&chatMember)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) IsChatExist(users []uint64) (bool, error) {
	hashID := users[0] + users[1]
	var dbHashID uint64
	result := d.db.Table("chats").Select("hash_id").Where("hash_id=?", hashID).Find(&dbHashID)
	if result.Error != nil {
		return true, result.Error
	}
	if dbHashID == hashID {
		return true, nil
	}
	return false, nil
}

func (d *Database) GetChatMessages(userID uint64, chatID uint64) ([]Message, error) {
	var messages []Message
	var joinedTime time.Time
	result := d.db.Table("chat_members").Select("joined_time").Where("user_id=? AND chat_id=?", userID, chatID).Find(&joinedTime)
	if result.Error != nil {
		return messages, result.Error
	}
	if err := d.db.Where("chat_id = ? AND Time > ?", chatID, joinedTime).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("no  message found for chat %w", err)
	}
	return messages, nil
}

func (d *Database) GetChatMembers(chatID uint64) ([]User, error) {
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
			return nil, err
		}
		name := user.FirstName + " " + user.LastName
		chatsName = append(chatsName, name)
	}
	return chatsName, nil
}

func (d *Database) IsAChatContact(userID uint64, chatID uint64) (bool, error) {
	var dbChatID int
	result := d.db.Table("chat_members").Select("chat_id").Where("chat_id=?", chatID).Find(&dbChatID)
	if result.Error != nil {
		return false, result.Error
	}
	if dbChatID == 0 {
		return false, errors.New("no record found with this chat ID")
	}
	var dbUserID int
	result = d.db.Table("chat_members").Select("user_id").Where("chat_id=? AND user_id=?", chatID, userID).Find(&dbUserID)
	if result.Error != nil {
		return false, result.Error
	}
	if dbUserID == 0 {
		return false, nil
	}
	return true, nil
}
