package fetcher

import (
	"context"
	"net/http"
)

type SoupDecorator struct {
	c *http.Client
}

func NewSoupDecorator(client *http.Client) *SoupDecorator {
	_, err := client.Get("https://cineplex.md")
	_ = err

	return &SoupDecorator{
		c: client,
	}
}

func (s *SoupDecorator) GetMovies(ctx context.Context) (CineplexMoviesResponse, error) {
	return get(s.c)
}
