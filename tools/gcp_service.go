package tools

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

func GetGoogleService(creds *google.Credentials, ctx context.Context) (*container.Service, error) {

	service, err := container.NewService(ctx, option.WithHTTPClient(oauth2.NewClient(ctx, creds.TokenSource)))

	if err != nil {
		return nil, err
	}

	return service, err
}
