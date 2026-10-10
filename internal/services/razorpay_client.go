package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const razorpayAPIBaseURL = "https://api.razorpay.com/v1"

type razorpayClient struct {
	keyID     string
	keySecret string
	baseURL   string
	http      *http.Client
}

type razorpayProviderOrder struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type razorpayProviderPayment struct {
	ID       string `json:"id"`
	OrderID  string `json:"order_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

func NewRazorpayClient(keyID, keySecret string) *razorpayClient {
	return &razorpayClient{
		keyID:     keyID,
		keySecret: keySecret,
		baseURL:   razorpayAPIBaseURL,
		http:      &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *razorpayClient) KeyID() string {
	return c.keyID
}

func (c *razorpayClient) CreateOrder(ctx context.Context, amount int64, receipt string) (*razorpayProviderOrder, error) {
	if err := c.validateCredentials(); err != nil {
		return nil, err
	}
	payload := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  receipt,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/orders", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.keyID, c.keySecret)

	var response razorpayProviderOrder
	if err := c.do(request, &response); err != nil {
		return nil, err
	}
	if response.ID == "" {
		return nil, errors.New("Razorpay returned an invalid order")
	}
	return &response, nil
}

func (c *razorpayClient) FetchPayment(ctx context.Context, paymentID string) (*razorpayProviderPayment, error) {
	if err := c.validateCredentials(); err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/payments/"+url.PathEscape(paymentID), nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(c.keyID, c.keySecret)
	var response razorpayProviderPayment
	if err := c.do(request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *razorpayClient) CapturePayment(ctx context.Context, paymentID string, amount int64) (*razorpayProviderPayment, error) {
	if err := c.validateCredentials(); err != nil {
		return nil, err
	}
	body, err := json.Marshal(map[string]interface{}{"amount": amount, "currency": "INR"})
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/payments/"+url.PathEscape(paymentID)+"/capture", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.keyID, c.keySecret)
	var response razorpayProviderPayment
	if err := c.do(request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *razorpayClient) validateCredentials() error {
	if strings.TrimSpace(c.keyID) == "" || strings.TrimSpace(c.keySecret) == "" {
		return errors.New("Razorpay credentials are not configured")
	}
	return nil
}

func (c *razorpayClient) do(request *http.Request, output interface{}) error {
	response, err := c.http.Do(request)
	if err != nil {
		return fmt.Errorf("Razorpay request failed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Razorpay returned HTTP %d", response.StatusCode)
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(output); err != nil {
		return fmt.Errorf("invalid Razorpay response: %w", err)
	}
	return nil
}
