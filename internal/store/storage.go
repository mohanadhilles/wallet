package store

import (
	"context"
	"gorm.io/gorm"
)

type Initializer struct {
	User interface {
		Me(ctx context.Context, userID uint64) (*User, error)
		GetUserByID(ctx context.Context, id uint64) (*User, error)
		CreateUser(ctx context.Context, user *User) (*User, error)
		GetUserByUsername(ctx context.Context, username string) (*User, error)
	}

	OTP interface {
		CreateOTP(ctx context.Context, otp *OTP) (*OTP, error)
		GetOTPByUserID(ctx context.Context, userID uint64) (*OTP, error)
		VerifyOTP(ctx context.Context, username string, code string) (bool, error)
	}

	Transaction interface {
		CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
		GetTransactionByID(ctx context.Context, id uint64) (*Transaction, error)
		GetTransactionsByUserID(ctx context.Context, userID uint64, page int, pageSize int, phone string, email string, username string) ([]Transaction, error)
	}

	Wallet interface {
		CreateWallet(ctx context.Context, wallet *Wallet) (*Wallet, error)
		GetWalletByUserID(ctx context.Context, userID uint64) (*Wallet, error)
		UpdateWalletBalance(ctx context.Context, userID uint64, amount float64) error
	}
}

func LoadInitializer(db *gorm.DB) Initializer {
	return Initializer{
		User:        &storeUser{db: db},
		OTP:         &storeOTP{db: db},
		Transaction: &storeTransaction{db: db},
		Wallet:      &storeWallet{db: db},
	}
}
