package main

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Login     string    `gorm:"unique;size:80" json:"login"`
	Password  string    `gorm:"size:255" json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
