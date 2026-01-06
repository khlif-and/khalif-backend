package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuthHandler_Login_Integration(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// 1. Setup: Register a User first
	registerFields := map[string]string{
		"username": "userlogin",
		"email":    "userlogin@example.com",
		"phone":    "08999999999",
		"password": "password123",
	}
	body, contentType := createMultipartRequest(registerFields, nil, "", "")
	resp, err := http.Post(server.URL+"/api/v1/users/auth/register", contentType, body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Test Login Success
	t.Run("Success", func(t *testing.T) {
		fields := map[string]string{
			"email":    "userlogin@example.com",
			"password": "password123",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/users/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			t.Logf("User Login failed: %v", result)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		data := result["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.NotEmpty(t, data["refresh_token"])
	})

	// 3. Test Login Failure (Wrong Password)
	t.Run("Fail_WrongPassword", func(t *testing.T) {
		fields := map[string]string{
			"email":    "userlogin@example.com",
			"password": "wrongpassword",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/users/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// 4. Test Login Failure (User Not Found)
	t.Run("Fail_UserNotFound", func(t *testing.T) {
		fields := map[string]string{
			"email":    "notfounduser@example.com",
			"password": "password123",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/users/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
