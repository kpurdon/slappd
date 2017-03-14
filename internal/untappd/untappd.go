package untappd

import (
	"fmt"
	"regexp"
	"strings"
)

// Beer defines the shape of a Beer from the Untappd API
// Note: only fields used by slappd are included here
type Beer struct {
	ID          int     `json:"bid"`
	Name        string  `json:"beer_name"`
	Label       string  `json:"beer_label"`
	Ibu         int     `json:"beer_ibu"`
	Abv         float64 `json:"beer_abv"`
	Style       string  `json:"beer_style"`
	Description string  `json:"beer_description"`

	// Slug is part of a Beer for some Untappd API endpoints
	Slug string `json:"beer_slug,omitempty"`

	// Brewery is nested in Beer for some Untappd API endpoints
	Brewery *Brewery `json:"brewery,omitempty"`
}

// Brewery defines the shape of a Brewery from the Untappd API
// Note: only fields used by slappd are included here
type Brewery struct {
	Name string `json:"brewery_name"`
}

func (b *Beer) slug() string {
	if b.Slug != "" {
		return b.Slug
	}

	re := regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(b.Name), "-"), "-")
}

func title(beer *Beer, brewery *Brewery) string {
	u := fmt.Sprintf("https://untappd.com/b/%s/%d", beer.slug(), beer.ID)
	return fmt.Sprintf("<%s|%s>", u, fmt.Sprintf("%s %s", brewery.Name, beer.Name))
}

func text(beer *Beer) string {
	return fmt.Sprintf("%s | %d IBU | %0.0f%% ABV \n%s", beer.Style, beer.Ibu, beer.Abv, beer.Description)
}
