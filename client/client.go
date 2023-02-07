package client

import (
	"context"
	"fmt"
	"github.com/BRUHItsABunny/webfireAutoTrainer/api"
	"net/http"
)

type WFTClient struct {
	Client      *http.Client
	LaravelAuth string
	XSRF        string
}

func NewWFTClient(hClient *http.Client, laravelAuth, xsrf string) *WFTClient {
	if hClient == nil {
		hClient = http.DefaultClient
	}

	return &WFTClient{
		Client:      hClient,
		LaravelAuth: laravelAuth,
		XSRF:        xsrf,
	}
}

func (c *WFTClient) StartClass(ctx context.Context, class *api.WFTClass) ([]byte, error) {
	req, err := api.GetParametersRequest(ctx, class, c.LaravelAuth, c.XSRF)
	if err != nil {
		return nil, fmt.Errorf("api.GetParametersRequest: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Do: %w", err)
	}

	result, err := api.BytesParser(resp)
	if err != nil {
		return nil, fmt.Errorf("api.BytesParser: %w", err)
	}
	return result, nil
}

func (c *WFTClient) InitClass(ctx context.Context, class *api.WFTClass) ([]byte, error) {
	req, err := api.PutParametersRequest(ctx, class, c.LaravelAuth, c.XSRF, api.NewAICCData(api.LessonStatusInitialized))
	if err != nil {
		return nil, fmt.Errorf("api.PutParametersRequest: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Do: %w", err)
	}

	result, err := api.BytesParser(resp)
	if err != nil {
		return nil, fmt.Errorf("api.BytesParser: %w", err)
	}
	return result, nil
}

func (c *WFTClient) FinishClass(ctx context.Context, class *api.WFTClass) ([]byte, error) {
	req, err := api.PutParametersRequest(ctx, class, c.LaravelAuth, c.XSRF, api.NewAICCData(api.LessonStatusProgressed))
	if err != nil {
		return nil, fmt.Errorf("api.PutParametersRequest: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Do: %w", err)
	}

	result, err := api.BytesParser(resp)
	if err != nil {
		return nil, fmt.Errorf("api.BytesParser: %w", err)
	}
	return result, nil
}

func (c *WFTClient) ExitClass(ctx context.Context, class *api.WFTClass) ([]byte, error) {
	req, err := api.ExitRequest(ctx, class, c.LaravelAuth, c.XSRF)
	if err != nil {
		return nil, fmt.Errorf("api.ExitClassRequest: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Do: %w", err)
	}

	result, err := api.BytesParser(resp)
	if err != nil {
		return nil, fmt.Errorf("api.BytesParser: %w", err)
	}
	return result, nil
}

func (c *WFTClient) FinishExam(ctx context.Context, class *api.WFTClass, score float64) ([]byte, error) {
	req, err := api.PutParametersRequest(ctx, class, c.LaravelAuth, c.XSRF, api.NewAICCDataWithScore(api.LessonStatusProgressed, score))
	if err != nil {
		return nil, fmt.Errorf("api.PutParametersRequest: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Do: %w", err)
	}

	result, err := api.BytesParser(resp)
	if err != nil {
		return nil, fmt.Errorf("api.BytesParser: %w", err)
	}
	return result, nil
}
