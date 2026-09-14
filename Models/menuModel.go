package Models

import (
	"time"
	"gorm.io/gorm"
    "github.com/google/uuid"
)

type Menu struct {
   	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Menu_ID    string    `gorm:"column:menu_id;unique;type:char(36)" json:"menu_id"`
    Name       string    `gorm:"column:name;not null" json:"name"`
    Category   string    `gorm:"column:category" json:"category"`
    Start_Date string    `gorm:"column:start_date" json:"start_date"`
    End_Date   string    `gorm:"column:end_date" json:"end_date"`
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
}

// Auto-generate UUID before creating
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
    m.Menu_ID = uuid.NewString()
    return nil
}