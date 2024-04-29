package api

import (
	"time"

	"github.com/mhghw/fara-message/db"
)

type HTTPError struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	ID       int        `json:"chatId"`
	Name     string     `json:"chatName"`
	Messages []Message  `json:"messages"`
	Users    []UserInfo `json:"users"`
}

type AnotherUserInfo struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type UserInfo struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Gender      int       `json:"gender"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	CreatedTime time.Time `json:"createdTime"`
}

func ConvertUserInfo(newInfo UserInfo) db.UserInfo {
	gender := db.Male
	if newInfo.Gender != 0 {
		gender = db.Female
	}
	return db.UserInfo{
		ID:          newInfo.ID,
		Username:    newInfo.Username,
		FirstName:   newInfo.FirstName,
		LastName:    newInfo.LastName,
		Gender:      gender,
		DateOfBirth: newInfo.DateOfBirth,
		CreatedTime: newInfo.CreatedTime,
	}
}

func ConvertdbUserInfo(newInfo db.UserInfo) UserInfo {
	return UserInfo{
		ID:          newInfo.ID,
		Username:    newInfo.Username,
		FirstName:   newInfo.FirstName,
		LastName:    newInfo.LastName,
		Gender:      int(newInfo.Gender),
		DateOfBirth: newInfo.DateOfBirth,
		CreatedTime: newInfo.CreatedTime,
	}
}

