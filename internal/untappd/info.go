package untappd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// InfoResponse defines the shape of the response from the UntappdAPI /v4/beer/info/:beerID endpoint
type InfoResponse struct {
	Response struct {
		Beer *Beer `json:"beer"`
	} `json:"response"`
}

// Title produces a slack formatted title for the InfoResponse
func (i *InfoResponse) Title() string {
	return title(i.Response.Beer, i.Response.Beer.Brewery)
}

// Text produces slack formatted text for the InfoResponse
func (i *InfoResponse) Text() string {
	return text(i.Response.Beer)
}

// Info calls the UntappdAPI /v4/beer/info/:beerID endpoint
func Info(id string) (*InfoResponse, error) {
	creds, err := getCredentials()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf("https://api.untappd.com/v4/beer/info/%s", id))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	q := u.Query()
	q.Set("client_id", creds.ClientID)
	q.Set("client_secret", creds.ClientSecret)
	q.Set("compact", "true")
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var ir InfoResponse
	if err := json.NewDecoder(res.Body).Decode(&ir); err != nil {
		return nil, errors.WithStack(err)
	}

	return &ir, nil
}
