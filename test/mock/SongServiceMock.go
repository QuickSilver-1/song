package mock

import (
	"net/url"
	"song/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockSongService - mock для SongService
type MockSongService struct {
	mock.Mock
}

func (m *MockSongService) GetLib(filter domain.Song, page domain.Page) (*[]domain.Song, error) {
	args := m.Called(filter, page)
	return args.Get(0).(*[]domain.Song), args.Error(1)
}

func (m *MockSongService) GetText(id uint64, page domain.Page) (*domain.SongText, error) {
	args := m.Called(id, page)
	return args.Get(0).(*domain.SongText), args.Error(1)
}

func (m *MockSongService) DelSong(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSongService) ChangeSong(song domain.Song) error {
	args := m.Called(song)
	return args.Error(0)
}

func (m *MockSongService) CreateSong(data domain.SongDataByUser, apiUrl *url.URL) (*domain.Id, error) {
	args := m.Called(data, apiUrl)
	return args.Get(0).(*domain.Id), args.Error(1)
}
