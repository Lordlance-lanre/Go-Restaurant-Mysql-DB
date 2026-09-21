package Models

import (
	"time"
	"gorm.io/gorm"
    "github.com/google/uuid"
)

type FoodItems struct {
    ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Food_ID    string    `json:"food_id" gorm:"type:char(36);uniqueIndex;not null"`
    Name       string    `json:"name" gorm:"type:varchar(100)" validate:"required,min=2,max=100"`
    Price      float64   `json:"price" validate:"required"`
    Food_image string    `json:"food_image" gorm:"type:varchar(255)" validate:"required"`
       MenuRef    string    `json:"menu_id" gorm:"column:menu_id;type:varchar(36);not null" validate:"required"`
    
    Menu       Menu      `json:"menu,omitempty" gorm:"foreignKey:MenuRef;references:Menu_ID"`
    Start_Date string    `json:"start_date" gorm:"column:start_date;type:date;not null"`
    End_Date   string    `json:"end_date" gorm:"column:end_date;type:date;not null"`
    Created_at time.Time `json:"created_at"`
    Updated_at time.Time `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (f *FoodItems) BeforeCreate(tx *gorm.DB) error {
	f.Food_ID = uuid.NewString()
	return nil
}