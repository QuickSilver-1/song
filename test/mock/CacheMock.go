package mock

import (
	"song/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockCacheRepo - mock для интерфейса CacheRepo
type MockCacheRepo struct {
	mock.Mock
}

func (m *MockCacheRepo) CreateKey(id domain.Id, text domain.SongText) error {
	args := m.Called(id, text)
	return args.Error(0)
}

func (m *MockCacheRepo) GetText(id domain.Id) (*domain.SongText, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.SongText), args.Error(1)
}

func (m *MockCacheRepo) DelKey(id domain.Id) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCacheRepo) Close() error {
	return nil
}
