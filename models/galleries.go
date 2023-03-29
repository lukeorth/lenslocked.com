package models

import "gorm.io/gorm"

// Gallery is our image container resources that visitors
// view
type Gallery struct {
    gorm.Model
    UserID uint `gorm:"not_null;index"`
    Title string `gorm:"not_null"`
}
