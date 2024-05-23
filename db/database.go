package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Mysql Database

type Database struct {
	db *gorm.DB
}

func InitialDatabase() {
	var err error
	dsn := "root:Yeganeh-2004@tcp(localhost:3306)/faramessage?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Mysql.db=gormDB
	if err != nil {
		panic("failed to connect to database")
	}
	err = Mysql.db.AutoMigrate(&Chat{}, &ChatMember{}, &Message{},&User{},&Contacts{},&OTP{})
	if err != nil {
		log.Printf("failed to migrate: %v", err)
		return
	}
	fmt.Println("Migration done ..")
}
