package admin

import (
	"errors"
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"
	"strings"
)

type adminService struct {
	repo ports.AdminRepository
}

func NewAdminService(repo ports.AdminRepository) ports.AdminService {
	return &adminService{repo: repo}
}

func (s *adminService) UpdateProfile(adminID uint, req *domain.UpdateProfileRequest) (*domain.Admin, error) {
	admin, err := s.repo.FindByID(adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("admin not found")
	}

	// Update fields if provided
	if req.Username != "" {
		// Check uniqueness if changed
		if req.Username != admin.Username {
			existing, _ := s.repo.FindByUsername(req.Username)
			if existing != nil {
				return nil, errors.New(messages.ErrUsernameExists)
			}
			admin.Username = req.Username
		}
	}

	if req.Email != "" {
		if req.Email != admin.Email {
			existing, _ := s.repo.FindByEmail(req.Email)
			if existing != nil {
				return nil, errors.New(messages.ErrEmailAlreadyExists)
			}
			admin.Email = req.Email
		}
	}

	if req.Phone != "" {
		admin.Phone = req.Phone
	}

	var oldProfilePicture string
	if req.ProfilePicture != "" && req.ProfilePicture != admin.ProfilePicture {
		oldProfilePicture = admin.ProfilePicture
		admin.ProfilePicture = req.ProfilePicture
	}

	if err := s.repo.Update(admin); err != nil {
		return nil, err
	}

	// Cleanup old profile picture if updated
	if oldProfilePicture != "" {
		// Convert URL to file path (e.g., /uploads/... -> uploads/...)
		filePath := strings.TrimPrefix(oldProfilePicture, "/")
		// Ignore error during cleanup (don't fail the request)
		_ = utils.DeleteUploadedFile(filePath)
	}

	return admin, nil
}
