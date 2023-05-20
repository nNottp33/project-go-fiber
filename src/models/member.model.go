package models

import "time"

type Members struct {
	Id           int8      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"index:constraint_member,unique;not null"`
	Password     string    `json:"password" gorm:"not null"`
	FullName     string    `json:"full_name" gorm:"index"`
	Email        string    `json:"email" gorm:"index"`
	PhoneNumber  string    `json:"phone_number" `
	Role         string    `json:"role"`
	SessionToken string    `json:"session_token" gorm:"index"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	IsActive     bool      `json:"is_active" gorm:"default:true;index"`
}
