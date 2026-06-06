package store

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID         uint64         `gorm:"primaryKey"`
	SenderID   uint64         `gorm:"index"`
	ReceiverID uint64         `gorm:"index"`
	Amount     float64        `gorm:"type:decimal(10,2)"`
	Currency   string         `gorm:"size:3"`
	Status     string         `gorm:"size:20;default:'pending'"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Sender   *User `gorm:"foreignKey:SenderID"`
	Receiver *User `gorm:"foreignKey:ReceiverID"`
}

func (t *Transaction) TableName() string {
	return "transactions"
}

type storeTransaction struct {
	db *gorm.DB
}

func (s *storeTransaction) CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	err := s.db.WithContext(ctx).Create(transaction).Error
	return transaction, err
}

func (s *storeTransaction) GetTransactionByID(ctx context.Context, id uint64) (*Transaction, error) {
	var transaction Transaction
	err := s.db.WithContext(ctx).Preload("Sender").Preload("Receiver").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *storeTransaction) GetTransactionsByUserID(ctx context.Context, userID uint64, page int, pageSize int, phone string, email string, username string) ([]Transaction, error) {
	var transactions []Transaction
	q := s.db.WithContext(ctx).Preload("Sender").Preload("Receiver").
		Joins("JOIN users s ON s.id = transactions.sender_id").
		Joins("JOIN users r ON r.id = transactions.receiver_id").
		Where("transactions.sender_id = ? OR transactions.receiver_id = ?", userID, userID)

	if phone != "" {
		q = q.Where("s.phone LIKE ? OR r.phone LIKE ?", "%"+phone+"%", "%"+phone+"%")
	}
	if email != "" {
		q = q.Where("s.email LIKE ? OR r.email LIKE ?", "%"+email+"%", "%"+email+"%")
	}
	if username != "" {
		q = q.Where("s.username LIKE ? OR r.username LIKE ?", "%"+username+"%", "%"+username+"%")
	}

	err := q.Order("transactions.id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
