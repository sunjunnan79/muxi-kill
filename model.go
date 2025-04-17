package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name string
	Num  int64
}

func InitDB() *gorm.DB {
	dsn := "root:xxxxx@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Item{})
	if err != nil {
		return nil
	}
	return db
}
