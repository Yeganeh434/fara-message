package db

import (
	"time"
)

type Message struct {
	ID       int `gorm:"primary_key"`
	SenderID int `gorm:"foreign_key"`
	ChatID   int `gorm:"foreign_key"`
	Content  string
}

type Gender int8

const (
	Male Gender = iota
	Female
)

type User struct {
	ID          string `gorm:"primary_key"`
	Username    string
	FirstName   string
	LastName    string
	Password    string
	Gender      Gender
	DateOfBirth time.Time
	CreatedTime time.Time
}

type UserInfo struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Gender      Gender    `json:"gender"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	CreatedTime time.Time `json:"createdTime"`
}

type Contacts struct {
	UserID           string `gorm:"foreignkey:ID"`
	ContantID        string `gorm:"foreignkey:ID"`
}

func ConvertUserToUserInfo(user User) UserInfo {
	return UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
		CreatedTime: user.CreatedTime,
	}
}
