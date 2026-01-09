package audio

import (
	"khalif-backend/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type MockAudioRepo struct {
	mock.Mock
}

func (m *MockAudioRepo) Create(audio *domain.Audio) error {
	args := m.Called(audio)
	return args.Error(0)
}

func (m *MockAudioRepo) FindByID(id uint) (*domain.Audio, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Audio), args.Error(1)
}

func (m *MockAudioRepo) FindByAudioID(audioID string) (*domain.Audio, error) {
	args := m.Called(audioID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Audio), args.Error(1)
}

func (m *MockAudioRepo) FindAll(page, limit int) ([]domain.Audio, int64, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Audio), args.Get(1).(int64), args.Error(2)
}

func (m *MockAudioRepo) Update(audio *domain.Audio) error {
	args := m.Called(audio)
	return args.Error(0)
}

func (m *MockAudioRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAudioRepo) IncrementListeningCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
