package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_AdminJourney(t *testing.T) {
	server := setupE2EServer(t)
	defer server.Close()

	client := &http.Client{}
	var accessToken, refreshToken string

	// Step 1: Register Admin
	t.Run("1_Register", func(t *testing.T) {
		fields := map[string]string{
			"username": "adminjourney",
			"email":    "adminjourney@example.com",
			"phone":    "08777777777",
			"password": "password123",
		}
		body, contentType := createMultipartRequest(fields)
		resp, err := http.Post(server.URL+"/api/v1/auth/register", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Step 2: Login to get Tokens
	t.Run("2_Login", func(t *testing.T) {
		fields := map[string]string{
			"email":    "adminjourney@example.com",
			"password": "password123",
		}
		body, contentType := createMultipartRequest(fields)
		resp, err := http.Post(server.URL+"/api/v1/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		accessToken = data["token"].(string)
		refreshToken = data["refresh_token"].(string)
	})

	// Step 3: Get Profile (Me)
	t.Run("3_GetMe", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/v1/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		assert.Equal(t, "adminjourney", data["username"])
	})

	// Step 4: Update Profile
	t.Run("4_UpdateProfile", func(t *testing.T) {
		updateData := map[string]interface{}{
			"username": "admin_updated",
			"email": "adminjourney@example.com",
			"phone": "08777777777", 
		}
		jsonBody, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", server.URL+"/api/v1/admin/update", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json") 
		
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode) 
	})

	// Step 5: Refresh Token
	t.Run("5_RefreshToken", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": refreshToken}
		jsonBody, _ := json.Marshal(reqBody)
		
		resp, err := http.Post(server.URL+"/api/v1/auth/refresh", "application/json", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		accessToken = data["token"].(string) // Update Access Token
	})

	// Step 6: Logout
	t.Run("6_Logout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Step 7: Verify Token Revoked (Refresh should fail)
	t.Run("7_VerifyLogout", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": refreshToken}
		jsonBody, _ := json.Marshal(reqBody)
		
		resp, err := http.Post(server.URL+"/api/v1/auth/refresh", "application/json", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		// Should be 401 Unauthorized because refresh token was revoked
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
