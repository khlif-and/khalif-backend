package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login_Integration(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	registerFields := map[string]string{
		"username": "loginuser",
		"email":    "login@example.com",
		"phone":    "08222222222",
		"password": "password123",
	}
	body, contentType := createMultipartRequest(registerFields, nil, "", "")
	resp, err := http.Post(server.URL+"/api/v1/auth/register", contentType, body)
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusCreated {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		t.Logf("Registration failed: %v", result)
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Registration failed in login test setup")
	resp.Body.Close()

	t.Run("Success", func(t *testing.T) {
		fields := map[string]string{
			"email":    "login@example.com",
			"password": "password123",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			t.Logf("Login failed: %v", result)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		data := result["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.NotEmpty(t, data["refresh_token"])
	})

	t.Run("Fail_WrongPassword", func(t *testing.T) {
		fields := map[string]string{
			"email":    "login@example.com",
			"password": "wrongpassword",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail_UserNotFound", func(t *testing.T) {
		fields := map[string]string{
			"email":    "notfound@example.com",
			"password": "password123",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
