package interfaces

// ServerRepo представляет интерфейс для управления сервером
type ServerRepo interface {
	// Start запускает сервер
	Start() error

	// Shutdown завершает работу сервера по паттерну Graceful Shutdown
	Shutdown() error
}
