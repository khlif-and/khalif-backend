package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuthHandler_Me_Integration(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	registerFields := map[string]string{
		"username": "meuser",
		"email":    "me@example.com",
		"phone":    "08444444444",
		"password": "password123",
	}
	body, contentType := createMultipartRequest(registerFields, nil, "", "")
	registerResp, _ := http.Post(server.URL+"/api/v1/users/auth/register", contentType, body)

	var registerResult map[string]interface{}
	json.NewDecoder(registerResp.Body).Decode(&registerResult)
	registerResp.Body.Close()

	data := registerResult["data"].(map[string]interface{})
	accessToken := data["token"].(string)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/v1/users/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		userData := result["data"].(map[string]interface{})
		assert.Equal(t, "meuser", userData["username"])
		assert.Equal(t, "me@example.com", userData["email"])
	})

	t.Run("Fail_NoToken", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/v1/auth/me", nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
