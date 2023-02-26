package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type RelayResponseObject struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data struct {
		Amount       float64   `json:"amount"`
		Currency     string    `json:"currency"`
		DonorName    string    `json:"donor_name"`
		DonorComment string    `json:"donor_comment"`
		Status       string    `json:"status"`
		UUID         uuid.UUID `json:"uuid"`
		Meta         struct {
			UUID uuid.UUID `json:"uuid"`
			Time time.Time `json:"time"`
			Hash string    `json:"hash"`
		} `json:"meta"`
	} `json:"data"`
}

// CheckRelay will check the V3 relay with the given relay key etc.
// This endpoint will do a best effort conversion of the metadata to a map of string,string
// If the conversion fails, the data will still be returned, but an error will also be returned
func CheckRelay(ctx context.Context, transport http.RoundTripper, baseURL string, relay string, uuid uuid.UUID, accessKey string) (rr *RelayResponseObject, err error) {

	re := resty.New()
	re.SetBaseURL(baseURL)

	re.SetTransport(transport)
	r := re.NewRequest()
	r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessKey))
	r.SetContext(ctx)

	rr = new(RelayResponseObject)

	r.SetPathParams(map[string]string{
		"uuid":     uuid.String(),
		"provider": relay,
	})

	// r.SetResult(rr)

	resp, err := r.Get("/api/v3/relays/{provider}/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("error in post: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got bad response code: %d, %s", resp.StatusCode(), resp.Body())
	}

	body := resp.Body()
	err = json.Unmarshal(body, rr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json (%s): %w", body, err)
	}

	return rr, nil
}
