package fetcher

import (
	"context"
	"net/http"
)

type SoupDecorator struct {
	c *http.Client
}

func NewSoupDecorator(client *http.Client) (*SoupDecorator, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://cineplex.md", nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return &SoupDecorator{
		c: client,
	}, nil
}

func (s *SoupDecorator) GetMovies(ctx context.Context) (CineplexMoviesResponse, error) {
	return get(ctx, s.c)
}
