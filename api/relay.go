package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type RelayResponseObject struct {
	Amount     string            `json:"amount"`
	Status     string            `json:"status"`
	UUID       string            `json:"uuid"`
	MetaRaw    string            `json:"meta"`
	MetaParsed map[string]string `json:"-"`
}

// CheckRelay will check the V3 relay with the given relay key etc.
// This endpoint will do a best effort conversion of the metadata to a map of string,string
// If the conversion fails, the data will still be returned, but an error will also be returned
func CheckRelay(ctx context.Context, transport http.RoundTripper, baseURL string, relay string, uuid uuid.UUID, accessKey string) (rr *RelayResponseObject, err error) {

	re := resty.New()
	re.SetBaseURL(baseURL)

	r := re.NewRequest()
	r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessKey))
	r.SetContext(ctx)

	rr = new(RelayResponseObject)

	r.SetPathParams(map[string]string{
		"uuid":     uuid.String(),
		"provider": relay,
	})

	r.SetResult(rr)

	resp, err := r.Get("/api/v3/relays/{provider}/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("error in post")
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got bad response code: %d, %s", resp.StatusCode(), resp.Body())
	}

	if rr.MetaRaw == "" {
		return rr, nil
	}
	d := new(map[string]string)

	err = json.Unmarshal([]byte(rr.MetaRaw), d)
	if err != nil {
		return rr, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return rr, nil
}
