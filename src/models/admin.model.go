package models

import "time"

type AdminUsers struct {
	Id           int8      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"index:constraint_admin,unique;not null"`
	Password     string    `json:"password" gorm:"not null"`
	FullName     string    `json:"full_name" gorm:"index"`
	Role         string    `json:"role"`
	SessionToken string    `json:"session_token" gorm:"index"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	IsActive     bool      `json:"is_active" gorm:"default:true;index"`
}
