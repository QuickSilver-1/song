package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"song/internal/domain"
	"song/internal/presentation/logger"
	"song/internal/presentation/postgres"
	"song/internal/presentation/realization"
	"song/internal/presentation/server"
	"song/internal/services"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	HOST string
	PORT string
	USER string
	PASS string
	NAME string
	Srv  *server.Server

	redisHost string
	redisPort string
	redisPass string

	apiUrl  string
	apiPort string
)

func init() {
	_ = logger.NewLogger()
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	HOST = os.Getenv("DB_HOST")
	PORT = os.Getenv("DB_PORT")
	USER = os.Getenv("DB_USER")
	PASS = os.Getenv("DB_PASSWORD")
	NAME = os.Getenv("DB_NAME")

	redisHost = os.Getenv("REDIS_HOST")
	redisPort = os.Getenv("REDIS_PORT")
	redisPass = os.Getenv("REDIS_PASSWORD")

	apiUrl = os.Getenv("API_URL")
	apiPort = os.Getenv("API_PORT")

	Srv = StartTestServer()
}

func StartTestServer() *server.Server {
	// Настройка базы данных
	_, err := postgres.CreateDB(HOST, PORT, USER, PASS, NAME)
	if err != nil {
		log.Fatalf("Could not create database connection: %v", err)
	}

	// Настройка сервисов
	cacheRepo := realization.NewConnectRedis(redisHost, redisPort, redisPass)
	songRepo := &realization.SongRepo{}
	songService := services.NewSongService(cacheRepo, songRepo)

	// Запуск сервера
	srv := server.NewServer()
	go func() {
		api, _ := url.Parse(fmt.Sprintf("%s:%s", apiUrl, apiPort))
		if err := srv.Start(songService, api, "8080"); err != nil {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Ждем запуска сервера
	time.Sleep(3 * time.Second)

	return srv
}

func TestGetLib(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/lib?page=1", nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []domain.Song
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestDelSong(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/song?id=1", nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
