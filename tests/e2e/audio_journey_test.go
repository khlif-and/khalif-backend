package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createAudioMultipartRequest(fields map[string]string, files map[string]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	for key, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, "", err
		}
		defer file.Close()

		// Determine Content-Type based on extension
		contentType := "application/octet-stream"
		ext := strings.ToLower(filepath.Ext(filePath))
		switch ext {
		case ".mp3":
			contentType = "audio/mpeg"
		case ".wav":
			contentType = "audio/wav"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		}

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, key, filepath.Base(filePath)))
		h.Set("Content-Type", contentType)

		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, "", err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, "", err
		}
	}

	writer.Close()
	return body, writer.FormDataContentType(), nil
}

func createDummyFile(t *testing.T, name, content string) string {
	path := filepath.Join(t.TempDir(), name)
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
	return path
}

func TestE2E_AudioJourney(t *testing.T) {
	server := setupE2EServer(t)
	defer server.Close()

	client := &http.Client{}
	var adminAccessToken string
	var audioID uint // Will be set after creation

	// Step 1: Register Admin
	t.Run("1_RegisterAdmin", func(t *testing.T) {
		fields := map[string]string{
			"username": "adminaudio",
			"email":    "adminaudio@example.com",
			"phone":    "08999999999",
			"password": "password123",
		}
		body, contentType := createMultipartRequest(fields)
		resp, err := http.Post(server.URL+"/api/v1/auth/register", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Step 2: Login Admin
	t.Run("2_LoginAdmin", func(t *testing.T) {
		fields := map[string]string{
			"email":    "adminaudio@example.com",
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
		adminAccessToken = data["token"].(string)
	})

	// Step 3: Create Audio
	t.Run("3_CreateAudio", func(t *testing.T) {
		// Create dummy files
		audioPath := createDummyFile(t, "test_audio.mp3", "dummy audio content")
		thumbnailPath := createDummyFile(t, "test_thumb.jpg", "dummy image content")

		fields := map[string]string{
			"title":                 "Kajian E2E Test",
			"name_ustadz":           "Ustadz E2E",
			"duration_audio":        "300",
			"color_thumbnail_audio": "#FF0000",
		}
		files := map[string]string{
			"audio_file":     audioPath,
			"thumbnail_file": thumbnailPath,
		}

		body, contentType, err := createAudioMultipartRequest(fields, files)
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", server.URL+"/api/v1/admin/audio", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+adminAccessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		
		assert.Equal(t, "Kajian E2E Test", data["title"])
		assert.Equal(t, "Ustadz E2E", data["name_ustadz"])
		
		// Capture ID for next steps
		idFloat := data["id"].(float64)
		audioID = uint(idFloat)
	})

	// Step 4: Get All Audios (Public)
	t.Run("4_GetAllAudios", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/audio")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		audios := data["audios"].([]interface{})
		
		assert.NotEmpty(t, audios)
	})

	// Step 5: Get Audio By ID (Public)
	t.Run("5_GetAudioByID", func(t *testing.T) {
		// Use ID captured from Create step
		url := server.URL + "/api/v1/audio/" + convertUintToString(audioID) // Need helper or just fmt.Sprintf
		resp, err := http.Get(url)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		
		assert.Equal(t, "Kajian E2E Test", data["title"])
	})

	// Step 6: Increment Listening Count
	t.Run("6_IncrementListeningCount", func(t *testing.T) {
		url := server.URL + "/api/v1/audio/" + convertUintToString(audioID) + "/listen"
		req, _ := http.NewRequest("POST", url, nil) // Empty body
		
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify count increased by getting details again
		urlGet := server.URL + "/api/v1/audio/" + convertUintToString(audioID)
		respGet, err := http.Get(urlGet)
		assert.NoError(t, err)
		defer respGet.Body.Close()
		
		var result map[string]interface{}
		json.NewDecoder(respGet.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		
		// In Create, initial is 0. Now should be 1.
		listeningCount := data["listening_count"].(float64)
		assert.Equal(t, float64(1), listeningCount)
	})

	// Step 7: Update Audio (Admin)
	t.Run("7_UpdateAudio", func(t *testing.T) {
		fields := map[string]string{
			"title": "Kajian E2E Updated",
		}
		// Try updating just title without files
		body, contentType, err := createAudioMultipartRequest(fields, nil)
		assert.NoError(t, err)

		url := server.URL + "/api/v1/admin/audio/" + convertUintToString(audioID)
		req, _ := http.NewRequest("PUT", url, body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+adminAccessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		
		assert.Equal(t, "Kajian E2E Updated", data["title"])
	})

	// Step 8: Delete Audio (Admin)
	t.Run("8_DeleteAudio", func(t *testing.T) {
		url := server.URL + "/api/v1/admin/audio/" + convertUintToString(audioID)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminAccessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Step 9: Verify Deletion
	t.Run("9_VerifyDeletion", func(t *testing.T) {
		url := server.URL + "/api/v1/audio/" + convertUintToString(audioID)
		resp, err := http.Get(url)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// Helper to convert uint to string manually to avoid heavy fmt import if possible, but fmt is fine
func convertUintToString(id uint) string {
	return fmt.Sprintf("%d", id)
}
