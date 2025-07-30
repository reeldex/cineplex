package fetcher

import (
	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
	"net/http"
	"scraper/dto"
)

type SoupDecorator struct {
	c *http.Client
}

func NewSoupDecorator(client *http.Client) *SoupDecorator {
	return &SoupDecorator{
		c: client,
	}
}

func (s *SoupDecorator) GetMovies() ([]dto.RawFilmData, error) {
	var response []dto.RawFilmData

	_, err := soup.GetWithClient("https://cineplex.md/lang/en", s.c)
	if err != nil {
		return response, err
	}

	movies, err := get(s.c)
	if err != nil {
		return response, err
	}

	for _, m := range movies.Movies {
		id, err := uuid.Parse(m.UUID)
		if err != nil {
			// todo: log
			continue
		}

		filmDTO := dto.RawFilmData{
			ID:    id,
			Title: m.OriginalTitle,
			Lang:  m.Languages,
			Date:  m.CinemaDate,
		}

		response = append(response, filmDTO)
	}

	return response, nil
}
