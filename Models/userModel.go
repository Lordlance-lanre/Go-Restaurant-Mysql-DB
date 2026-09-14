package Models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	FirstName string    `gorm:"" json:"first_name"`
	LastName  string    `gorm:"" json:"last_name"`
	Email     string    `gorm:"unique" json:"email"`
	Phone     string    `gorm:"unique" json:"phone"`
	Password  []byte    `gorm:"column:password" json:"-"`
	CreatedAt time.Time `gorm:"" json:"created_at"`
	UpdatedAt time.Time `gorm:"" json:"updated_at"`
}

func (user *User) SetPassword(password string) {
	// hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return err == nil
}
