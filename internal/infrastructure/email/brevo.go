package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"khalif-backend/internal/platform/config"
	"khalif-backend/internal/platform/logger"

	"go.uber.org/zap"
)

// EmailService interface for sending emails
type EmailService interface {
	SendOTP(toEmail, toName, otpCode string) error
	SendPasswordReset(toEmail, toName, resetToken string) error
}

// BrevoEmailService implements EmailService using Brevo API
type BrevoEmailService struct {
	apiKey      string
	senderEmail string
	senderName  string
}

// NewBrevoEmailService creates a new Brevo email service
func NewBrevoEmailService(cfg *config.Config) EmailService {
	return &BrevoEmailService{
		apiKey:      cfg.BrevoAPIKey,
		senderEmail: cfg.BrevoSenderEmail,
		senderName:  cfg.BrevoSenderName,
	}
}

// BrevoEmailRequest represents the Brevo API request structure
type BrevoEmailRequest struct {
	Sender      BrevoContact   `json:"sender"`
	To          []BrevoContact `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent"`
}

// BrevoContact represents sender/recipient in Brevo
type BrevoContact struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// SendOTP sends OTP verification email via Brevo
func (s *BrevoEmailService) SendOTP(toEmail, toName, otpCode string) error {
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px;">
    <div style="max-width: 500px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px;">Khalif App</h1>
            <p style="color: rgba(255,255,255,0.9); margin: 10px 0 0;">Email Verification</p>
        </div>
        <div style="padding: 40px 30px; text-align: center;">
            <p style="color: #333; font-size: 16px; margin-bottom: 30px;">
                Halo <strong>%s</strong>,<br><br>
                Gunakan kode OTP berikut untuk verifikasi akun Anda:
            </p>
            <div style="background-color: #f8f9fa; border: 2px dashed #667eea; border-radius: 10px; padding: 20px; margin: 20px 0;">
                <span style="font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #667eea;">%s</span>
            </div>
            <p style="color: #666; font-size: 14px; margin-top: 30px;">
                Kode ini berlaku selama <strong>10 menit</strong>.<br>
                Jangan bagikan kode ini kepada siapapun.
            </p>
        </div>
        <div style="background-color: #f8f9fa; padding: 20px; text-align: center; border-top: 1px solid #eee;">
            <p style="color: #999; font-size: 12px; margin: 0;">
                © 2026 Khalif App. All rights reserved.
            </p>
        </div>
    </div>
</body>
</html>
`, toName, otpCode)

	reqBody := BrevoEmailRequest{
		Sender: BrevoContact{
			Email: s.senderEmail,
			Name:  s.senderName,
		},
		To: []BrevoContact{
			{Email: toEmail, Name: toName},
		},
		Subject:     fmt.Sprintf("Kode Verifikasi Khalif App: %s", otpCode),
		HTMLContent: htmlContent,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.Log.Error("Failed to marshal email request", zap.Error(err))
		return fmt.Errorf("failed to prepare email: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log.Error("Failed to create HTTP request", zap.Error(err))
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Failed to send email via Brevo", zap.Error(err))
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logger.Log.Error("Brevo API returned error", zap.Int("status_code", resp.StatusCode))
		return fmt.Errorf("email service returned error: %d", resp.StatusCode)
	}

	logger.Log.Info("OTP email sent successfully", zap.String("to", toEmail))
	return nil
}

// SendPasswordReset sends password reset email via Brevo
func (s *BrevoEmailService) SendPasswordReset(toEmail, toName, resetToken string) error {
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px;">
    <div style="max-width: 500px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <div style="background: linear-gradient(135deg, #e74c3c 0%%, #c0392b 100%%); padding: 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px;">Khalif App</h1>
            <p style="color: rgba(255,255,255,0.9); margin: 10px 0 0;">Password Reset</p>
        </div>
        <div style="padding: 40px 30px; text-align: center;">
            <p style="color: #333; font-size: 16px; margin-bottom: 30px;">
                Halo <strong>%s</strong>,<br><br>
                Anda menerima email ini karena ada permintaan reset password untuk akun Anda.
            </p>
            <p style="color: #333; font-size: 16px; margin-bottom: 20px;">
                Gunakan token berikut untuk reset password:
            </p>
            <div style="background-color: #f8f9fa; border: 2px dashed #e74c3c; border-radius: 10px; padding: 20px; margin: 20px 0; word-break: break-all;">
                <span style="font-size: 14px; font-weight: bold; color: #e74c3c;">%s</span>
            </div>
            <p style="color: #666; font-size: 14px; margin-top: 30px;">
                Token ini berlaku selama <strong>30 menit</strong>.<br>
                Jika Anda tidak meminta reset password, abaikan email ini.
            </p>
        </div>
        <div style="background-color: #f8f9fa; padding: 20px; text-align: center; border-top: 1px solid #eee;">
            <p style="color: #999; font-size: 12px; margin: 0;">
                © 2026 Khalif App. All rights reserved.
            </p>
        </div>
    </div>
</body>
</html>
`, toName, resetToken)

	reqBody := BrevoEmailRequest{
		Sender: BrevoContact{
			Email: s.senderEmail,
			Name:  s.senderName,
		},
		To: []BrevoContact{
			{Email: toEmail, Name: toName},
		},
		Subject:     "Reset Password - Khalif App",
		HTMLContent: htmlContent,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.Log.Error("Failed to marshal email request", zap.Error(err))
		return fmt.Errorf("failed to prepare email: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log.Error("Failed to create HTTP request", zap.Error(err))
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Failed to send email via Brevo", zap.Error(err))
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logger.Log.Error("Brevo API returned error", zap.Int("status_code", resp.StatusCode))
		return fmt.Errorf("email service returned error: %d", resp.StatusCode)
	}

	logger.Log.Info("Password reset email sent successfully", zap.String("to", toEmail))
	return nil
}
