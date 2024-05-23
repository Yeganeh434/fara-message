package db

import (
	"database/sql"
	"time"
)

// 0 equals direct chat and 1 equals group chat
type Chat struct {
	ID          int `gorm:"primary_key"`
	HashID      string
	Name        string
	CreatedTime time.Time
	Type        int
}

type ChatMember struct {
	ChatID     int `gorm:"foreignkey:ID"`
	UserID     int `gorm:"foreignkey:ID"`
	JoinedTime time.Time
}

type Message struct {
	ID       int `gorm:"primary_key"`
	SenderID int `gorm:"foreign_key"`
	ChatID   int `gorm:"foreign_key"`
	Content  string
	Time     time.Time
}

// 0 equals male and 1 equals female
type User struct {
	ID          string `gorm:"primary_key"`
	Username    string
	FirstName   string
	LastName    string
	Password    string
	Gender      int
	DateOfBirth time.Time
	Email       string
	CreatedTime time.Time
	DeletedAt   sql.NullTime
}

// type UserInfo struct {
// 	ID          string    `json:"id"`
// 	Username    string    `json:"username"`
// 	FirstName   string    `json:"firstname"`
// 	LastName    string    `json:"lastname"`
// 	Gender      int       `json:"gender"`
// 	DateOfBirth time.Time `json:"dateOfBirth"`
// 	CreatedTime time.Time `json:"createdTime"`
// }

type Contacts struct {
	ID        int `gorm:"primary_key"`
	UserID    int `gorm:"foreignkey:ID"`
	ContactID int `gorm:"foreignkey:ID"`
}

// func ConvertUserToUserInfo(user User) UserInfo {
// 	return UserInfo{
// 		ID:          user.ID,
// 		Username:    user.Username,
// 		FirstName:   user.FirstName,
// 		LastName:    user.LastName,
// 		Gender:      user.Gender,
// 		DateOfBirth: user.DateOfBirth,
// 		CreatedTime: user.CreatedTime,
// 	}
// }

type OTP struct {
	ID    int `gorm:"primary_key"`
	OTP   int
	Email string `gorm:"foreignkey:Email"`
}
