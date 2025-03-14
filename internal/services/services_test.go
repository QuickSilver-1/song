package services

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"song/internal/domain"
	"song/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockSongRepo  *mock.MockSongRepo
	mockCacheRepo *mock.MockCacheRepo
	service       *SongService
)

func init() {
	mockSongRepo = new(mock.MockSongRepo)
	mockCacheRepo = new(mock.MockCacheRepo)
	service = NewSongService(mockCacheRepo, mockSongRepo)
}

// Тест для метода GetLib
func TestSongService_GetLib(t *testing.T) {
	filter := domain.Song{}
	page := domain.Page(1)
	expectedSongs := &[]domain.Song{{ID: 1, Name: "Test Song"}}

	mockSongRepo.On("GetLib", filter, page).Return(expectedSongs, nil)

	result, err := service.GetLib(filter, page)

	assert.Nil(t, err)
	assert.Equal(t, expectedSongs, result)
	mockSongRepo.AssertExpectations(t)
}

// Тест для метода GetText
func TestSongService_GetText(t *testing.T) {
	id := domain.Id(1)
	page := domain.Page(1)
	text := domain.SongText("Verse 1\\n\\nVerse 2")
	expectedText := "Verse 1"

	mockCacheRepo.On("GetText", id).Return(&text, nil)

	result, err := service.GetText(uint64(id), page)

	assert.Nil(t, err)
	assert.Equal(t, expectedText, *result)
	mockCacheRepo.AssertExpectations(t)
}

// Тест для метода DelSong
func TestSongService_DelSong(t *testing.T) {
	id := domain.Id(1)

	mockSongRepo.On("DelSong", id).Return(nil)
	mockCacheRepo.On("DelKey", id).Return(nil)

	err := service.DelSong(uint64(id))

	assert.Nil(t, err)
	mockSongRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

// Тест для метода ChangeSong
func TestSongService_ChangeSong(t *testing.T) {
	song := domain.Song{ID: 1, Name: "Updated Song"}

	mockSongRepo.On("ChangeSong", song).Return(nil)

	err := service.ChangeSong(song)

	assert.Nil(t, err)
	mockSongRepo.AssertExpectations(t)
}

// Тест для метода CreateSong - различные ошибки
func TestSongService_CreateSong_Errors(t *testing.T) {
	data := domain.SongDataByUser{
		Group: "Test Group",
		Name:  "Test Song",
	}
	apiUrl, _ := url.Parse("http://example.com")

	tests := []struct {
		name           string
		apiUrl         *url.URL
		mockHTTPServer func() *httptest.Server
		expectedError  *domain.RequestError
	}{
		{
			name:   "Ошибка создания запроса",
			apiUrl: nil,
			expectedError: &domain.RequestError{
				Err:  "error creating request - test error",
				Code: http.StatusInternalServerError,
			},
		},
		{
			name: "Ошибка выполнения запроса",
			mockHTTPServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			expectedError: &domain.RequestError{
				Err:  "error making request - test error",
				Code: http.StatusInternalServerError,
			},
		},
		{
			name: "Ошибка декодирования ответа",
			mockHTTPServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("invalid json"))
				}))
			},
			expectedError: &domain.RequestError{
				Err:  "error decoding response - test error",
				Code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockHTTPServer != nil {
				ts := tt.mockHTTPServer()
				defer ts.Close()
				apiUrl, _ = url.Parse(ts.URL)
			}

			_, err := service.CreateSong(data, apiUrl)

			assert.NotNil(t, err)
			assert.Equal(t, tt.expectedError.Code, err.(*domain.RequestError).Code)
		})
	}
}
