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

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Password    string    `json:"password"`
	Gender      int       `json:"gender"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	CreatedTime time.Time `json:"createdTime"`
}

type ContactsResponse struct {
	Contacts []AnotherUserInfo `json:"contacts"`
}

func ConvertUserInfo(newInfo UserInfo) db.UserInfo {
	return db.UserInfo{
		ID:          newInfo.ID,
		Username:    newInfo.Username,
		FirstName:   newInfo.FirstName,
		LastName:    newInfo.LastName,
		Gender:      newInfo.Gender,
		DateOfBirth: newInfo.DateOfBirth,
		CreatedTime: newInfo.CreatedTime,
	}
}

func ConvertdbUser(newInfo db.User) UserInfo {
	return UserInfo{
		ID:          newInfo.ID,
		Username:    newInfo.Username,
		FirstName:   newInfo.FirstName,
		LastName:    newInfo.LastName,
		Gender:      newInfo.Gender,
		DateOfBirth: newInfo.DateOfBirth,
		CreatedTime: newInfo.CreatedTime,
	}
}

func ConvertUpdateUser(newInfo UpdateUser) db.User {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, newInfo.DateOfBirth)    //handle error!!!!!!!!!!!!!!!
	return db.User{
		Username:    newInfo.Username,
		FirstName:   newInfo.FirstName,
		LastName:    newInfo.LastName,
		Gender:      newInfo.Gender,
		DateOfBirth: date,
	}
}
