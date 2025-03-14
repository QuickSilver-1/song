package domain

import "time"

// Song - объект песни
type Song struct {
	ID    uint64    `json:"id"`          // Идентификатор песни
	Name  string    `json:"name"`        // Название песни
	Group string    `json:"group"`       // Группа или исполнитель
	Date  time.Time `json:"releaseDate"` // Дата выпуска песни
	Text  string    `json:"text"`        // Текст песни
	Link  string    `json:"link"`        // Ссылка на песню
}

// NewSong создает новый объект Song
func NewSong(name, group, text, link string, date time.Time) *Song {
	return &Song{
		Name:  name,
		Group: group,
		Date:  date,
		Text:  text,
		Link:  link,
	}
}
