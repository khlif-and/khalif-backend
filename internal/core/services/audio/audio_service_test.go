package audio

import (
	"io"
	"khalif-backend/internal/core/domain"
	"khalif-backend/pkg/messages"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) UploadFile(file io.Reader, filename string, directory string) (string, string, error) {
	args := m.Called(file, filename, directory)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockStorageProvider) DeleteFile(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func TestAudioService_Create(t *testing.T) {
	req := &domain.CreateAudioRequest{
		Title:               "Test Audio",
		AudioFile:           "/uploads/audio/test.mp3",
		ThumbnailFile:       "/uploads/thumbnails/test.jpg",
		ColorThumbnailAudio: "#FF5733",
		NameUstadz:          "Ustadz Test",
		DurationAudio:       300,
	}

	t.Run("Success_WithDirectURL", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("Create", mock.AnythingOfType("*domain.Audio")).Return(nil).Once()

		// Case where no file stream is provided, only URL in request
		audio, err := service.Create(req, nil, "", nil, "")

		assert.NoError(t, err)
		assert.NotNil(t, audio)
		assert.Equal(t, req.Title, audio.Title)
		assert.NotEmpty(t, audio.AudioID)
	})

	t.Run("Fail_TitleRequired", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		emptyReq := &domain.CreateAudioRequest{
			AudioFile:  "/uploads/audio/test.mp3",
			NameUstadz: "Ustadz Test",
		}

		audio, err := service.Create(emptyReq, nil, "", nil, "")

		assert.Error(t, err)
		assert.Nil(t, audio)
		assert.Equal(t, messages.ErrAudioTitleRequired, err.Error())
	})

	t.Run("Fail_AudioFileRequired", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		emptyReq := &domain.CreateAudioRequest{
			Title:      "Test Audio",
			NameUstadz: "Ustadz Test",
		}

		// Both Request URL and Stream are empty
		audio, err := service.Create(emptyReq, nil, "", nil, "")

		assert.Error(t, err)
		assert.Nil(t, audio)
		assert.Equal(t, messages.ErrAudioFileRequired, err.Error())
	})
	
	t.Run("Fail_UstadzNameRequired", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		emptyReq := &domain.CreateAudioRequest{
			Title:     "Test Audio",
			AudioFile: "/uploads/audio/test.mp3",
		}

		audio, err := service.Create(emptyReq, nil, "", nil, "")

		assert.Error(t, err)
		assert.Nil(t, audio)
		assert.Equal(t, messages.ErrUstadzNameRequired, err.Error())
	})
}

func TestAudioService_GetByID(t *testing.T) {
	mockAudio := &domain.Audio{
		ID:         1,
		AudioID:    "test-uuid",
		Title:      "Test Audio",
		NameUstadz: "Ustadz Test",
	}

	t.Run("Success", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("FindByID", uint(1)).Return(mockAudio, nil).Once()

		audio, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, audio)
		assert.Equal(t, mockAudio.Title, audio.Title)
	})

	t.Run("Fail_NotFound", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("FindByID", uint(999)).Return(nil, nil).Once()

		audio, err := service.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, audio)
		assert.Equal(t, messages.ErrAudioNotFound, err.Error())
	})
}

func TestAudioService_GetAll(t *testing.T) {
	mockAudios := []domain.Audio{
		{ID: 1, Title: "Audio 1"},
		{ID: 2, Title: "Audio 2"},
	}

	t.Run("Success", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("FindAll", 1, 10).Return(mockAudios, int64(2), nil).Once()

		response, err := service.GetAll(1, 10)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Audios, 2)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 1, response.TotalPages)
	})
}

func TestAudioService_Update(t *testing.T) {
	mockAudio := &domain.Audio{
		ID:         1,
		AudioID:    "test-uuid",
		Title:      "Old Title",
		NameUstadz: "Ustadz Test",
	}

	t.Run("Success", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("FindByID", uint(1)).Return(mockAudio, nil).Once()
		mockAudioRepo.On("Update", mock.AnythingOfType("*domain.Audio")).Return(nil).Once()

		req := &domain.UpdateAudioRequest{
			Title: "New Title",
		}

		audio, err := service.Update(1, req, nil, "", nil, "")

		assert.NoError(t, err)
		assert.NotNil(t, audio)
		assert.Equal(t, "New Title", audio.Title)
	})

	t.Run("Fail_NotFound", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("FindByID", uint(999)).Return(nil, nil).Once()

		req := &domain.UpdateAudioRequest{
			Title: "New Title",
		}

		audio, err := service.Update(999, req, nil, "", nil, "")

		assert.Error(t, err)
		assert.Nil(t, audio)
		assert.Equal(t, messages.ErrAudioNotFound, err.Error())
	})
}

func TestAudioService_Delete(t *testing.T) {
	mockAudio := &domain.Audio{
		ID:      1,
		AudioID: "test-uuid",
	}

	t.Run("Success", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("FindByID", uint(1)).Return(mockAudio, nil).Once()
		mockAudioRepo.On("Delete", uint(1)).Return(nil).Once()

		err := service.Delete(1)

		assert.NoError(t, err)
	})
}

func TestAudioService_IncrementListeningCount(t *testing.T) {
	mockAudio := &domain.Audio{
		ID:      1,
		AudioID: "test-uuid",
	}

	t.Run("Success", func(t *testing.T) {
		mockAudioRepo := new(MockAudioRepo)
		mockStorage := new(MockStorageProvider)
		service := NewAudioService(mockAudioRepo, mockStorage)

		mockAudioRepo.On("FindByID", uint(1)).Return(mockAudio, nil).Once()
		mockAudioRepo.On("IncrementListeningCount", uint(1)).Return(nil).Once()

		err := service.IncrementListeningCount(1)

		assert.NoError(t, err)
	})
}
