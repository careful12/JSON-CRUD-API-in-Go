package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title string
	Body  string
}

// gorm.Model definition
// type Model struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
//   }
