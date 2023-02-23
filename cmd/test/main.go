package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/k0kubun/pp/v3"
	"github.com/streemtech/tiltify/v5/api"
)

func main() {
	cmd, err := api.NewClientWithResponses("https://v5api.tiltify.com", true, api.WithRequestEditorFn(ref))
	if err != nil {
		log.Fatal(err.Error())
	}

	resp, err := cmd.CampaignDonations(context.Background(), "7d47a0ee-49cf-4ed2-9996-ec77c1f847bd", &api.CampaignDonationsParams{
		Limit: aws.Int(10),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	pp.Println(resp)
}

func ref(ctx context.Context, req *http.Request) error {
	Token, _ := os.LookupEnv("TILTIFY_TOKEN")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))
	return nil
}
