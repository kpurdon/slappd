package untappd

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// SearchResponse defines the shape of the response from the UntappdAPI /v4/search/beer endpoint
type SearchResponse struct {
	Response struct {
		Beers struct {
			Items []*SearchItem
		}
	}
}

// SearchItem defines the shape if the nested Items in SearchResponse
type SearchItem struct {
	Beer    *Beer
	Brewery *Brewery
}

// Title produces a slack formatted title for the SearchItem
func (i *SearchItem) Title() string {
	return title(i.Beer, i.Brewery)
}

// Text produces slack formatted text for the SearchItem
func (i *SearchItem) Text() string {
	return text(i.Beer)
}

// Search calls the UntappdAPI /v4/search/beer endpoint
func Search(searchStr string) (*SearchResponse, error) {
	creds, err := getCredentials()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("https://api.untappd.com/v4/search/beer")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	q := u.Query()
	q.Set("client_id", creds.ClientID)
	q.Set("client_secret", creds.ClientSecret)
	q.Set("q", searchStr)
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	var ur SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&ur); err != nil {
		return nil, err
	}

	return &ur, nil
}
