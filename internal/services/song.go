package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"song/internal/domain"
	"song/internal/interfaces"
	"song/internal/presentation/logger"
	"strings"
)

// SongService - сервис для работы с песнями
type SongService struct {
	cacheDb interfaces.CacheRepo
	song    interfaces.SongRepo
}

// NewSongService создает новый объект SongService
func NewSongService(cache interfaces.CacheRepo, song interfaces.SongRepo) *SongService {
	return &SongService{
		cacheDb: cache,
		song:    song,
	}
}

// GetLib получает библиотеку песен
// filter - фильтр для песен
// page - номер страницы для пагинации
func (s *SongService) GetLib(filter domain.Song, page domain.Page) (*[]domain.Song, error) {
	songs, err := s.song.GetLib(filter, page)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

// GetText получает текст песни по идентификатору и номеру страницы
// id - идентификатор песни
// page - номер страницы для текста песни
func (s *SongService) GetText(id uint64, page domain.Page) (*domain.SongText, error) {
	text, err := s.cacheDb.GetText(id)
	if err != nil {
		return nil, err
	}

	if text == nil {
		text, err = s.song.GetText(id)
		if err != nil {
			return nil, err
		}

		err = s.cacheDb.CreateKey(id, *text)
		if err != nil {
			return nil, err
		}
	}

	verses := strings.Split(*text, "\\n\\n")
	fmt.Println(verses)
	n := len(verses)
	if page > n {
		return nil, &domain.InputDataError{
			Err:  fmt.Sprintf("This song have only %d verses", n),
			Code: http.StatusBadRequest,
		}
	}

	return &verses[page-1], nil
}

// DelSong удаляет песню по идентификатору
// id - идентификатор песни
func (s *SongService) DelSong(id uint64) error {
	err := s.song.DelSong(id)
	if err != nil {
		return err
	}

	err = s.cacheDb.DelKey(id)
	if err != nil {
		return err
	}

	return nil
}

// ChangeSong изменяет данные песни
// song - объект песни с новыми данными
func (s *SongService) ChangeSong(song domain.Song) error {
	return s.song.ChangeSong(song)
}

// CreateSong создает новую песню
// data - данные о песне, предоставленные пользователем
// apiUrl - URL API для получения дополнительных данных о песне
func (s *SongService) CreateSong(data domain.SongDataByUser, apiUrl *url.URL) (*domain.Id, error) {
	apiUrl.Path += "/info"
	params := url.Values{}
	params.Add("group", data.Group)
	params.Add("song", data.Name)
	apiUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", apiUrl.String(), nil)
	if err != nil {
		return nil, &domain.RequestError{
			Err:  fmt.Sprintf("error creating request - %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &domain.RequestError{
			Err:  fmt.Sprintf("error making request - %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Error(fmt.Sprintf("Error closing response body: %v", err))
		}
	}()
	var apiData domain.SongDataByApi
	err = json.NewDecoder(resp.Body).Decode(&apiData)
	if err != nil {
		return nil, &domain.RequestError{
			Err:  fmt.Sprintf("error decoding response - %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	id, err := s.song.CreateSong(domain.Song{
		Name:  data.Name,
		Group: data.Group,
		Date:  apiData.Date,
		Text:  apiData.Text,
		Link:  apiData.Link,
	})
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *SongService) Close() error {
	return s.cacheDb.Close()
}
