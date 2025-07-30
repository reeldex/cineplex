package sender

import (
	"context"
	"go.uber.org/zap"
	"scraper/internal/services/fetcher"
	"scraper/pkg/telemetry"
	"sync"
)

type MovieFetcher interface {
	GetMovies() (fetcher.CineplexMoviesResponse, error)
}

type Scraper struct {
	s MovieFetcher
	l zap.Logger
}

func NewScraper(mv MovieFetcher, l zap.Logger) *Scraper {
	return &Scraper{
		s: mv,
		l: l,
	}
}

func (c *Scraper) Publish(ctx context.Context) (context.Context, <-chan dto.FilmResponse) {}

func (c *Scraper) GetFilms(ctx context.Context) (context.Context, <-chan dto.FilmResponse) {
	response := make(chan dto.FilmResponse)

	rawMovies, err := c.s.GetMovies()
	if err != nil {
		telemetry.ScraperErrors.WithLabelValues("unable_to_scrape_movies").Inc()
		c.l.Println("unable to scrape movies", err)

		close(response)
		return ctx, response
	}

	var wg sync.WaitGroup
	for _, rawMovie := range rawMovies {

		wg.Add(1)

		go func(movieData dto.RawFilmData) {
			defer wg.Done()
			movieModel, err := converter.FromDTO(movieData)
			if err != nil {
				telemetry.ScraperErrors.WithLabelValues("unable_to_convert_raw_movie_dto").Inc()
				c.l.Println("Unable to convert raw movie DTO", err)
				response <- dto.FilmResponse{Film: movieModel, Error: err}
				return
			}

			if !c.r.IsExists(movieModel) {
				_, err = c.r.Insert(movieModel)
				if err != nil {
					telemetry.ScraperErrors.WithLabelValues("unable_to_record_movie").Inc()
					response <- dto.FilmResponse{Film: movieModel, Error: err}
					c.l.Println("unable to insert new rawMovie", movieData, err)
					return
				}

				select {
				case <-ctx.Done():
				case response <- dto.FilmResponse{Film: movieModel}:
				}
			}
		}(rawMovie)

	}

	// goroutine that waits for all others to complete and safely close the channel
	go func() {
		wg.Wait()
		close(response)
	}()

	return ctx, response
}
