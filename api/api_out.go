package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BRUHItsABunny/gOkHttp/requests"
	"github.com/BRUHItsABunny/webfireAutoTrainer/constants"
	"net/http"
	"net/url"
)

func GetParametersRequest(ctx context.Context, class *WFTClass, laravelAuth, xsrf string) (*http.Request, error) {
	parameters := url.Values{
		"session_id": {class.SessionID},
		"version":    {"3.5"},
		"command":    {"GETPARAM"},
		"aicc_data":  {""},
	}

	req, err := requests.MakePOSTRequest(ctx, constants.EndpointStoreMarks,
		requests.NewHeaderOption(DefaultHeaders(class, laravelAuth, xsrf)),
		requests.NewURLParamOption(parameters),
	)
	if err != nil {
		return nil, fmt.Errorf("requests.MakePOSTRequest: %w", err)
	}
	return req, nil
}

func PutParametersRequest(ctx context.Context, class *WFTClass, laravelAuth, xsrf string, data *AICCData) (*http.Request, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	parameters := url.Values{
		"session_id": {class.SessionID},
		"version":    {"3.5"},
		"command":    {"PUTPARAM"},
		"aicc_data":  {string(dataStr)},
	}

	req, err := requests.MakePOSTRequest(ctx, constants.EndpointStoreMarks,
		requests.NewHeaderOption(DefaultHeaders(class, laravelAuth, xsrf)),
		requests.NewURLParamOption(parameters),
	)
	if err != nil {
		return nil, fmt.Errorf("requests.MakePOSTRequest: %w", err)
	}
	return req, nil
}

func ExitRequest(ctx context.Context, class *WFTClass, laravelAuth, xsrf string) (*http.Request, error) {
	parameters := url.Values{
		"session_id": {class.SessionID},
		"version":    {"3.5"},
		"command":    {"EXITAU"},
		"aicc_data":  {""},
	}

	req, err := requests.MakePOSTRequest(ctx, constants.EndpointStoreMarks,
		requests.NewHeaderOption(DefaultHeaders(class, laravelAuth, xsrf)),
		requests.NewURLParamOption(parameters),
	)
	if err != nil {
		return nil, fmt.Errorf("requests.MakePOSTRequest: %w", err)
	}
	return req, nil
}
