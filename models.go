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

type Page struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"unique" json:"slug"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Page) TableName() string {
	return "page"
}
