package config

import (
	"net/http"
	"time"
)

const (
	EARTH_RADIUS_KM = 6371.0
	MAX_RETRIES     = 3
	TIMEOUT_SECONDS = 10
)

var (
	MapboxRateLimiter chan struct{}
	HTTPClient        *http.Client
)

func InitRateLimiter() {
	MapboxRateLimiter = make(chan struct{}, 5)
	for i := 0; i < 5; i++ {
		MapboxRateLimiter <- struct{}{}
	}

	HTTPClient = &http.Client{
		Timeout: TIMEOUT_SECONDS * time.Second,
	}
}
