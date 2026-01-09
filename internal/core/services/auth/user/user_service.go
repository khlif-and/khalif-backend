package user

import (
	"errors"
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"
	"strings"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) UpdateProfile(userID uint, req *domain.UpdateProfileRequest) (*domain.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Username != "" {
		// Check uniqueness if changed
		if req.Username != user.Username {
			existing, _ := s.repo.FindByUsername(req.Username)
			if existing != nil {
				return nil, errors.New(messages.ErrUsernameExists)
			}
			user.Username = req.Username
		}
	}

	if req.Email != "" {
		if req.Email != user.Email {
			existing, _ := s.repo.FindByEmail(req.Email)
			if existing != nil {
				return nil, errors.New(messages.ErrEmailAlreadyExists)
			}
			user.Email = req.Email
		}
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	var oldProfilePicture string
	if req.ProfilePicture != "" && req.ProfilePicture != user.ProfilePicture {
		oldProfilePicture = user.ProfilePicture
		user.ProfilePicture = req.ProfilePicture
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	// Cleanup old profile picture if updated
	if oldProfilePicture != "" {
		// Convert URL to file path (e.g., /uploads/... -> uploads/...)
		filePath := strings.TrimPrefix(oldProfilePicture, "/")
		// Ignore error during cleanup (don't fail the request)
		_ = utils.DeleteUploadedFile(filePath)
	}

	return user, nil
}
