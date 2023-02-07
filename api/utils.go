package api

import (
	"fmt"
	"github.com/BRUHItsABunny/webfireAutoTrainer/constants"
	"net/http"
)

func DefaultHeaders(class *WFTClass, laravelAuth, xsrf string) http.Header {
	return http.Header{
		"User-Agent":      {constants.UserAgent},
		"Accept":          {"*/*"},
		"Sec-GPC":         {"1"},
		"Accept-Language": {"en-US,en;q=0.6"},
		"Origin":          {"https://webfiretraining.com"},
		"Sec-Fetch-Site":  {"same-origin"},
		"Sec-Fetch-Mode":  {"no-cors"},
		"Sec-Fetch-Dest":  {"empty"},
		"Referer":         {class.Referer},
		"Cookie":          {fmt.Sprintf("XSRF-TOKEN=%s; laravel_session=%s", xsrf, laravelAuth)},
	}
}
