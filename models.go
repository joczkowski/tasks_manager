package main

import "time"

type User struct {
	Id             int `gorm:"primaryKey"`
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Task struct {
	Id          int `gorm:"primaryKey"`
	Title       string
	Description string
	UserId      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string
}
