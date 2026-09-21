package Models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
    ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Order_ID   string    `json:"order_id" gorm:"type:char(36);uniqueIndex;not null"`
    User_ID    uint      `json:"user_id" gorm:"column:user_id;not null;index" validate:"required"`
    User       User      `json:"user,omitempty" gorm:"foreignKey:User_ID;references:ID"`
   FoodRef   string    `json:"food_id" gorm:"column:food_id;type:char(36);not null" validate:"required"`
FoodItems FoodItems `json:"food,omitempty" gorm:"foreignKey:FoodRef;references:Food_ID"`
    
    Order_Date time.Time `json:"order_date" gorm:"column:order_date;type:date;not null" validate:"required"`
    Created_at time.Time `json:"created_at"`
    Updated_at time.Time `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
// Food_ID    string    `json:"food_id" gorm:"uniqueIndex;not null"`

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	o.Order_ID = uuid.NewString()
	return nil
}
