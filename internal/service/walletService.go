package service
import (
	"context"
	"wallet/internal/store"
)

type WalletService struct {
	store store.Initializer
}

type WalletServiceInterface interface {
	GetWalletByUserID(ctx context.Context, userID uint64) (*store.Wallet, error)
}

func ExportWalletService(s store.Initializer) WalletServiceInterface {
	return &WalletService{store: s}
}


// WalletService provides methods to manage user wallets, including creating wallets, retrieving wallet information,
//  and updating wallet balances. It interacts with the underlying data store through the store.Initializer interface to perform these operations.	


func (s *WalletService) GetWalletByUserID(ctx context.Context, userID uint64) (*store.Wallet, error) {
	return s.store.Wallet.GetWalletByUserID(ctx, userID)
}
