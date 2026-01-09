package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"khalif-backend/internal/platform/config"
	"khalif-backend/internal/platform/logger"
	"khalif-backend/pkg/messages"

	"go.uber.org/zap"
)

// EmailService interface for sending emails
type EmailService interface {
	SendOTP(toEmail, toName, otpCode string) error
	SendPasswordReset(toEmail, toName, resetToken string) error
	SendWelcome(toEmail, username string) error
}

// BrevoEmailService implements EmailService using Brevo API
type BrevoEmailService struct {
	apiKey       string
	apiURL       string
	senderEmail  string
	senderName   string
	templatePath string
	templates    map[string]*template.Template
	templateMu   sync.RWMutex
}

// NewBrevoEmailService creates a new Brevo email service with template caching
func NewBrevoEmailService(cfg *config.Config) EmailService {
	service := &BrevoEmailService{
		apiKey:       cfg.BrevoAPIKey,
		apiURL:       cfg.BrevoAPIURL,
		senderEmail:  cfg.BrevoSenderEmail,
		senderName:   cfg.BrevoSenderName,
		templatePath: cfg.EmailTemplatePath,
		templates:    make(map[string]*template.Template),
	}

	// Pre-load templates at startup
	templateNames := []string{
		messages.EmailTemplateOTP,
		messages.EmailTemplatePasswordReset,
		messages.EmailTemplateWelcome,
	}

	for _, name := range templateNames {
		if err := service.cacheTemplate(name); err != nil {
			logger.Log.Warn("Failed to pre-cache email template", zap.String("template", name), zap.Error(err))
		}
	}

	logger.Log.Info("Email service initialized with template caching", zap.Int("cached_templates", len(service.templates)))
	return service
}

// BrevoEmailRequest represents the Brevo API request with custom HTML
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

// OTPData holds data for OTP email template
type OTPData struct {
	OTPCode string
}

// PasswordResetData holds data for password reset email template
type PasswordResetData struct {
	Username   string
	ResetToken string
}

// WelcomeData holds data for welcome email template
type WelcomeData struct {
	Username string
}

// cacheTemplate loads and caches a template
func (s *BrevoEmailService) cacheTemplate(templateName string) error {
	templateFile := filepath.Join(s.templatePath, templateName)

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	s.templateMu.Lock()
	s.templates[templateName] = tmpl
	s.templateMu.Unlock()

	return nil
}

// getTemplate gets a cached template or loads it if not cached
func (s *BrevoEmailService) getTemplate(templateName string) (*template.Template, error) {
	s.templateMu.RLock()
	tmpl, exists := s.templates[templateName]
	s.templateMu.RUnlock()

	if exists {
		return tmpl, nil
	}

	// Template not cached, load it
	if err := s.cacheTemplate(templateName); err != nil {
		return nil, err
	}

	s.templateMu.RLock()
	tmpl = s.templates[templateName]
	s.templateMu.RUnlock()

	return tmpl, nil
}

// executeTemplate executes a cached template with data
func (s *BrevoEmailService) executeTemplate(templateName string, data interface{}) (string, error) {
	tmpl, err := s.getTemplate(templateName)
	if err != nil {
		logger.Log.Error("Failed to get template", zap.String("template", templateName), zap.Error(err))
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		logger.Log.Error("Failed to execute template", zap.String("template", templateName), zap.Error(err))
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// sendEmail sends an email via Brevo API
func (s *BrevoEmailService) sendEmail(toEmail, toName, subject, htmlContent string) error {
	reqBody := BrevoEmailRequest{
		Sender: BrevoContact{
			Email: s.senderEmail,
			Name:  s.senderName,
		},
		To: []BrevoContact{
			{Email: toEmail, Name: toName},
		},
		Subject:     subject,
		HTMLContent: htmlContent,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.Log.Error("Failed to marshal email request", zap.Error(err))
		return fmt.Errorf("failed to prepare email: %w", err)
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonBody))
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

	return nil
}

// SendOTP sends OTP verification email via Brevo
func (s *BrevoEmailService) SendOTP(toEmail, toName, otpCode string) error {
	data := OTPData{OTPCode: otpCode}

	htmlContent, err := s.executeTemplate(messages.EmailTemplateOTP, data)
	if err != nil {
		return err
	}

	if err := s.sendEmail(toEmail, toName, messages.EmailSubjectOTP, htmlContent); err != nil {
		return err
	}

	logger.Log.Info("OTP email sent successfully", zap.String("to", toEmail))
	return nil
}

// SendPasswordReset sends password reset email via Brevo
func (s *BrevoEmailService) SendPasswordReset(toEmail, toName, resetToken string) error {
	data := PasswordResetData{
		Username:   toName,
		ResetToken: resetToken,
	}

	htmlContent, err := s.executeTemplate(messages.EmailTemplatePasswordReset, data)
	if err != nil {
		return err
	}

	if err := s.sendEmail(toEmail, toName, messages.EmailSubjectPasswordReset, htmlContent); err != nil {
		return err
	}

	logger.Log.Info("Password reset email sent successfully", zap.String("to", toEmail))
	return nil
}

// SendWelcome sends welcome email after account activation
func (s *BrevoEmailService) SendWelcome(toEmail, username string) error {
	data := WelcomeData{Username: username}

	htmlContent, err := s.executeTemplate(messages.EmailTemplateWelcome, data)
	if err != nil {
		return err
	}

	if err := s.sendEmail(toEmail, username, messages.EmailSubjectWelcome, htmlContent); err != nil {
		return err
	}

	logger.Log.Info("Welcome email sent successfully", zap.String("to", toEmail))
	return nil
}
