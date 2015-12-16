package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type SlackAttachment struct {
	Title    string `json:"title,omitempty"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type SlackResponse struct {
	ResponseType string             `json:"response_type,omitempty"`
	Text         string             `json:"text"`
	Attachments  *[]SlackAttachment `json:"attachments,omitempty"`
}

type UntappdBeer struct {
	ID          int     `json:"bid"`
	Name        string  `json:"beer_name"`
	Label       string  `json:"beer_label"`
	Ibu         int     `json:"beer_ibu"`
	Abv         float64 `json:"beer_abv"`
	Style       string  `json:"beer_style"`
	Description string  `json:"beer_description"`
}

type UntappdBrewery struct {
	Name string `json:"brewery_name"`
}

type UntappdBeerResponse struct {
	Beers struct {
		Items []struct {
			Beer    *UntappdBeer    `json:"beer"`
			Brewery *UntappdBrewery `json:"brewery"`
		}
	}
}

type UntappdResponse struct {
	Meta struct {
		StatusCode int `json:"code"`
	}
	Beer *UntappdBeerResponse `json:"response"`
}

func untappdRequest(searchString string) (untappdData *UntappdResponse, err error) {

	untappdResponse := &UntappdResponse{}

	untappdClientID := os.Getenv("UNTAPPD_CLIENT_ID")
	if untappdClientID == "" {
		return untappdResponse, errors.New("unable to read environment variable UNTAPPD_CLIENT_ID")
	}
	untappdClientSecret := os.Getenv("UNTAPPD_CLIENT_SECRET")
	if untappdClientSecret == "" {
		return untappdResponse, errors.New("unable to read environment variable UNTAPPD_CLIENT_SECRET")
	}

	requestURL := url.URL{
		Scheme: "https",
		Host:   "api.untappd.com",
		Path:   "/v4/search/beer",
	}
	q := requestURL.Query()
	q.Set("client_id", untappdClientID)
	q.Set("client_secret", untappdClientSecret)
	q.Set("q", searchString)
	requestURL.RawQuery = q.Encode()

	res, err := http.Get(requestURL.String())
	if err != nil {
		return untappdResponse, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return untappdResponse, err
	}
	if err := json.Unmarshal(body, &untappdResponse); err != nil {
		return untappdResponse, err
	}

	return untappdResponse, nil

}

func slug(s string) string {
	var re = regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func handler(w http.ResponseWriter, r *http.Request) {

	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Println("unable to read environment variable SLACK_TOKEN")
		http.Error(w, "UNABLE_TO_READ_SLACK_TOKEN", http.StatusInternalServerError)
		return
	}

	reqSlackToken := r.FormValue("token")
	if reqSlackToken == "" {
		http.Error(w, "MISSING_ARG_TOKEN", http.StatusBadRequest)
		return
	}

	validToken := false
	slackTokens := strings.Split(slackToken, ",")
	for _, token := range slackTokens {
		if token == reqSlackToken {
			validToken = true
			break
		}
	}

	if !validToken {
		http.Error(w, "INVALID_TOKEN", http.StatusUnauthorized)
		return
	}

	reqUserName := r.FormValue("user_name")
	if reqUserName == "" {
		http.Error(w, "MISSING_ARG_USER_NAME", http.StatusBadRequest)
		return
	}
	if reqUserName == "slackbot" {
		return // this is a bot response that we simply want to ignore
	}

	searchString := r.FormValue("text")
	if searchString == "" {
		http.Error(w, "MISSING_ARG_TEXT", http.StatusBadRequest)
		return
	}

	untappdData, err := untappdRequest(searchString)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	response := &SlackResponse{
		ResponseType: "in_channel",
		Text:         "Your Untappd Response",
	}

	if len(untappdData.Beer.Beers.Items) == 0 {
		response.ResponseType = "ephemeral"
		response.Text = "No Results Found"
	} else {
		beer := untappdData.Beer.Beers.Items[0].Beer
		brewery := untappdData.Beer.Beers.Items[0].Brewery

		untappdURL := fmt.Sprintf("https://untappd.com/b/%s/%d", slug(beer.Name), beer.ID)

		newAttachment := SlackAttachment{
			Title:    fmt.Sprintf("<%s|%s>", untappdURL, fmt.Sprintf("%s %s", brewery.Name, beer.Name)),
			Text:     fmt.Sprintf("%s | %d IBU | %0.0f%% ABV \n%s", beer.Style, beer.Ibu, beer.Abv, beer.Description),
			ImageURL: beer.Label,
		}
		attachments := append([]SlackAttachment{}, newAttachment)
		response.Attachments = &attachments
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func main() {
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
