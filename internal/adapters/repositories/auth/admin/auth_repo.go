package admin

import (
	"context"
	"fmt"
	"khalif-backend/internal/core/ports"
	"khalif-backend/internal/platform/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) ports.AuthRepository {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CheckLockStatus(email string) (bool, string, error) {
	var isLocked bool
	var message string

	sqlDB, err := r.db.DB()
	if err != nil {
		return false, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use QueryRow to fetch results from the stored procedure
	// Postgres CALL with OUT params returns a row
	row := sqlDB.QueryRowContext(ctx, "CALL sp_check_lock_status($1, NULL, NULL)", email)
	err = row.Scan(&isLocked, &message)

	if err != nil {
		logger.Log.Error("Failed to check lock status", zap.Error(err))
		return false, "", err
	}

	return isLocked, message, nil
}

func (r *AuthRepo) RecordLoginFailure(email string) error {
	var isLocked bool
	var message string

	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use QueryRow to fetch results from the stored procedure
	// Postgres CALL with OUT params returns a row
	row := sqlDB.QueryRowContext(ctx, "CALL sp_handle_login_failure($1, NULL, NULL)", email)
	err = row.Scan(&isLocked, &message)

	if err != nil {
		logger.Log.Error("Failed to record login failure", zap.Error(err))
		return err
	}

	return nil
}

func (r *AuthRepo) StoreRefreshToken(userID uint, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error {
	query := `INSERT INTO refresh_tokens (admin_id, token_hash, expires_at, user_agent, ip_address) VALUES ($1, $2, $3, $4, $5)`
	return r.db.Exec(query, userID, tokenHash, expiresAt, userAgent, ipAddress).Error
}

func (r *AuthRepo) ValidateRefreshToken(tokenHash string) (uint, error) {
	var adminID uint
	var isRevoked bool
	var expiresAt time.Time

	query := `SELECT admin_id, is_revoked, expires_at FROM refresh_tokens WHERE token_hash = $1`
	err := r.db.Raw(query, tokenHash).Row().Scan(&adminID, &isRevoked, &expiresAt)
	if err != nil {
		return 0, err
	}

	if isRevoked {
		return 0, fmt.Errorf("token revoked")
	}

	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("token expired")
	}

	return adminID, nil
}

func (r *AuthRepo) RevokeRefreshToken(tokenHash string) error {
	query := `UPDATE refresh_tokens SET is_revoked = TRUE WHERE token_hash = $1`
	return r.db.Exec(query, tokenHash).Error
}

func (r *AuthRepo) RevokeAllTokens(adminID uint) error {
	// Call the 1-arg version of the SP for Admins
	return r.db.Exec("CALL sp_revoke_user_tokens($1)", adminID).Error
}
