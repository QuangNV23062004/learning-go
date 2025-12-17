package application

import (
	"errors"
	"learning-go/internal/config"
	"learning-go/internal/pkg/users/domain"
	"learning-go/internal/pkg/users/dtos"
	"learning-go/internal/pkg/users/enums"
	"learning-go/internal/pkg/users/infrastructure"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *infrastructure.UserRepository
}

type UserCredentials struct {
	AccessToken  string
	RefreshToken string
	User         *domain.User
}

// Constructor liked
func NewUserService(repo *infrastructure.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Helper
func GenerateTokens(claims jwt.MapClaims, key []byte, exp string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Functions
func (s *UserService) Register(registerDto dtos.RegisterDto) (*domain.User, error) {
	var user = domain.User{}

	checkUser, err := s.repo.FindByEmail(registerDto.Email)
	if err == nil {
		return nil, err
	}

	if checkUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := HashPassword(registerDto.Password)
	if err != nil {
		return nil, err
	}

	users, err := s.repo.FindAll(true)
	if err != nil {
		return nil, err
	}

	role := enums.User
	if len(users) == 0 {
		role = enums.Admin
	}

	user = domain.User{
		Username:  registerDto.Username,
		Email:     registerDto.Email,
		Password:  string(hashedPassword),
		Birthdate: registerDto.Birthdate,
		Role:      string(role),
	}

	createdUser, err := s.repo.Create(&user)
	if err != nil {
		return nil, err
	}

	createdUser.Password = ""

	return createdUser, nil
}

func (s *UserService) Login(loginDto dtos.LoginDto) (*UserCredentials, error) {
	user, err := s.repo.FindByEmail(loginDto.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"iss":  config.GetEnv("JWT_ISSUER", "my-golang-app"),
		"sub":  user.ID,
		"role": user.Role,
	}

	accessSecret := []byte(config.GetEnv("JWT_ACCESS_SECRET", "acd7aa40a1f9c7df48adfaf48c69e4c3203a6a17acc60638a7bdc9103d0a499997180cb8d8"))
	accessExpiry := config.GetEnv("JWT_ACCESS_EXPIRY", "1d")

	refreshSecret := []byte(config.GetEnv("JWT_REFRESH_SECRET", "8ee21eb52bf23900b62daa6fc7f7b8a09e6548cc8ef5ecc894fa26ff49022b8b8327ba6ce9914b"))
	refreshExpiry := config.GetEnv("JWT_REFRESH_EXPIRY", "30d")

	accessToken, err := GenerateTokens(claims, accessSecret, accessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateTokens(claims, refreshSecret, refreshExpiry)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &UserCredentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
