// @title Song API
// @version 1.0
// @description This is a sample server for a song management application.
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"net/url"
	"os"
	"song/internal/presentation/logger"
	"song/internal/presentation/postgres"
	"song/internal/presentation/realization"
	"song/internal/presentation/server"
	"song/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	err := logger.NewLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = godotenv.Load("../../../.env")
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to load env fail - %v", err))
		return
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	redisHost := os.Getenv("REDIS_HOST")
	fmt.Println(redisHost)
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASSWORD")

	apiUrl := os.Getenv("API_URL")
	apiPort := os.Getenv("API_PORT")

	serverPort := os.Getenv("SERVER_PORT")

	_, err = postgres.CreateDB(host, port, user, password, name)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Database creating error - %v", err))
		return
	}

	err = postgres.DbService.CreateSchema()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Schema creating error - %v", err))
		return
	}

	songRepo := realization.NewSongRepo()
	cacheRepo := realization.NewConnectRedis(redisHost, redisPort, redisPass)
	songService := services.NewSongService(cacheRepo, songRepo)

	api, err := url.Parse(fmt.Sprintf("%s:%s", apiUrl, apiPort))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Api url parsing error - %v", err))
		return
	}

	srv := server.NewServer()
	err = srv.Start(songService, api, serverPort)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Server working error - %v", err))
	}

	err = srv.Shutdown()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Server stopping error - %v", err))
	}
}
