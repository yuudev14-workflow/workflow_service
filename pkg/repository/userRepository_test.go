package repository_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/yuudev14-workflow/workflow-service/environment"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
)

var (
	sqlxDB       *sqlx.DB
	mock         sqlmock.Sqlmock
	repo         repository.UserRepository
	expectedUser *models.User
)

func TestMain(m *testing.M) {
	environment.Setup()
	logging.Setup("DEBUG")

	// Create a new mock database connection
	mockDB, sqlmock, mockErr := sqlmock.New()
	mock = sqlmock
	if mockErr != nil {
		logging.Sugar.Fatalf("an error '%s' was not expected when opening a stub database connection", mockErr)
	}
	defer mockDB.Close()

	// Wrap the mock database with sqlx
	sqlxDB = sqlx.NewDb(mockDB, "sqlmock")

	// Create an instance of UserRepositoryImpl with the mock database
	repo = repository.NewUserRepository(sqlxDB)
	// Set up the expectation
	id, _ := uuid.NewUUID()
	expectedUser = &models.User{
		ID:       id,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}

func setMockRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email)
}

func checkExpectedOutput(t *testing.T, err error, expected interface{}, output interface{}) {
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}
func TestGetUserByEmailOrUsername(t *testing.T) {

	tests := []struct {
		name     string
		username string
		expected *models.User
	}{
		{
			name:     "user exist",
			username: "testuser",
			expected: expectedUser,
		},
		{
			name:     "user does not exist",
			username: "testuser1",
			expected: nil,
		},
	}

	// Define the expected query and result
	expectedQuery := "SELECT \\* from users WHERE email=\\$1 OR username=\\$1"
	rows := setMockRows()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery(expectedQuery).
				WithArgs(test.username).
				WillReturnRows(rows)
			// Call the method being tested
			user, err := repo.GetUserByEmailOrUsername(test.username)
			// Assert the results
			checkExpectedOutput(t, err, test.expected, user)
		})

	}

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {

	id2, _ := uuid.NewUUID()
	tests := []struct {
		name     string
		id       uuid.UUID
		expected *models.User
	}{
		{
			name:     "user exist",
			id:       expectedUser.ID,
			expected: expectedUser,
		},
		{
			name:     "user doenst exist",
			id:       id2,
			expected: nil,
		},
	}

	expectedQuery := "SELECT id, username, email, first_name, last_name from users WHERE id=\\$1"
	rows := setMockRows()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery(expectedQuery).
				WithArgs(test.id).
				WillReturnRows(rows)
			user, err := repo.GetUserByID(test.id)
			checkExpectedOutput(t, err, test.expected, user)

		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername(t *testing.T) {
	// create the tests
	tests := []struct {
		name     string
		username string
		expected *models.User
	}{
		{
			name:     "user exist",
			username: expectedUser.Username,
			expected: expectedUser,
		},
		{
			name:     "user does not exist",
			username: "testuser1",
			expected: nil,
		},
	}

	expectedQuery := "SELECT id, username, email, first_name, last_name from users WHERE username=\\$1"

	rows := setMockRows()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery(expectedQuery).
				WithArgs(test.username).
				WillReturnRows(rows)
			user, err := repo.GetUserByUsername(test.username)
			checkExpectedOutput(t, err, test.expected, user)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	// create the tests
	tests := []struct {
		name     string
		email    string
		expected *models.User
	}{
		{
			name:     "user exist",
			email:    expectedUser.Email,
			expected: expectedUser,
		},
		{
			name:     "user does not exist",
			email:    "testuser1",
			expected: nil,
		},
	}

	expectedQuery := "SELECT \\* from users WHERE email=\\$1"

	rows := setMockRows()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery(expectedQuery).
				WithArgs(test.email).
				WillReturnRows(rows)
			user, err := repo.GetUserByEmail(test.email)
			checkExpectedOutput(t, err, test.expected, user)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserWhereUserIsNotAvaliable(t *testing.T) {
	test := struct {
		user     models.User
		expected *models.User
	}{
		user: models.User{
			Email:    "test111@gmail.com",
			Username: "test111",
			Password: "password",
		},
		expected: &models.User{
			Email:    "test111@gmail.com",
			Username: "test111",
			Password: "password",
		},
	}

	expectedQuery := "INSERT INTO users \\(email, username, password\\) VALUES \\(\\$1, \\$2, \\$3\\)"

	rows := sqlmock.NewRows([]string{"username", "email"}).AddRow(test.user.Username, test.user.Email)

	mock.ExpectExec(expectedQuery).
		WithArgs(test.user.Email, test.user.Username, test.user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT id, username, email, first_name, last_name from users WHERE username=\\$1").
		WithArgs(test.user.Username).
		WillReturnRows(rows)
	user, err := repo.CreateUser(&test.user)
	t.Logf("error: %v, user, %v", err, user)

	assert.Equal(t, test.expected.Email, user.Email)
	assert.Equal(t, test.expected.Username, user.Username)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserWhereUserIsAvaliable(t *testing.T) {
	user := &models.User{
		Email:    "test111@gmail.com",
		Username: "test111",
		Password: "password",
	}

	expectedQuery := "INSERT INTO users \\(email, username, password\\) VALUES \\(\\$1, \\$2, \\$3\\)"

	mock.ExpectExec(expectedQuery).
		WithArgs(user.Email, user.Username, user.Password).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(fmt.Errorf("some error"))

	user, err := repo.CreateUser(user)
	t.Logf("error: %v, user, %v", err, user)

	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
