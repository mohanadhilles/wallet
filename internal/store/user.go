package store

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint64         `gorm:"primaryKey"`
	Username     string         `gorm:"unique;not null"`
	Email        *string        `gorm:"unique"`
	Phone        *string        `gorm:"unique"`
	Status       bool           `gorm:"default:false"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Wallet       Wallet         `gorm:"foreignKey:UserID"`
	Transactions []Transaction  `gorm:"foreignKey:SenderID;foreignKey:ReceiverID"`
}

func (u *User) TableName() string {
	return "users"
}

type storeUser struct {
	db *gorm.DB
}

func (s *storeUser) Me(ctx context.Context, userID uint64) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).Preload("Wallet").Preload("Transactions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *storeUser) CreateUser(ctx context.Context, user *User) (*User, error) {
	err := s.db.WithContext(ctx).Where("username = ?", user.Username).FirstOrCreate(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *storeUser) GetUserByID(ctx context.Context, id uint64) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).Preload("Wallet").Preload("Transactions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *storeUser) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).Preload("Wallet").Preload("Transactions").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
