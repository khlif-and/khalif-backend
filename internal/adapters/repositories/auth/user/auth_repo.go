package user

import (
	"context"
	"database/sql"
	"fmt"
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/internal/platform/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) ports.UserAuthRepository {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CheckLockStatus(email string) (bool, string, error) {
	var isLocked bool
	var remainingTime sql.NullString
	var message string

	sqlDB, err := r.db.DB()
	if err != nil {
		return false, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use QueryRow to fetch results from the stored function
	row := sqlDB.QueryRowContext(ctx, "SELECT * FROM sp_check_user_lock_status($1)", email)
	err = row.Scan(&isLocked, &remainingTime)

	if err != nil {
		logger.Log.Error("Failed to check user lock status", zap.Error(err))
		return false, "", err
	}

	if isLocked {
		message = fmt.Sprintf("Account locked. Try again in %s", remainingTime.String)
	}

	return isLocked, message, nil
}

func (r *AuthRepo) RecordLoginFailure(email string) error {
	var isLocked bool
	var remainingTime sql.NullString

	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use QueryRow scan as function returns a table
	row := sqlDB.QueryRowContext(ctx, "SELECT * FROM sp_handle_user_login_failure($1)", email)
	err = row.Scan(&isLocked, &remainingTime)

	if err != nil {
		logger.Log.Error("Failed to record login failure", zap.Error(err))
		return err
	}

	return nil
}

func (r *AuthRepo) StoreRefreshToken(userID uint, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error {
	query := `INSERT INTO user_refresh_tokens (user_id, token_hash, expires_at, user_agent, ip_address) VALUES ($1, $2, $3, $4, $5)`
	return r.db.Exec(query, userID, tokenHash, expiresAt, userAgent, ipAddress).Error
}

func (r *AuthRepo) ValidateRefreshToken(tokenHash string) (uint, error) {
	var userID uint
	var isRevoked bool
	var expiresAt time.Time

	query := `SELECT user_id, is_revoked, expires_at FROM user_refresh_tokens WHERE token_hash = $1`
	err := r.db.Raw(query, tokenHash).Row().Scan(&userID, &isRevoked, &expiresAt)
	if err != nil {
		return 0, err
	}

	if isRevoked {
		return 0, fmt.Errorf("token revoked")
	}

	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("token expired")
	}

	return userID, nil
}

func (r *AuthRepo) RevokeRefreshToken(tokenHash string) error {
	query := `UPDATE user_refresh_tokens SET is_revoked = TRUE WHERE token_hash = $1`
	return r.db.Exec(query, tokenHash).Error
}

func (r *AuthRepo) RevokeAllUserTokens(userID uint) error {
	// Call the 1-arg version of the SP for Users
	return r.db.Exec("CALL sp_revoke_user_tokens($1)", userID).Error
}

// StoreOTP stores a new OTP for user verification
func (r *AuthRepo) StoreOTP(userID uint, otpCode string, expiresAt time.Time) error {
	query := `INSERT INTO otp_tokens (user_id, otp_code, expires_at) VALUES ($1, $2, $3)`
	return r.db.Exec(query, userID, otpCode, expiresAt).Error
}

// ValidateOTP validates the OTP code for the given email
func (r *AuthRepo) ValidateOTP(email, otpCode string) (*domain.User, error) {
	var user domain.User
	var otp domain.OTPToken

	// Find user by email
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Find valid OTP
	query := `SELECT id, user_id, otp_code, expires_at, is_used FROM otp_tokens 
			  WHERE user_id = $1 AND otp_code = $2 AND is_used = false AND expires_at > NOW()
			  ORDER BY created_at DESC LIMIT 1`
	
	if err := r.db.Raw(query, user.ID, otpCode).Scan(&otp).Error; err != nil {
		return nil, fmt.Errorf("invalid OTP")
	}

	if otp.ID == 0 {
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	return &user, nil
}

// MarkOTPUsed marks all OTPs for the user as used
func (r *AuthRepo) MarkOTPUsed(userID uint) error {
	query := `UPDATE otp_tokens SET is_used = TRUE WHERE user_id = $1`
	return r.db.Exec(query, userID).Error
}

// InvalidateOldOTPs invalidates all previous OTPs for the user
func (r *AuthRepo) InvalidateOldOTPs(userID uint) error {
	query := `UPDATE otp_tokens SET is_used = TRUE WHERE user_id = $1 AND is_used = FALSE`
	return r.db.Exec(query, userID).Error
}

// StorePasswordResetToken stores a new password reset token
func (r *AuthRepo) StorePasswordResetToken(userID uint, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	return r.db.Exec(query, userID, tokenHash, expiresAt).Error
}

// ValidatePasswordResetToken validates the password reset token and returns user ID
func (r *AuthRepo) ValidatePasswordResetToken(tokenHash string) (uint, error) {
	var userID uint
	var isUsed bool
	var expiresAt time.Time

	query := `SELECT user_id, is_used, expires_at FROM password_reset_tokens WHERE token_hash = $1`
	err := r.db.Raw(query, tokenHash).Row().Scan(&userID, &isUsed, &expiresAt)
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	if isUsed {
		return 0, fmt.Errorf("token already used")
	}

	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("token expired")
	}

	return userID, nil
}

// MarkPasswordResetTokenUsed marks the password reset token as used
func (r *AuthRepo) MarkPasswordResetTokenUsed(tokenHash string) error {
	query := `UPDATE password_reset_tokens SET is_used = TRUE WHERE token_hash = $1`
	return r.db.Exec(query, tokenHash).Error
}

// InvalidateOldPasswordResetTokens invalidates all previous password reset tokens for the user
func (r *AuthRepo) InvalidateOldPasswordResetTokens(userID uint) error {
	query := `UPDATE password_reset_tokens SET is_used = TRUE WHERE user_id = $1 AND is_used = FALSE`
	return r.db.Exec(query, userID).Error
}
