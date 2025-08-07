package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/slash3b/cache"
	"go.uber.org/zap"
)

const retryCount = 8

type CineplexApi struct {
	client *http.Client
	cache  cache.LRU[string, CineplexMoviesResponse]
	lg     *zap.Logger
}

func NewCineplex(c *http.Client, lg *zap.Logger) *CineplexApi {
	return &CineplexApi{client: c}
}

func (s *CineplexApi) GetMovies(ctx context.Context) (CineplexMoviesResponse, error) {
	filter := FilterRequest{
		ByFormat:      "",
		ByQuality:     "",
		ByAudioFormat: "",
		ByCinema:      "",
		ByStatus:      "0",
	}

	response := CineplexMoviesResponse{}

	jsonData, err := json.Marshal(filter)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://cineplex.md/api/getMoviesFiltered", bytes.NewBuffer(jsonData))
	if err != nil {
		return response, err
	}

	var res *http.Response

	for r := range retryCount {
		res, err = s.client.Do(req)
		if err == nil {
			break
		}

		if err != nil {
			select {
			case <-ctx.Done():
				s.lg.Warn("context canceled", zap.Error(err))

				return response, ctx.Err()
			default:
				//nolint:gosec
				sleepTime := time.Duration(100<<r+rand.Intn(200)) * time.Millisecond

				s.lg.Warn("failed to fetch movies",
					zap.String("err", err.Error()),
					zap.String("sleepTime", sleepTime.String()),
					zap.Int("retryN", r),
				)

				time.Sleep(sleepTime)
			}
		}
	}

	if err != nil {
		cres, cacheErr := s.cache.Get("getMoviesFiltered")
		if cacheErr != nil {
			s.lg.Warn("cache is empty", zap.String("key", "getMoviesFiltered"))

			return response, err
		}

		s.lg.Info("fetched response from cache")

		return cres, nil
	}

	defer func() {
		_ = res.Body.Close()
	}()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return response, err
	}

	s.cache.Set("getMoviesFiltered", response)

	return response, nil
}
