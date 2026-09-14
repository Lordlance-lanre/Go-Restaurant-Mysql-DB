package Models

import (
	"time"
	"gorm.io/gorm"
    "github.com/google/uuid"
)

type FoodItems struct{
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Food_ID    string    `json:"food_id" gorm:"uniqueIndex;not null"`
	Name       string    `json:"name" validate:"required,min=2,max=100"`
	Price      float64   `json:"price" validate:"required"`
	Food_image string    `json:"food_image" validate:"required"`
	MenuID     string    `json:"menu_id" gorm:"column:menu_id;not null" validate:"required"`
	Menu       Menu      `json:"menu,omitempty" gorm:"foreignKey:MenuID;references:Menu_ID"`
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