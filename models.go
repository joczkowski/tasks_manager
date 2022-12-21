package main

type User struct {
	Id             int `gorm:"primaryKey"`
	Email          string
	HashedPassword string
}
