package postgres

import (
	"database/sql"
	"fmt"
	"net/http"
	e "song/internal/presentation/customError"
	"song/internal/presentation/logger"
)

var (
	DbService *DB
)

// DB - структура для работы с базой данных
type DB struct {
	Db *sql.DB
}

// CreateDB создает подключение к базе данных и возвращает экземпляр DB
func CreateDB(ip, port, user, pass, nameDB string) (*DB, error) {
	logger.Logger.Debug("Database connection creating...")
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ip, port, user, pass, nameDB)
	conn, err := sql.Open("postgres", sqlInfo)

	if err != nil {
		return nil, &e.DbConnectionError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Database connection error: %v", err),
		}
	}

	DbService = &DB{
		Db: conn,
	}

	logger.Logger.Info("Database connection has been created")
	return DbService, nil
}

// CloseDB закрывает подключение к базе данных
func (db *DB) CloseDB() error {
	logger.Logger.Debug("Closing database connection")
	return db.Db.Close()
}
