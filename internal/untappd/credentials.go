package untappd

import (
	"errors"
	"os"
)

type credentials struct {
	ClientID     string
	ClientSecret string
}

func getCredentials() (*credentials, error) {
	cid := os.Getenv("UNTAPPD_CLIENT_ID")
	if cid == "" {
		return nil, errors.New("missing environment variable: UNTAPPD_CLIENT_ID")
	}

	cs := os.Getenv("UNTAPPD_CLIENT_SECRET")
	if cs == "" {
		return nil, errors.New("missing environment variable: UNTAPPD_CLIENT_SECRET")
	}

	return &credentials{ClientID: cid, ClientSecret: cs}, nil
}
