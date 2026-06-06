package service

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"
	"wallet/internal/store"
	"wallet/provider/jwt"
	"wallet/provider/mailer"
)

type UserService struct {
	store store.Initializer
}

type UserServiceInterface interface {
	GenerateOTP(ctx context.Context, username string) (*store.OTP, error)
	VerifyOTP(ctx context.Context, userID uint64, code string) (string, error)
	Me(ctx context.Context, userID uint64) (*store.User, error)
	CreateUser(ctx context.Context, user *store.User) (*store.User, error)
	GetUserByID(ctx context.Context, id uint64) (*store.User, error)
	GetUserByUsername(ctx context.Context, username string) (*store.User, error)
}

func ExportUserService(s store.Initializer) UserServiceInterface {
	return &UserService{store: s}
}

// UserService provides methods to manage user accounts, including creating users, retrieving user information,
//
//	and handling OTP generation and verification for authentication.
//	It interacts with the underlying data store through the store.Initializer interface to perform these operations.
func (s *UserService) GenerateOTP(ctx context.Context, username string) (*store.OTP, error) {
	var user *store.User
	user, err := s.store.User.GetUserByUsername(ctx, username)
	if err != nil {
		user, err = s.store.User.CreateUser(ctx, &store.User{Username: username})
		if err != nil {
			return nil, err
		}
		go (func() {
			err = mailer.Send(user.Username, "Login Notification", "You have successfully logged in to your account. If this wasn't you, please secure your account immediately. Welcome back!")
			if err != nil {
				log.Println("Error sending login notification email:", err)
				return
			}
			_, err := s.store.Wallet.CreateWallet(ctx, &store.Wallet{UserID: user.ID, Balance: 0})
			if err != nil {
				log.Println("Error creating wallet for user:", err)
				return
			}
		})()
	}

	otp := &store.OTP{}
	otp.Username = user.Username
	otp.Code = generateOTPCode()

	err = mailer.Send(user.Username, "OTP Code", "Your OTP code is: "+otp.Code+". If this wasn't you, please secure your account immediately.")
	if err != nil {
		log.Println("Error sending OTP email:", err)
		return nil, err
	}

	return s.store.OTP.CreateOTP(ctx, otp)
}

func (s *UserService) VerifyOTP(ctx context.Context, userID uint64, code string) (string, error) {
	user, err := s.store.User.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	_, err = s.store.OTP.VerifyOTP(ctx, user.Username, code)
	if err != nil {
		return "", err
	}

	// inside your login handler, after credentials are verified

	jwtService := jwt.JWTServices(os.Getenv("JWT_SECRET_KEY"), 24*time.Hour)
	token, err := jwtService.GenerateToken(userID)
	if err != nil {
		return err.Error(), err
	}
	return token, nil
}

func (s *UserService) Me(ctx context.Context, userID uint64) (*store.User, error) {
	return s.store.User.Me(ctx, userID)
}

func (s *UserService) CreateUser(ctx context.Context, user *store.User) (*store.User, error) {
	return s.store.User.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint64) (*store.User, error) {
	return s.store.User.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	return s.store.User.GetUserByUsername(ctx, username)
}

// generate 6 digit OTP code
func generateOTPCode() string {
	const otpLength = 6
	const digits = "0123456789"
	otp := make([]byte, otpLength)
	for i := range otp {
		otp[i] = digits[rand.Int63n(int64(len(digits)))]
	}
	return string(otp)
}
