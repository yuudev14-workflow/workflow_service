package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GetUserByEmailOrUsername(usernameOrEmail string) (*models.User, error)
	ValidateUserSignUp(username string, email string) error
	VerifyUser(form dto.LoginForm) (*models.User, error)
	CheckUserByUsername(username string) error
	CheckUserByEmail(email string) error
	CreateUser(form dto.SignupForm) (*models.User, error)
}

type AuthServiceImpl struct {
	*sqlx.DB
}

// Auth Service Constructor
func NewUserService(db *sqlx.DB) AuthService {
	return &AuthServiceImpl{
		DB: db,
	}
}

// Get User by providing username
func (a *AuthServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	return repository.DbSelectOne[models.User](
		a.DB,
		"SELECT id, username, email, first_name, last_name from users WHERE username=$1",
		username,
	)
}

// VerifyUser implements UserService.
func (a *AuthServiceImpl) VerifyUser(form dto.LoginForm) (*models.User, error) {
	user, usernameError := repository.DbSelectOne[models.User](
		a.DB,
		"SELECT * from users WHERE email=$1 OR username=$1",
		form.Username,
	)

	if usernameError != nil {
		return nil, usernameError
	}

	if user == nil {
		return nil, fmt.Errorf("user is not found")
	}

	isNotMatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))

	if isNotMatch != nil {
		return nil, fmt.Errorf("password is not correct")
	}
	return user, nil
}

// ValidateUserSignUp implements UserService.
func (a *AuthServiceImpl) ValidateUserSignUp(username string, email string) error {
	// check if username already exist
	usernameError := a.CheckUserByUsername(username)

	if usernameError != nil {
		return usernameError
	}

	// check if email already exist
	emailError := a.CheckUserByEmail(email)

	if emailError != nil {
		return emailError
	}
	return nil
}

// CheckUserByEmail implements UserService.
func (a *AuthServiceImpl) CheckUserByEmail(email string) error {
	// check if email already exist
	user, emailError := repository.DbSelectOne[models.User](
		a.DB,
		"SELECT * from users WHERE email=$1",
		email,
	)

	if emailError != nil {
		return emailError
	}

	if user != nil {
		return fmt.Errorf("email already exist")
	}
	return nil
}

// CheckUserByUsername implements UserService.
func (a *AuthServiceImpl) CheckUserByUsername(username string) error {
	user, usernameError := repository.DbSelectOne[models.User](
		a.DB,
		"SELECT id, username, email, first_name, last_name from users WHERE username=$1",
		username,
	)

	if usernameError != nil {
		return usernameError
	}

	if user != nil {
		return fmt.Errorf("username already exist")
	}
	return nil
}

// Get User by providing username or email
func (a *AuthServiceImpl) GetUserByEmailOrUsername(usernameOrEmail string) (*models.User, error) {
	return repository.DbSelectOne[models.User](
		a.DB,
		"SELECT * from users WHERE email=$1 OR username=$1",
		usernameOrEmail,
	)
}

// create user
func (a *AuthServiceImpl) CreateUser(form dto.SignupForm) (*models.User, error) {
	// encode password
	newPassword, passwordErr := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)

	if passwordErr != nil {
		return nil, passwordErr
	}
	// save user
	excryptedPassword := string(newPassword)

	_, err := a.DB.Exec(`INSERT INTO users (email, username, password) VALUES ($1, $2, $3)`, form.Email, form.Username, excryptedPassword)
	if err != nil {
		return nil, err
	}
	user, err := a.GetUserByUsername(form.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
