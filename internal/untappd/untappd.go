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

	// only available from the BeerInfo API
	Slug        string   `json:"beer_slug,omitempty"`
	RatingCount int      `json:"rating_count,omitempty"`
	RatingScore float64  `json:"rating_score,omitempty"`
	Brewery     *Brewery `json:"brewery,omitempty"`
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
	format := "%s | %d IBU | %0.0f%% ABV | %0.02f rating (%d votes)\n%s"
	return fmt.Sprintf(format, beer.Style, beer.Ibu, beer.Abv, beer.RatingScore, beer.RatingCount, beer.Description)
}
