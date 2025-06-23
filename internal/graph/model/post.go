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

func (p *Post) GetID() int {
	return p.ID
}
