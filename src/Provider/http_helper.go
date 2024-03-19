package Provider

import (
	"errors"
	"io"
	"net/http"
)

func getHttpBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Check the status code
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New(resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	// Read the response body
	return io.ReadAll(resp.Body)
}
