package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id              string
	Name            string
	Email           string
	Password        string
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (user *User) Empty() bool {
	return user.Id == ""
}

func (user *User) EmailVerified() bool {
	if user.EmailVerifiedAt != nil {
		if user.EmailVerifiedAt.IsZero() {
			return false
		}

		return true
	}

	return false
}
