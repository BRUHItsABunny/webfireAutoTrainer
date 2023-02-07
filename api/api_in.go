package api

import (
	"fmt"
	"github.com/BRUHItsABunny/gOkHttp/responses"
	"net/http"
)

func BytesParser(resp *http.Response) ([]byte, error) {
	err := responses.CheckHTTPCode(resp, 200)
	if err != nil {
		return nil, fmt.Errorf("responses.CheckHTTPCode: %w", err)
	}

	result, err := responses.ResponseBytes(resp)
	if err != nil {
		return nil, fmt.Errorf("responses.ResponseJSON: %w", err)
	}
	return result, nil
}
