package service

import (
	"errors"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         *repository.UserRepository
	simpananRepo *repository.SimpananRepository
}

func NewAuthService(repo *repository.UserRepository, simpananRepo *repository.SimpananRepository) *AuthService {
	return &AuthService{repo: repo, simpananRepo: simpananRepo}
}

func (s *AuthService) Register(user *model.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	// Create user first
	if err := s.repo.Create(user); err != nil {
		return err
	}

	// Initialize user wallets (3 types)
	if err := s.simpananRepo.InitializeUserWallets(user.ID); err != nil {
		// If wallet initialization fails, we should rollback user creation
		// For now, we'll just return the error
		return err
	}

	return nil
}

func (s *AuthService) Login(email, password, jwtSecret string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return GenerateToken(user.ID, jwtSecret)
}

func GenerateToken(userID uint, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *AuthService) GetUserWithRole(id uint) (*model.User, error) {
	return s.repo.FindByIDWithRole(id)
}
