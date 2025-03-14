package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"song/internal/domain"
	e "song/internal/presentation/customError"
	"song/internal/presentation/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TIME_FORMAT = "02.01.2006"
)

// Handlers определяет хендлеры для обработки HTTP-запросов
type Handlers struct{}

// NewHandlers создает новый экземпляр Handlers
func NewHandlers() *Handlers {
	return &Handlers{}
}

// @Summary		Get library
// @Description	Get songs library
// @Tags			library
// @Accept			json
// @Produce		json
// @Param			song		query		string	false	"Song name"
// @Param			group		query		string	false	"Group name"
// @Param			releaseDate	query		string	false	"Release date in format dd.mm.yyyy"
// @Param			text		query		string	false	"Song text"
// @Param			link		query		string	false	"Link"
// @Param			page		query		int		false	"Page number"
// @Success		200			{array}		domain.Song
// @Failure		400			{object}	map[string]string
// @Failure		500			{object}	map[string]string
// @Router			/lib [get]
func (h *Handlers) GetLib(ctx *gin.Context) {
	song, err := parseSong(ctx)
	if err != nil {
		answerError(ctx, err)
		return
	}

	pageStr := ctx.Request.URL.Query().Get("page")
	var page int = 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			answerError(ctx, &e.InvalidInputData{
				Err:  "Invalid page",
				Code: http.StatusBadRequest,
			})
			return
		}
	}

	lib, err := SongService.GetLib(*song, page)

	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, lib)
}

// @Summary		Get song text
// @Description	Get the text of a song
// @Tags			song
// @Accept			json
// @Produce		json
// @Param			id		query		uint64	true	"Song ID"
// @Param			page	query		int		false	"Page number"
// @Success		200		{object}	domain.SongText
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			/text [get]
func (h *Handlers) GetText(ctx *gin.Context) {
	idStr := ctx.Request.URL.Query().Get("id")
	var id uint64
	var err error
	if idStr == "" {
		answerError(ctx, &e.InvalidInputData{
			Err:  "id is a required parameter",
			Code: http.StatusBadRequest,
		})
		return
	} else {
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			answerError(ctx, &e.InvalidInputData{
				Err:  "Invalid id",
				Code: http.StatusBadRequest,
			})
			return
		}
	}

	pageStr := ctx.Request.URL.Query().Get("page")
	var page int = 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			answerError(ctx, &e.InvalidInputData{
				Err:  "Invalid page",
				Code: http.StatusBadRequest,
			})
			return
		}
	}

	text, err := SongService.GetText(id, page)
	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, text)
}

// @Summary		Delete song
// @Description	Delete a song by ID
// @Tags			song
// @Accept			json
// @Produce		json
// @Param			id	query		uint64	true	"Song ID"
// @Success		200	{object}	nil
// @Failure		400	{object}	map[string]string
// @Failure		500	{object}	map[string]string
// @Router			/song [delete]
func (h *Handlers) DelSong(ctx *gin.Context) {
	idStr := ctx.Request.URL.Query().Get("id")
	var id uint64
	var err error
	if idStr == "" {
		answerError(ctx, &e.InvalidInputData{
			Err:  "id is a required parameter",
			Code: http.StatusBadRequest,
		})
		return
	} else {
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			answerError(ctx, &e.InvalidInputData{
				Err:  "Invalid id",
				Code: http.StatusBadRequest,
			})
			return
		}
	}

	err = SongService.DelSong(id)
	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// @Summary		Change song
// @Description	Change song details by ID
// @Tags			song
// @Accept			json
// @Produce		json
// @Param			id		query		uint64		true	"Song ID"
// @Param			body	body		domain.Song	true	"Song details"
// @Success		200		{object}	nil
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			/song [patch]
func (h *Handlers) ChangeSong(ctx *gin.Context) {
	idStr := ctx.Request.URL.Query().Get("id")
	var id uint64
	var err error
	if idStr == "" {
		answerError(ctx, &e.InvalidInputData{
			Err:  "id is a required parameter",
			Code: http.StatusBadRequest,
		})
		return
	} else {
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			answerError(ctx, &e.InvalidInputData{
				Err:  "Invalid id",
				Code: http.StatusBadRequest,
			})
			return
		}
	}

	var song domain.Song
	err = json.NewDecoder(ctx.Request.Body).Decode(&song)
	defer func() {
		if err := ctx.Request.Body.Close(); err != nil {
			logger.Logger.Error(fmt.Sprintf("Error closing response body: %v", err))
		}
	}()
	if err != nil {
		answerError(ctx, &e.InvalidInputData{
			Err:  "Invalid body",
			Code: http.StatusBadRequest,
		})
		return
	}

	if song.Date.IsZero() && song.Group == "" && song.Link == "" && song.Name == "" && song.Text == "" {
		answerError(ctx, &e.InvalidInputData{
			Err:  "Invalid body",
			Code: http.StatusBadRequest,
		})
		return
	}

	song.ID = id
	err = SongService.ChangeSong(song)
	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// @Summary		Create song
// @Description	Create a new song
// @Tags			song
// @Accept			json
// @Produce		json
// @Param			body	body		domain.SongDataByUser	true	"Song details"
// @Success		200		{object}	map[string]domain.Id
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			/song [post]
func (h *Handlers) CreateSong(ctx *gin.Context) {
	var data domain.SongDataByUser
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	defer func() {
		if err := ctx.Request.Body.Close(); err != nil {
			logger.Logger.Error(fmt.Sprintf("Error closing response body: %v", err))
		}
	}()
	if err != nil {
		answerError(ctx, &e.InvalidInputData{
			Err:  "Invalid body",
			Code: http.StatusBadRequest,
		})
		return
	}

	id, err := SongService.CreateSong(data, ApiUrl)
	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]domain.Id{"song_id": *id})
}

func parseSong(ctx *gin.Context) (*domain.Song, error) {
	idStr := ctx.Request.URL.Query().Get("id")
	var id uint64
	var err error
	if idStr != "" {
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, &e.InvalidInputData{
				Err:  "Invalid id",
				Code: http.StatusBadRequest,
			}
		}
	}

	name := ctx.Request.URL.Query().Get("song")
	group := ctx.Request.URL.Query().Get("group")
	dateStr := ctx.Request.URL.Query().Get("releaseDate")
	var date time.Time
	if dateStr != "" {
		date, err = time.Parse(TIME_FORMAT, dateStr)
		if err != nil {
			return nil, &e.InvalidInputData{
				Err:  "Invalid date format, correct format - 16.07.2006",
				Code: http.StatusBadRequest,
			}
		}
	}

	text := ctx.Request.URL.Query().Get("text")
	link := ctx.Request.URL.Query().Get("link")

	return &domain.Song{
		ID:    id,
		Name:  name,
		Group: group,
		Date:  date,
		Text:  text,
		Link:  link,
	}, nil
}

// answerError обрабатывает ошибки и возвращает соответствующий HTTP-статус
func answerError(ctx *gin.Context, err error) {
	baseErr := err.(*domain.BaseError)

	switch baseErr.Code {
	case http.StatusInternalServerError:
		logger.Logger.Error(baseErr.Error())
		ctx.JSON(http.StatusInternalServerError, map[string]string{"errors": STATUS_INTERNAL_SERVER})
	case http.StatusBadRequest:
		logger.Logger.Debug("Invalid data from user")
		ctx.JSON(http.StatusBadRequest, map[string]string{"errors": fmt.Sprintf("%s: %s", STATUS_BAD_REQUEST, baseErr.Error())})
	}

	ctx.Abort()
}
