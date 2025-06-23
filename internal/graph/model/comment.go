package model

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postID"`
	UserID    int       `json:"userID"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
}

func (c *Comment) GetID() int {
	return c.ID
}
