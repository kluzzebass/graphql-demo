package model

import "time"

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Title     string    `json:"title"`
	Ingress   string    `json:"ingress"`
	Content   string    `json:"content"`
	Category  Category  `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phoneNumber"`
	Address     *Address     `json:"address"`
	Role        Role         `json:"role"`
	LastLogin   time.Time    `json:"lastLogin"`
	Preferences *Preferences `json:"preferences"`
	PostIDs     []int        `json:"postIDs"`
}
