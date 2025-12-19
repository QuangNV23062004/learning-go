package utils

import (
	"learning-go/internal/config"
	"learning-go/internal/pkg/users/domain"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	config *config.JWTConfig
}

type VerifyEmailClaims struct {
	Email     string
	Password  string
	Username  string
	Birthdate string
	ExpiresAt time.Time
}

func NewJwtService(config *config.JWTConfig) *JwtService {
	return &JwtService{
		config: config,
	}
}

func (j *JwtService) GetJwtToken(tokenString string, key string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithIssuer(j.config.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (j *JwtService) PrepareAuthClaims(user *domain.User) jwt.MapClaims {
	claims := jwt.MapClaims{
		"iss":  j.config.Issuer,
		"sub":  user.ID,
		"role": user.Role,
	}
	return claims
}

func (j *JwtService) PrepareVerifyClaims(user *domain.User) jwt.MapClaims {
	claims := jwt.MapClaims{
		"iss":       j.config.Issuer,
		"email":     user.Email,
		"password":  user.Password,
		"username":  user.Username,
		"birthdate": user.Birthdate,
	}

	return claims
}

func (j *JwtService) GenerateTokens(claims jwt.MapClaims, key []byte, exp string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	expDuration, err := time.ParseDuration(exp)

	if err != nil {
		return "", err
	}

	claims["exp"] = time.Now().Add(expDuration).Unix()

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtService) ParseExp(exp string) int64 {
	expInt, _ := strconv.ParseInt(exp, 10, 64)
	return expInt
}

func (j *JwtService) GenerateAccessToken(user *domain.User) (string, error) {

	claims := j.PrepareAuthClaims(user)

	accessSecret := []byte(j.config.AccessSecret)
	accessExpiry := j.config.AccessExpiry

	accessToken, err := j.GenerateTokens(claims, accessSecret, accessExpiry)
	if err != nil {
		return "", err
	}

	return accessToken, nil

}

func (j *JwtService) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := j.PrepareAuthClaims(user)

	refreshSecret := []byte(j.config.RefreshSecret)
	refreshExpiry := j.config.RefreshExpiry

	refreshToken, err := j.GenerateTokens(claims, refreshSecret, refreshExpiry)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (j *JwtService) GenerateVerifyToken(user *domain.User) (string, string, error) {
	claims := j.PrepareVerifyClaims(user)

	verifySecret := []byte(j.config.VerifySecret)
	verifyExpiry := j.config.VerifyExpiry
	token, err := j.GenerateTokens(claims, verifySecret, verifyExpiry)
	if err != nil {
		return "", "", err
	}

	return token, verifyExpiry, nil
}

func (j *JwtService) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {

	token, err := j.GetJwtToken(tokenString, j.config.AccessSecret)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil
}

func (j *JwtService) VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, err := j.GetJwtToken(tokenString, j.config.RefreshSecret)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, nil
}

func (j *JwtService) VerifyVerificationToken(tokenString string) (VerifyEmailClaims, error) {
	token, err := j.GetJwtToken(tokenString, j.config.VerifySecret)
	if err != nil {
		return VerifyEmailClaims{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return VerifyEmailClaims{
			Email:     claims["email"].(string),
			Password:  claims["password"].(string),
			Username:  claims["username"].(string),
			Birthdate: claims["birthdate"].(string),
			ExpiresAt: time.Unix(int64(claims["exp"].(float64)), 0),
		}, nil
	}
	return VerifyEmailClaims{}, nil
}

func (j *JwtService) GetAccessToken(tokenString string) (*jwt.Token, error) {
	accessSecret := []byte(j.config.AccessSecret)

	return j.GetJwtToken(tokenString, string(accessSecret))
}
