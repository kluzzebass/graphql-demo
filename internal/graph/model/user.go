package model

import "time"

type User struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phoneNumber"`
	Address     *Address     `json:"address"`
	Role        Role         `json:"role"`
	CreatedAt   time.Time    `json:"createdAt"`
	LastLogin   time.Time    `json:"lastLogin"`
	Preferences *Preferences `json:"preferences"`
	PostIDs     []int        `json:"postIDs"`
}

func (u *User) GetID() int {
	return u.ID
}
