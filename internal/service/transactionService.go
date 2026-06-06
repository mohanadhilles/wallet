package service

import (
	"context"
	"errors"
	"wallet/internal/store"
	"wallet/provider/mailer"
	"strconv"

"log"
)

type TransactionService struct {
	store store.Initializer
}


type TransactionServiceInterface interface {
	CreateTransaction(ctx context.Context, transaction *store.Transaction) (*store.Transaction, error)
	GetTransactionByID(ctx context.Context, id uint64) (*store.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uint64, page int, pageSize int, phone string, email string, username string) ([]store.Transaction, error)
}

func ExportTransactionService(s store.Initializer) TransactionServiceInterface {
	return &TransactionService{store: s}
}


// TransactionService provides methods to manage transactions between users, including creating transactions,
//  retrieving transaction details by ID, and listing transactions for a specific user with pagination support.
//  It interacts with the underlying data store through the store.Initializer interface to perform these operations.

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *store.Transaction) (*store.Transaction, error) {
	walletSender, err := s.store.Wallet.GetWalletByUserID(ctx, transaction.SenderID)
	if err != nil {
		log.Printf("Error fetching sender wallet: %v", err)
		return nil, err
	}
	if walletSender.Balance < transaction.Amount {
		return nil, errors.New("insufficient balance")
	}

	walletReceiver, err := s.store.Wallet.GetWalletByUserID(ctx, transaction.ReceiverID)
	if err != nil {
		return nil, errors.New("receiver wallet not found")
	}

	err = s.store.Wallet.UpdateWalletBalance(ctx, transaction.SenderID, -transaction.Amount)
	if err != nil {
		return nil, err
	}

	err = s.store.Wallet.UpdateWalletBalance(ctx, transaction.ReceiverID, transaction.Amount)
	if err != nil {
		return nil, err
	}

	transaction.Status = "completed"
	walletReceiver.Balance += transaction.Amount
	walletSender.Balance -= transaction.Amount

	go func() {
		err = mailer.Send(walletSender.User.Username, "Transaction Alert", "You have sent " + strconv.FormatFloat(transaction.Amount, 'f', -1, 64) + " to " + walletReceiver.User.Username + ". Your new balance is " + strconv.FormatFloat(walletSender.Balance, 'f', -1, 64) + ". If this wasn't you, please secure your account immediately.")
		if err != nil {
			log.Printf("Error sending email to sender: %v", err)
			return
		}
		err = mailer.Send(walletReceiver.User.Username, "Transaction Alert", "You have received " + strconv.FormatFloat(transaction.Amount, 'f', -1, 64) + " from " + walletSender.User.Username + ". Your new balance is " + strconv.FormatFloat(walletReceiver.Balance, 'f', -1, 64) + ". If this wasn't you, please secure your account immediately.")
		if err != nil {
			log.Printf("Error sending email to receiver: %v", err)
			return
		}
	}()
	
	return s.store.Transaction.CreateTransaction(ctx, transaction)
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id uint64) (*store.Transaction, error) {
	return s.store.Transaction.GetTransactionByID(ctx, id)
}

func (s *TransactionService) GetTransactionsByUserID(ctx context.Context, userID uint64, page int, pageSize int, phone string, email string, username string) ([]store.Transaction, error) {
	return s.store.Transaction.GetTransactionsByUserID(ctx, userID, page, pageSize, phone, email, username)
}
