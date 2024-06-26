package domain

import (
	"time"
)

type User struct {
	UserID           uint `gorm:"primarykey"`
	Username         string
	Email            *string
	PasswordHash     string
	RegistrationDate time.Time
	LastLoginDate    time.Time
	IsActive         bool
}

type IUserRepository interface {
	FindAll() ([]User, error)
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByCredential(username string) (*User, error)
	IsExistsByUsername(username string) (bool, error)
	Update(user *User) error
	Delete(id uint) error
}
