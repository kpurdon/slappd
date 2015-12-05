package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type SlackAttachment struct {
	Title    string `json:"title,omitempty"`
	Text     string `json:"test,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type SlackResponse struct {
	ResponseType string           `json:"response_type,omitempty"`
	Text         string           `json:"text"`
	Attachment   *SlackAttachment `json:"attachments,omitempty"`
}

type UntappdBeer struct {
	Name        string  `json:"beer_name"`
	Label       string  `json:"beer_label"`
	Ibu         float64 `json:"beer_ibu"`
	Abv         float64 `json:"beer_abv"`
	Description string  `json:"beer_description"`
}

type UntappdBeerResponse struct {
	Beers struct {
		Items []struct {
			Beer *UntappdBeer `json:"beer"`
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

	requestURL.Path = "/v4/search/beer"
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

func handler(w http.ResponseWriter, r *http.Request) {

	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Println("unable to read environment variable SLACK_TOKEN")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	reqSlackToken := r.FormValue("token")
	if reqSlackToken == "" {
		http.Error(w, "MISSING_ARG_TOKEN", http.StatusBadRequest)
		return
	}

	if slackToken != r.FormValue("token") {
		http.Error(w, "", http.StatusUnauthorized)
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

	// set the default response data
	response := &SlackResponse{
		ResponseType: "in_channel",
		Text:         "Your Untappd Response",
		Attachment:   &SlackAttachment{},
	}

	if len(untappdData.Beer.Beers.Items) == 0 {
		response.ResponseType = "ephemeral"
		response.ResponseType = "No Results Found"
	} else {
		beer := untappdData.Beer.Beers.Items[0].Beer
		response.Attachment.Title = beer.Name
		response.Attachment.Text = beer.Description
		response.Attachment.ImageURL = beer.Label
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
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
