package interfaces

import "song/internal/domain"

// SongRepo представляет интерфейс для работы с песнями
type SongRepo interface {
	// GetLib получает библиотеку песен с пагинацией
	GetLib(song domain.Song, page domain.Page) (*[]domain.Song, error)

	// GetText получает текст песни по идентификатору
	GetText(id domain.Id) (*domain.SongText, error)

	// DelSong удаляет песню по идентификатору
	DelSong(id domain.Id) error

	// ChangeSong изменяет данные песни
	ChangeSong(song domain.Song) error

	// CreateSong создает новую песню
	CreateSong(song domain.Song) (*domain.Id, error)
}
