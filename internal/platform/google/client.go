package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"khalif-backend/internal/platform/logger"
	"khalif-backend/pkg/messages"

	"go.uber.org/zap"
)

// TokenInfo represents the response from Google's tokeninfo endpoint
type TokenInfo struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"` // Google ID
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// VerifyIDToken calls Google's tokeninfo endpoint to validate the ID token
func VerifyIDToken(idToken string, clientID string) (*TokenInfo, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	resp, err := httpClient.Get(url)
	if err != nil {
		logger.Log.Error("Failed to call Google API", zap.Error(err))
		return nil, fmt.Errorf("failed to call google api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Warn("Google API returned non-200 status", zap.Int("status", resp.StatusCode))
		return nil, errors.New(messages.ErrInvalidGoogleToken)
	}

	var tokenInfo TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		logger.Log.Error("Failed to decode Google response", zap.Error(err))
		return nil, fmt.Errorf("failed to decode google response: %w", err)
	}

	// Basic validation
	if clientID != "" && tokenInfo.Aud != clientID {
		logger.Log.Warn("Google Token audience mismatch", zap.String("expected", clientID), zap.String("got", tokenInfo.Aud))
		return nil, fmt.Errorf("token audience mismatch: expected %s, got %s", clientID, tokenInfo.Aud)
	}

	if tokenInfo.EmailVerified != "true" {
		logger.Log.Warn("Google email not verified", zap.String("email", tokenInfo.Email))
		return nil, errors.New(messages.ErrGoogleEmailNotVerified)
	}

	return &tokenInfo, nil
}
