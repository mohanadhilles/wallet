package store

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type OTP struct {
	ID        uint64         `gorm:"primaryKey"`
	Username  string         `gorm:"size:255;index"`
	Code      string         `gorm:"size:6"`
	Used      bool           `gorm:"default:false"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (o *OTP) TableName() string {
	return "otps"
}

type storeOTP struct {
	db *gorm.DB
}

func (s *storeOTP) CreateOTP(ctx context.Context, otp *OTP) (*OTP, error) {
	err := s.db.WithContext(ctx).Create(otp).Error
	return otp, err
}

func (s *storeOTP) GetOTPByUserID(ctx context.Context, userID uint64) (*OTP, error) {
	var otp OTP
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (s *storeOTP) VerifyOTP(ctx context.Context, username string, code string) (bool, error) {
	var otp OTP
	err := s.db.WithContext(ctx).Where("username = ? AND code = ? AND used = false", username, code).Order("id DESC").First(&otp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	otp.Used = true
	err = s.db.WithContext(ctx).Save(&otp).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
