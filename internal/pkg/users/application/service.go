package application

import (
	"errors"
	"fmt"
	"learning-go/internal/config"
	httpError "learning-go/internal/http"
	"learning-go/internal/pkg/users/domain"
	"learning-go/internal/pkg/users/dtos"
	"learning-go/internal/pkg/users/enums"
	"learning-go/internal/pkg/users/infrastructure"
	"learning-go/internal/types"
	"learning-go/internal/utils"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	repo            *infrastructure.UserRepository
	jwtService      *utils.JwtService
	emailService    *utils.EmailService
	passwordService *utils.PasswordService
	serverConfig    *config.ServerConfig
}

type UserCredentials struct {
	AccessToken  string
	RefreshToken string
	User         *domain.User
}

// Constructor liked
func NewUserService(repo *infrastructure.UserRepository, JwtService *utils.JwtService, EmailService *utils.EmailService, PasswordService *utils.PasswordService, serverConfig *config.ServerConfig) *UserService {
	return &UserService{
		repo:            repo,
		jwtService:      JwtService,
		emailService:    EmailService,
		passwordService: PasswordService,

		serverConfig: serverConfig,
	}
}

// Auth Functions
func (s *UserService) Register(registerDto dtos.RegisterDto) (*domain.User, error) {
	// var user = domain.User{}

	checkUser, err := s.repo.FindByEmail(registerDto.Email)
	// If no error, it means user was found - email already exists
	if err == nil && checkUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// If error is something other than "record not found", return the error
	// "record not found" is expected for new registrations
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// hashedPassword, err := s.passwordService.HashPassword(registerDto.Password)
	// if err != nil {
	// 	return nil, err
	// }

	user := &domain.User{
		Email:     registerDto.Email,
		Password:  registerDto.Password,
		Username:  registerDto.Username,
		Birthdate: registerDto.Birthdate,
	}

	verifyToken, verifyExpiry, err := s.jwtService.GenerateVerifyToken(user)

	if err != nil {
		return nil, err
	}

	verifyLink := s.serverConfig.Host + "/auth/verify?token=" + verifyToken

	duration, err := time.ParseDuration(verifyExpiry)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(duration)

	// Prepare template data
	templateData := map[string]interface{}{
		"AppName":    s.serverConfig.AppName,
		"Username":   registerDto.Username,
		"VerifyLink": verifyLink,
		"ExpiryTime": expiresAt.Format(time.RFC1123),
	}

	// Render the email template
	templatePath := filepath.Join("internal", "pkg", "users", "templates", "verify_email.html")
	emailBody, err := s.emailService.RenderEmailTemplate(templatePath, templateData)
	if err != nil {
		return nil, domain.ErrFailedToRenderHTML
	}

	err = s.emailService.SendEmail(registerDto.Email, "Email Verification", emailBody)
	if err != nil {
		return nil, err
	}

	// Return a temporary user object with email for the response
	// The actual user will be created after email verification
	return &domain.User{Email: registerDto.Email, Username: registerDto.Username, Birthdate: registerDto.Birthdate}, nil
}

func (s *UserService) VerifyEmail(tokenString string) (*domain.User, error) {

	verifyVerificationClaims, err := s.jwtService.VerifyVerificationToken(tokenString)
	if verifyVerificationClaims == (utils.VerifyEmailClaims{}) {
		return nil, domain.ErrInvalidVerificationToken
	}
	if err != nil {
		return nil, err
	}

	// fmt.Println("Token verified successfully,", verifyVerificationClaims)

	var email, password, username, birthdate, role string

	email = verifyVerificationClaims.Email
	password = verifyVerificationClaims.Password
	username = verifyVerificationClaims.Username
	birthdate = verifyVerificationClaims.Birthdate

	hashedPassword, err := s.passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	role = string(enums.User)

	checkUser, err := s.repo.FindByEmail(email)
	// If no error, it means user was found - email already exists
	if err == nil && checkUser != nil {
		return nil, domain.ErrUserAlreadyExists
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
		Password:  hashedPassword,
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
		return nil, domain.ErrInvalidCredentials
	}

	err = s.passwordService.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
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
		return false, httpError.ErrForbidden
	}
	return s.repo.Delete(id)
}

func (s *UserService) UpdateUser(id string, userDto dtos.UpdateUserDto, sub string) (*domain.User, error) {
	if id != sub {
		return nil, httpError.ErrForbidden
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
	allowed := map[string]bool{"email": true, "username": true, "created_at": true}
	if !allowed[searchField] {
		searchField = "email"
	}
	if !allowed[sortBy] {
		sortBy = "created_at"
	}

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
	claims, err := s.jwtService.VerifyRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	var sub, exp string

	sub = claims["sub"].(string)
	exp = fmt.Sprintf("%v", claims["exp"])

	checkUser, err := s.repo.FindByID(sub, false)
	if checkUser == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	var accessToken, refreshToken string

	accessToken, err = s.jwtService.GenerateAccessToken(checkUser)
	if err != nil {
		return nil, err
	}

	refreshToken = refreshTokenString

	if time.Unix(int64(s.jwtService.ParseExp(exp)), 0).Sub(time.Now()) < 24*time.Hour {
		refreshToken, err = s.jwtService.GenerateRefreshToken(checkUser)
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
