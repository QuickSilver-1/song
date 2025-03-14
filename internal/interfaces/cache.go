package interfaces

import "song/internal/domain"

// CacheRepo представляет интерфейс для работы с кэшом
type CacheRepo interface {
	// CreateKey создает ключ в кэше
	CreateKey(domain.Id, domain.SongText) error

	// GetText получает текст песни из кэша по идентификатору
	GetText(domain.Id) (*domain.SongText, error)

	// DelKey удаляет ключ из кэша
	DelKey(domain.Id) error

	// Close закрывает подключение
	Close() error
}
