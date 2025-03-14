package domain

import (
	"testing"
	"time"
)

// TestNewSong проверяет функцию NewSong
func TestNewSong(t *testing.T) {
	name := "Test Song"
	group := "Test Group"
	text := "Test Text"
	link := "http://testlink.com"
	date := time.Now()

	song := NewSong(name, group, text, link, date)

	if song.Name != name {
		t.Errorf("Expected Name to be %s, but got %s", name, song.Name)
	}

	if song.Group != group {
		t.Errorf("Expected Group to be %s, but got %s", group, song.Group)
	}

	if !song.Date.Equal(date) {
		t.Errorf("Expected Date to be %v, but got %v", date, song.Date)
	}

	if song.Text != text {
		t.Errorf("Expected Text to be %s, but got %s", text, song.Text)
	}

	if song.Link != link {
		t.Errorf("Expected Link to be %s, but got %s", link, song.Link)
	}
}
