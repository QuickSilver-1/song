package realization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"song/internal/domain"
	e "song/internal/presentation/customError"
	"song/internal/presentation/logger"
	"song/internal/presentation/postgres"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

// SongRepo - реализация репозитория для работы с песнями в базе данных
type SongRepo struct{}

func NewSongRepo() *SongRepo {
	return &SongRepo{}
}

// GetLib получает библиотеку песен по фильтру и номеру страницы
// filter - фильтр для песен
// page - номер страницы для пагинации
func (r *SongRepo) GetLib(filter domain.Song, page domain.Page) (*[]domain.Song, error) {
	query := sq.Select("id", "group_name", "song_name", "release_date", "text", "link").From("song")
	count := 1

	if filter.ID != 0 {
		query = query.Where(fmt.Sprintf("id = $%d", count), filter.ID)
	} else {
		if filter.Name != "" {
			query = query.Where(fmt.Sprintf("song_name LIKE $%d", count), "%"+filter.Name+"%")
			count++
		}
		if filter.Group != "" {
			query = query.Where(fmt.Sprintf("group_name LIKE $%d", count), "%"+filter.Group+"%")
			count++
		}
		if !filter.Date.IsZero() {
			query = query.Where(fmt.Sprintf("release_date = $%d", count), filter.Date)
			count++
		}
		if filter.Text != "" {
			query = query.Where(fmt.Sprintf("text LIKE $%d", count), "%"+filter.Text+"%")
			count++
		}
		if filter.Link != "" {
			query = query.Where(fmt.Sprintf("link LIKE $%d", count), "%"+filter.Link+"%")
		}
	}

	query = query.Offset(uint64(20 * (page - 1))).Limit(20)
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, &e.DbQueryError{
			Err:  fmt.Sprintf("Ошибка генерации SQL-запроса: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := postgres.DbService.Db.QueryContext(timeoutCtx, sqlQuery, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &e.RowsNotFoundError{
				Err:  "Песни с таким фильтром не существуют",
				Code: http.StatusBadRequest,
			}
		}
		return nil, &e.DbQueryError{
			Err:  fmt.Sprintf("Ошибка выполнения запроса к базе данных: %v", err),
			Code: http.StatusInternalServerError,
		}
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Logger.Error(fmt.Sprintf("Error closing response body: %v", err))
		}
	}()

	var result []domain.Song
	for rows.Next() {
		var song domain.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Name, &song.Date, &song.Text, &song.Link); err != nil {
			return nil, &e.DbQueryError{
				Err:  fmt.Sprintf("Ошибка сканирования строки: %v", err),
				Code: http.StatusInternalServerError,
			}
		}
		result = append(result, song)
	}

	return &result, nil
}

// GetText получает текст песни по идентификатору
// id - идентификатор песни
func (r *SongRepo) GetText(id domain.Id) (*domain.SongText, error) {
	var text string
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := postgres.DbService.Db.QueryRowContext(timeoutCtx, `SELECT text FROM song WHERE id = $1`, id).Scan(&text)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &e.RowsNotFoundError{
				Err:  "Песня с таким идентификатором не существует",
				Code: http.StatusBadRequest,
			}
		}
		return nil, &e.DbQueryError{
			Err:  fmt.Sprintf("Ошибка выполнения запроса к базе данных: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	return &text, nil
}

// DelSong удаляет песню по идентификатору
// id - идентификатор песни
func (r *SongRepo) DelSong(id domain.Id) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postgres.DbService.Db.ExecContext(timeoutCtx, "DELETE FROM song WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &e.RowsNotFoundError{
				Err:  "Песня с таким идентификатором не существует",
				Code: http.StatusBadRequest,
			}
		}
		return &e.DbQueryError{
			Err:  fmt.Sprintf("Ошибка выполнения запроса к базе данных: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

// ChangeSong изменяет данные песни
// song - объект песни с новыми данными
func (r *SongRepo) ChangeSong(song domain.Song) error {
	query := sq.Update("song").Where("id = ?", song.ID)

	if song.Name != "" {
		query = query.Set("song_name", song.Name)
	}
	if song.Group != "" {
		query = query.Set("group_name", song.Group)
	}
	if !song.Date.IsZero() {
		query = query.Set("release_date", song.Date)
	}
	if song.Text != "" {
		query = query.Set("text", song.Text)
	}
	if song.Link != "" {
		query = query.Set("link", song.Link)
	}

	sqlQuery, args, err := query.ToSql()
	placeholderCount := strings.Count(sqlQuery, "?")
	for i := 0; i <= placeholderCount; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", fmt.Sprintf("$%d", i+1), 1)
	}

	if err != nil {
		return &e.DbQueryError{
			Err:  fmt.Sprintf("Sql query generation error: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = postgres.DbService.Db.ExecContext(timeoutCtx, sqlQuery, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &e.RowsNotFoundError{
				Err:  "Song with this id not exist",
				Code: http.StatusBadRequest,
			}
		}
		return &e.DbQueryError{
			Err:  fmt.Sprintf("DB query error: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

// CreateSong создает новую песню
// song - объект новой песни
func (r *SongRepo) CreateSong(song domain.Song) (*domain.Id, error) {
	var id uint64
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := postgres.DbService.Db.QueryRowContext(timeoutCtx, `INSERT INTO song (song_name, group_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`, song.Name, song.Group, song.Date, song.Text, song.Link).Scan(&id)

	if err != nil {
		return nil, &e.DbQueryError{
			Err:  fmt.Sprintf("Ошибка выполнения запроса к базе данных: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	return &id, nil
}
