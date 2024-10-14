package repository_test

import (
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
	sqlxDB *sqlx.DB
	mock   sqlmock.Sqlmock
	repo   repository.UserRepository
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

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}
func TestGetUserByEmailOrUsername(t *testing.T) {

	id, _ := uuid.NewUUID()
	expectedUser := &models.User{
		ID:       id,
		Username: "testuser",
		Email:    "test@example.com",
	}
	tests := []struct {
		username string
		expected *models.User
	}{
		{
			username: "testuser",
			expected: expectedUser,
		},
		{
			username: "testuser1",
			expected: nil,
		},
	}

	// Define the expected query and result
	expectedQuery := "SELECT \\* from users WHERE email=\\$1 OR username=\\$1"

	// Set up the expectation
	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email)

	for _, test := range tests {
		// Call the method being tested

		mock.ExpectQuery(expectedQuery).
			WithArgs(test.username).
			WillReturnRows(rows)
		user, err := repo.GetUserByEmailOrUsername(test.username)
		// Assert the results
		assert.NoError(t, err)
		assert.Equal(t, test.expected, user)
		if user != nil {
			assert.Equal(t, expectedUser.Username, user.Username)
			assert.Equal(t, expectedUser.Email, user.Email)
		}
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// func TestGetUserByID(t *testing.T) {

// 	// Define the expected query and result
// 	expectedQuery := "SELECT \\* from users WHERE email=\\$1 OR username=\\$1"
// 	id, _ := uuid.NewUUID()
// 	expectedUser := &models.User{
// 		ID:       id,
// 		Username: "testuser",
// 		Email:    "test@example.com",
// 	}

// 	// Set up the expectation
// 	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
// 		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email)

// 	mock.ExpectQuery(expectedQuery).
// 		WithArgs("testuser").
// 		WillReturnRows(rows)

// 	// Call the method being tested
// 	user, err := repo.GetUserByEmailOrUsername("testuser")

// 	// Assert the results
// 	assert.NoError(t, err)
// 	assert.NotNil(t, user)
// 	if user != nil {
// 		assert.Equal(t, expectedUser.Username, user.Username)
// 		assert.Equal(t, expectedUser.Email, user.Email)
// 	}

// 	// Ensure all expectations were met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
