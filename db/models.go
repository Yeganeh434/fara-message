package db

import (
	"database/sql"
	"time"
)

// 0 equals direct chat and 1 equals group chat
type Chat struct {
	ID          uint64 `gorm:"primary_key"`
	HashID      uint64
	Name        string
	CreatedTime time.Time
	Type        int
}

type ChatMember struct {
	ChatID     uint64 `gorm:"foreignkey:ID"`
	UserID     uint64 `gorm:"foreignkey:ID"`
	JoinedTime time.Time
}

type Message struct {
	ID       uint64 `gorm:"primary_key"`
	SenderID uint64 `gorm:"foreign_key"`
	ChatID   uint64 `gorm:"foreign_key"`
	Content  string
	Time     time.Time
}

// 0 equals male and 1 equals female
type User struct {
	ID          uint64 `gorm:"primary_key"`
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

type Contacts struct {
	ID        uint64 `gorm:"primary_key"`
	UserID    uint64 `gorm:"foreignkey:ID"`
	ContactID uint64 `gorm:"foreignkey:ID"`
}

type OTP struct {
	ID             uint64 `gorm:"primary_key"`
	OTP            int
	Email          string `gorm:"foreignkey:Email"`
	ExpirationTime time.Time
}
