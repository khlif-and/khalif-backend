package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register_Integration(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	t.Run("Success_WithoutProfilePic", func(t *testing.T) {
		fields := map[string]string{
			"username": "testuser1",
			"email":    "test1@example.com",
			"phone":    "08111111111",
			"password": "password123",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/auth/register", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.NotEmpty(t, data["refresh_token"])
	})

	t.Run("Fail_MissingFields", func(t *testing.T) {
		fields := map[string]string{
			"username": "testuser2",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/auth/register", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Fail_InvalidEmail", func(t *testing.T) {
		fields := map[string]string{
			"username": "testuser3",
			"email":    "invalid-email",
			"phone":    "08111111113",
			"password": "password123",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/auth/register", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Fail_DuplicateEmail", func(t *testing.T) {
		fields := map[string]string{
			"username": "testuser4",
			"email":    "test1@example.com",
			"phone":    "08111111114",
			"password": "password123",
		}

		body, contentType := createMultipartRequest(fields, nil, "", "")

		resp, err := http.Post(server.URL+"/api/v1/auth/register", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
