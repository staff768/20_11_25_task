package checkers

import (
	"net/http"
	"time"
)

const (
	StatusAvailable    = "available"
	StatusNotAvailable = "not available"
)

type Checker struct {
	client *http.Client
}

func NewChecker(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (checker *Checker) Check(url string) string {

	if url[:4] != "http" {
		url = "http://" + url
	}

	resp, err := checker.client.Head(url)
	if err != nil {
		return StatusNotAvailable
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return StatusAvailable
	}

	return StatusNotAvailable
}
