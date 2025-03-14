package mock

import (
	"song/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockSongRepo - mock для интерфейса SongRepo
type MockSongRepo struct {
	mock.Mock
}

func (m *MockSongRepo) GetLib(filter domain.Song, page domain.Page) (*[]domain.Song, error) {
	args := m.Called(filter, page)
	return args.Get(0).(*[]domain.Song), args.Error(1)
}

func (m *MockSongRepo) GetText(id domain.Id) (*domain.SongText, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.SongText), args.Error(1)
}

func (m *MockSongRepo) DelSong(id domain.Id) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSongRepo) ChangeSong(song domain.Song) error {
	args := m.Called(song)
	return args.Error(0)
}

func (m *MockSongRepo) CreateSong(song domain.Song) (*domain.Id, error) {
	args := m.Called(song)
	return args.Get(0).(*domain.Id), args.Error(1)
}
