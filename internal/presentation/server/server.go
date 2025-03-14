package server

import (
	"net/url"
	"song/internal/presentation/logger"
	"song/internal/presentation/postgres"
	"song/internal/services"

	"github.com/gin-gonic/gin"
)

// Переменные для доступа к сервисам
var (
	SongService *services.SongService
	ApiUrl      *url.URL
)

// Константы http ответов
const (
	STATUS_UNAUTHORIZED    = "Authorization required"
	STATUS_INTERNAL_SERVER = "Sorry, something went wrong, we are already solving the problem"
	STATUS_BAD_REQUEST     = "Invalid data"
)

// Server определяет сервер с сервисами
type Server struct {
	srv *gin.Engine
}

// NewServer создает новый экземпляр Server
func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	srv := gin.New()

	h := NewHandlers()

	srv.Use(LoggerMiddleware())
	srv.GET("/lib", h.GetLib)
	srv.GET("/text", h.GetText)
	srv.DELETE("/song", h.DelSong)
	srv.PATCH("/song", h.ChangeSong)
	srv.POST("/song", h.CreateSong)

	logger.Logger.Info("Server has been created")
	return &Server{
		srv: srv,
	}
}

// Start запускает сервер
func (s *Server) Start(songService *services.SongService, api *url.URL, port string) error {
	SongService = songService
	ApiUrl = api

	logger.Logger.Debug("Starting server")
	err := s.srv.Run(":" + port)

	if err != nil {
		return err
	}

	logger.Logger.Info("Stopping server")
	return nil
}

// Shutdown завершает работу сервера и закрывает подключение к базе данных
func (s *Server) Shutdown() error {
	err := postgres.DbService.CloseDB()

	if err != nil {
		return err
	}

	return SongService.Close()
}
