package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuthHandler_RefreshToken_Integration(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	registerFields := map[string]string{
		"username": "refreshuser",
		"email":    "refresh@example.com",
		"phone":    "08333333333",
		"password": "password123",
	}
	body, contentType := createMultipartRequest(registerFields, nil, "", "")
	registerResp, _ := http.Post(server.URL+"/api/v1/users/auth/register", contentType, body)

	var registerResult map[string]interface{}
	json.NewDecoder(registerResp.Body).Decode(&registerResult)
	registerResp.Body.Close()

	data := registerResult["data"].(map[string]interface{})
	refreshToken := data["refresh_token"].(string)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": refreshToken}
		jsonBody, _ := json.Marshal(reqBody)

		resp, err := http.Post(server.URL+"/api/v1/users/auth/refresh", "application/json", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		newData := result["data"].(map[string]interface{})
		assert.NotEmpty(t, newData["token"])
		assert.NotEmpty(t, newData["refresh_token"])
		assert.NotEqual(t, refreshToken, newData["refresh_token"])
	})

	t.Run("Fail_RevokedToken", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": refreshToken}
		jsonBody, _ := json.Marshal(reqBody)

		resp, err := http.Post(server.URL+"/api/v1/users/auth/refresh", "application/json", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
