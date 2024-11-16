// models/post.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           `gorm:"primaryKey"`
    Title     string         `gorm:"size:255;not null"`
    Content   string         `gorm:"type:text"`
    Author    string         `gorm:"size:100"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

