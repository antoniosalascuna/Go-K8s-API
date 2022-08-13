package tools

import (
	"context"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
)

func GetCredentialFromJson(context context.Context, jsonUrl string) (*google.Credentials, error) {

	data, err := ioutil.ReadFile(jsonUrl)

	if err != nil {
		return nil, err
	}

	cred, err := google.CredentialsFromJSON(context, data, container.CloudPlatformScope)

	if err != nil {
		return nil, err
	}

	return cred, err

}
