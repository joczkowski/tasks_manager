package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id             int `gorm:"primaryKey"`
	Email          string
	HashedPassword string
}
