package application

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"learning-go/internal/config"
	"learning-go/internal/pkg/users/domain"
	"learning-go/internal/pkg/users/dtos"
	"learning-go/internal/pkg/users/enums"
	"learning-go/internal/pkg/users/infrastructure"
	"learning-go/internal/types"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
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

// Helpers
func GenerateTokens(claims jwt.MapClaims, key []byte, exp string) (string, error) {
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.GetEnv("MAIL_USERNAME", "example@example.com"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	port, _ := strconv.Atoi(config.GetEnv("MAIL_PORT", "587"))
	d := gomail.NewDialer(
		config.GetEnv("MAIL_HOST", "smtp.gmail.com"),
		port,
		config.GetEnv("MAIL_USERNAME", "example@example.com"),
		config.GetEnv("MAIL_PASSWORD", "password"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// RenderEmailTemplate renders an HTML email template with the provided data
func RenderEmailTemplate(templatePath string, data interface{}) (string, error) {
	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func parseExp(exp string) int64 {
	expInt, _ := strconv.ParseInt(exp, 10, 64)
	return expInt
}

func VerifyJwtToken(tokenString string, key string) (*jwt.Token, error) {
	issuer := config.GetEnv("JWT_ISSUER", "my-golang-app")
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}
	return token, nil
}

// Auth Functions
func (s *UserService) Register(registerDto dtos.RegisterDto) (*domain.User, error) {
	// var user = domain.User{}

	checkUser, err := s.repo.FindByEmail(registerDto.Email)
	// If no error, it means user was found - email already exists
	if err == nil && checkUser != nil {
		return nil, errors.New("user already exists")
	}

	// If error is something other than "record not found", return the error
	// "record not found" is expected for new registrations
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := HashPassword(registerDto.Password)
	if err != nil {
		return nil, err
	}

	user := jwt.MapClaims{
		"iss":       config.GetEnv("JWT_ISSUER", "my-golang-app"),
		"email":     registerDto.Email,
		"password":  string(hashedPassword),
		"username":  registerDto.Username,
		"birthdate": registerDto.Birthdate,
	}

	verifySecret := []byte(config.GetEnv("JWT_VERIFY_SECRET", "5f4dcc3b5aa765d61d8327deb882cf99acd7aa40a1f9c7df48adfaf48c69e4c3203a6a17acc60638a7bdc9103d0a499997180cb8d8"))
	verifyExpiry := config.GetEnv("JWT_VERIFY_EXPIRY", "30m")
	verifyToken, err := GenerateTokens(user, verifySecret, verifyExpiry)

	if err != nil {
		return nil, err
	}

	verifyLink := config.GetEnv("SERVER_HOST", "http://localhost:2000") + "/auth/verify?token=" + verifyToken

	// Prepare template data
	templateData := map[string]interface{}{
		"AppName":    config.GetEnv("APP_NAME", "My Golang App"),
		"Username":   registerDto.Username,
		"VerifyLink": verifyLink,
		"ExpiryTime": verifyExpiry,
	}

	// Render the email template
	templatePath := filepath.Join("internal", "pkg", "users", "templates", "verify_email.html")
	emailBody, err := RenderEmailTemplate(templatePath, templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render email template: %w", err)
	}

	err = SendEmail(registerDto.Email, "Email Verification", emailBody)
	if err != nil {
		return nil, err
	}

	// Return a temporary user object with email for the response
	// The actual user will be created after email verification
	return &domain.User{Email: registerDto.Email, Username: registerDto.Username, Birthdate: registerDto.Birthdate}, nil
}

func (s *UserService) VerifyEmail(tokenString string) (*domain.User, error) {
	verifySecret := (config.GetEnv("JWT_VERIFY_SECRET", "verify_secret_key"))

	token, err := VerifyJwtToken(tokenString, verifySecret)
	if err != nil {
		return nil, err
	}
	fmt.Println("Token verified successfully,", token)

	var email, password, username, birthdate, role string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email = claims["email"].(string)
		password = claims["password"].(string)
		username = claims["username"].(string)
		birthdate = claims["birthdate"].(string)
	}

	checkUser, err := s.repo.FindByEmail(email)
	// If no error, it means user was found - email already exists
	if err == nil && checkUser != nil {
		return nil, errors.New("user already verified, you can login now")
	}

	// If error is something other than "record not found", return the error
	// "record not found" is expected for new registrations
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	users, err := s.repo.FindAll(false)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		role = string(enums.Admin)
	}

	newUser := &domain.User{
		Email:     email,
		Password:  password,
		Username:  username,
		Birthdate: birthdate,
		Role:      role,
	}

	createdUser, err := s.repo.Create(newUser)
	if err != nil {
		return nil, err
	}
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
	accessExpiry := config.GetEnv("JWT_ACCESS_EXPIRY", "24h")

	refreshSecret := []byte(config.GetEnv("JWT_REFRESH_SECRET", "8ee21eb52bf23900b62daa6fc7f7b8a09e6548cc8ef5ecc894fa26ff49022b8b8327ba6ce9914b"))
	refreshExpiry := config.GetEnv("JWT_REFRESH_EXPIRY", "720h")

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

// only admin can see all users
func (s *UserService) GetAllUsers(includeDeleted bool) ([]domain.User, error) {
	users, err := s.repo.FindAll(includeDeleted)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// shared between admin and user
func (s *UserService) GetUserByID(id string, role string, includeDeleted bool) (*domain.User, error) {
	safeIncludeDeleted := false
	if role == string(enums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	user, err := s.repo.FindByID(id, safeIncludeDeleted)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

// only admin can restore users
func (s *UserService) RestoreUser(id string) (bool, error) {
	return s.repo.Restore(id)
}

func (s *UserService) DeleteUser(id string, role string, sub string) (bool, error) {
	if role != string(enums.Admin) && id != sub {
		return false, errors.New("only admin or owner can delete this account")
	}
	return s.repo.Delete(id)
}

func (s *UserService) UpdateUser(id string, userDto dtos.UpdateUserDto, sub string) (*domain.User, error) {
	if id != sub {
		return nil, errors.New("you can only update your own profile")
	}

	user, err := s.repo.FindByID(id, false)
	if err != nil {
		return nil, err
	}
	if userDto.Username != "" {
		user.Username = userDto.Username
	}

	if userDto.Birthdate != "" {
		user.Birthdate = userDto.Birthdate
	}

	return s.repo.Update(user)
}

// admin only
func (s *UserService) PaginatedUsers(page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool) (*types.Paginated[domain.User], error) {
	data, err := s.repo.Paginated(page, limit, search, searchField, order, sortBy, includeDeleted)
	if err != nil {
		return nil, err
	}

	data.Data = func(users []domain.User) []domain.User {
		for i := range users {
			users[i].Password = ""
		}
		return users
	}(data.Data)

	return data, nil
}

func (s *UserService) RefreshTokens(refreshTokenString string) (*UserCredentials, error) {
	refreshSecret := (config.GetEnv("JWT_REFRESH_SECRET", "refresh_secret_key"))

	token, err := VerifyJwtToken(refreshTokenString, refreshSecret)
	if err != nil {
		return nil, err
	}

	var sub, exp string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub = claims["sub"].(string)
		exp = fmt.Sprintf("%v", claims["exp"])
	}

	checkUser, err := s.repo.FindByID(sub, false)
	if checkUser == nil {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"iss":  config.GetEnv("JWT_ISSUER", "my-golang-app"),
		"sub":  checkUser.ID,
		"role": checkUser.Role,
	}

	accessSecret := []byte(config.GetEnv("JWT_ACCESS_SECRET", "acd7aa40a1f9c7df48adfaf48c69e4c3203a6a17acc60638a7bdc9103d0a499997180cb8d8"))
	accessExpiry := config.GetEnv("JWT_ACCESS_EXPIRY", "24h")
	var accessToken, refreshToken string

	accessToken, err = GenerateTokens(claims, accessSecret, accessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken = refreshTokenString

	if time.Unix(int64(parseExp(exp)), 0).Sub(time.Now()) < 24*time.Hour {
		refreshSecret := []byte(config.GetEnv("JWT_REFRESH_SECRET", "8ee21eb52bf23900b62daa6fc7f7b8a09e6548cc8ef5ecc894fa26ff49022b8b8327ba6ce9914b"))
		refreshExpiry := config.GetEnv("JWT_REFRESH_EXPIRY", "720h")
		refreshToken, err = GenerateTokens(claims, refreshSecret, refreshExpiry)
		if err != nil {
			return nil, err
		}
	}

	return &UserCredentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         checkUser,
	}, nil
}
