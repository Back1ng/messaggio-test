package entity

import "time"

type Message struct {
	ID          int64      `json:"id"`
	Message     string     `json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at"`
}
