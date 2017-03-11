package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kpurdon/slappd/internal/slack"
	"github.com/kpurdon/slappd/internal/untappd"
)

func isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st := os.Getenv("SLACK_TOKEN")
		if st == "" {
			log.Printf("missing environment variable: SLACK_TOKEN")
			http.Error(w, http.StatusText(500), 500)
			return
		}

		rt := r.FormValue("token")
		if rt == "" {
			log.Printf("missing form value: token")
			http.Error(w, http.StatusText(400), 400)
			return
		}

		var authorized bool
		for _, t := range strings.Split(st, ",") {
			if t == rt {
				authorized = true
			}
		}

		if !authorized {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		u := r.FormValue("user_name")
		if u == "" {
			log.Printf("missing form value: user_name")
			http.Error(w, http.StatusText(400), 400)
			return
		}

		if u == "slackbot" {
			w.WriteHeader(200)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	searchText := r.FormValue("text")
	if searchText == "" {
		log.Printf("missing form value: text")
		http.Error(w, http.StatusText(400), 400)
		return
	}

	ud, err := untappd.Search(searchText)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	sr := slack.NewResponse()
	for _, item := range ud.Response.Beers.Items {
		attachment := &slack.Attachment{
			Title:    item.Title(),
			Text:     item.Text(),
			ImageURL: item.Beer.Label,
		}
		sr.Attachments = append(sr.Attachments, attachment)

		// TODO: add in actions to select from the list of options
		// for now break after the first attachment so we only return
		// a single result
		break
	}

	b, err := json.Marshal(sr)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)

	log.Printf("slappd listening on 0.0.0.0%s", addr)

	http.Handle("/", isAuthorized(http.HandlerFunc(handler)))
	http.ListenAndServe(addr, nil)
}
