package domain

import "time"

// SongText - алиас для текста песни
type SongText = string

// Page - алиас номер страницы
type Page = int

// SongDataByUser - данные о песне от пользователя
type SongDataByUser struct {
	Group string `json:"group"` // Группа или исполнитель
	Name  string `json:"song"`  // Название песни
}

// SongDataByApi - дополнительные данные о песне, полученные из поднятого API
type SongDataByApi struct {
	Date time.Time `json:"releaseDate"` // Дата выпуска песни
	Text string    `json:"text"`        // Текст песни
	Link string    `json:"link"`        // Ссылка на песню
}

// BaseError - шаблон ошибки
type BaseError struct {
	Err  string `json:"err"`  // Сообщение об ошибке
	Code int    `json:"code"` // Код ошибки
}

func (e *BaseError) Error() string {
	return e.Err
}

// RequestError - ошибка запроса к API
type RequestError = BaseError

type InputDataError = BaseError

// Id - алиас для идентификатора
type Id = uint64
