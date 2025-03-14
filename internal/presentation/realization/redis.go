package realization

import (
	"context"
	"fmt"
	"net/http"
	"song/internal/domain"
	e "song/internal/presentation/customError"
	"song/internal/presentation/logger"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// RedisRepo представляет репозиторий для работы с Redis
type RedisRepo struct {
	db *redis.Client
}

// NewConnectRedis создает новое подключение к Redis
// addr - адрес Redis сервера
// pass - пароль для подключения к Redis
func NewConnectRedis(addr, port, pass string) *RedisRepo {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: pass,
		DB:       0,
	})

	logger.Logger.Info("Redis connection create")
	return &RedisRepo{
		db: r,
	}
}

// CreateKey создает новый ключ в Redis
// id - идентификатор ключа
// text - текст песни
func (r *RedisRepo) CreateKey(id domain.Id, text domain.SongText) error {
	ctx := context.Background()
	res := r.db.Set(ctx, strconv.FormatUint(id, 10), text, 0)
	_, err := res.Result()

	if err != nil {
		return &e.RedisQueryError{
			Err:  fmt.Sprintf("Ошибка создания ключа в Redis: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	logger.Logger.Debug(fmt.Sprintf("key: %d was created", id))
	return nil
}

// GetText получает значение ключа из Redis по идентификатору
// id - идентификатор ключа
func (r *RedisRepo) GetText(id domain.Id) (*domain.SongText, error) {
	ctx := context.Background()
	res := r.db.Get(ctx, strconv.FormatUint(id, 10))

	err := res.Err()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, &e.RedisQueryError{
			Err:  fmt.Sprintf("Ошибка получения ключа в Redis: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	logger.Logger.Debug(fmt.Sprintf("key: %d was got", id))
	text := domain.SongText(res.Val())
	return &text, nil
}

// DelKey удаляет ключ из Redis
// id - идентификатор ключа
func (r *RedisRepo) DelKey(id domain.Id) error {
	ctx := context.Background()
	err := r.db.Del(ctx, strconv.FormatUint(id, 10)).Err()

	if err != nil {
		return &e.RedisQueryError{
			Err:  fmt.Sprintf("Ошибка удаления ключа в Redis: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	logger.Logger.Debug(fmt.Sprintf("key: %d was deleted", id))
	return nil
}

func (r *RedisRepo) Close() error {
	logger.Logger.Info("Redis connection was closed")
	return r.db.Close()
}
