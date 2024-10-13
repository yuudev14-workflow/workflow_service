package auth_api_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/yuudev14-workflow/workflow-service/api"
// 	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
// )

// func TestMain(m *testing.M) {
// 	test_helpers.SetupTestEnvironment("../../.env.test")
// 	exitCode := m.Run()

// 	os.Exit(exitCode)
// }

// func TestSignupRoute(t *testing.T) {
// 	router := api.InitRouter()

// 	w := httptest.NewRecorder()

// 	users := map[string]string{
// 		"username": "john_auth_doe",
// 		"email":    "john_auth@example.com",
// 		"password": "password",
// 	}

// 	// Marshal the slice of users to JSON
// 	jsonData, err := json.Marshal(users)
// 	if err != nil {
// 		panic(err)
// 	}
// 	req, _ := http.NewRequest("POST", "/api/auth/v1/sign-up", bytes.NewBuffer(jsonData))
// 	router.ServeHTTP(w, req)
// 	var responseBody map[string]interface{}
// 	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
// 	if err != nil {
// 		t.Errorf("error in unmarshaling response body... %v", err)
// 		return
// 	}
// 	logging.Logger.Info(responseBody)

// 	assert.Equal(t, 200, w.Code)
// }

// func ptrStr(s string) *string {
// 	return &s
// }

// func TestLoginRoute(t *testing.T) {
// 	tests := []struct {
// 		username     string
// 		password     string
// 		expectedCode int
// 		expectedMsg  *string
// 	}{
// 		{
// 			username:     "john_auth_doe",
// 			password:     "password",
// 			expectedCode: 200,
// 			expectedMsg:  nil,
// 		},
// 		{
// 			username:     "john_auth_doe",
// 			password:     "passwor",
// 			expectedCode: 400,
// 			expectedMsg:  ptrStr("password is not correct"),
// 		},
// 		{
// 			username:     "john_auth_e",
// 			password:     "password",
// 			expectedCode: 400,
// 			expectedMsg:  ptrStr("user is not found"),
// 		},
// 	}
// 	router := api.InitRouter()

// 	for _, e := range tests {
// 		w := httptest.NewRecorder()
// 		users := map[string]string{
// 			"username": e.username,
// 			"password": e.password,
// 		}

// 		logging.Logger.Debug(users, e)

// 		// Marshal the slice of users to JSON
// 		jsonData, err := json.Marshal(users)
// 		if err != nil {
// 			panic(err)
// 		}
// 		req, _ := http.NewRequest("POST", "/api/auth/v1/login", bytes.NewBuffer(jsonData))
// 		router.ServeHTTP(w, req)
// 		var responseBody map[string]interface{}
// 		err = json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		if err != nil {
// 			t.Errorf("error in unmarshaling response body... %v", err)
// 			return
// 		}
// 		logging.Logger.Info(responseBody)

// 		assert.Equal(t, e.expectedCode, w.Code)

// 		value, ok := responseBody["error"]
// 		if ok {
// 			assert.Equal(t, *e.expectedMsg, value)
// 		}
// 	}

// }
