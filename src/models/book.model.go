package models

import "time"

type Book struct {
	Id         int8      `json:"id" gorm:"primaryKey"`
	BookNameTh string    `json:"book_name_th" gorm:"index;not null"`
	BookNameEn string    `json:"book_name_en" gorm:"index;not null"`
	Price      float32   `json:"price" gorm:"default:0"`
	ImageUrl   string    `json:"image_url"`
	Author     string    `json:"author" gorm:"index"`
	Publisher  string    `json:"publisher" gorm:"index"`
	BookCode   string    `json:"book_code" gorm:"index:constraint_books,unique;not null"`
	CreatedBy  string    `json:"created_by"`
	UpdatedBy  string    `json:"updated_by"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	IsActive   bool      `json:"is_active" gorm:"default:true;index"`
}
