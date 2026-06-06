package store

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	ID        uint64         `gorm:"primaryKey"`
	UserID    uint64         `gorm:"uniqueIndex"`
	Balance   float64        `gorm:"type:decimal(10,2)"`
	Currency  string         `gorm:"size:3"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      *User          `gorm:"foreignKey:UserID"`
}

func (w *Wallet) TableName() string {
	return "wallets"
}

type storeWallet struct {
	db *gorm.DB
}

func (s *storeWallet) CreateWallet(ctx context.Context, wallet *Wallet) (*Wallet, error) {
	err := s.db.WithContext(ctx).Create(wallet).Error
	return wallet, err
}

func (s *storeWallet) GetWalletByUserID(ctx context.Context, userID uint64) (*Wallet, error) {
	var wallet Wallet
	err := s.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

/**
Breaking it down:
.Model(&Wallet{}) — targets the wallets table
.Where("user_id = ?", userID) — filters to the specific user's wallet
.UpdateColumn(...) — updates a single column, skipping GORM hooks and UpdatedAt auto-update
gorm.Expr("balance + ?", amount) — uses a raw SQL expression so the addition happens in the database, not in Go code
The key reason to use gorm.Expr here instead of reading the balance first and adding in Go:

Race condition safe — if two transactions run simultaneously, balance = balance + 50 computed by the DB is atomic, whereas read→add→write in Go would cause one update to overwrite the other
amount can be negative (for deductions) or positive (for deposits)
The tradeoff of using .UpdateColumn vs .Update is that it bypasses GORM's BeforeUpdate/AfterUpdate hooks and won't touch UpdatedAt — worth knowing if you rely on those.
*/

func (s *storeWallet) UpdateWalletBalance(ctx context.Context, userID uint64, amount float64) error {
	return s.db.WithContext(ctx).Model(&Wallet{}).Where("user_id = ?", userID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}
