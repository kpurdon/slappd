package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/kpurdon/slappd/internal/slack"
	"github.com/kpurdon/slappd/internal/untappd"
)

const callbackID = "slappd"

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

func selectHandler(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")

	var action slack.ActionPayload
	if err := json.Unmarshal([]byte(payload), &action); err != nil {
		log.Printf("invalid action json: %v", err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if action.CallbackID != callbackID {
		log.Printf("invalid callback_id: %s", action.CallbackID)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	id, err := untappd.Info(action.Actions[0].Value)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	sr := slack.NewResponse()
	attachment := &slack.Attachment{
		Title:    id.Title(),
		Text:     id.Text(),
		ImageURL: id.Response.Beer.Label,
	}
	sr.Attachments = append(sr.Attachments, attachment)

	b, err := json.Marshal(sr)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
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

	var sr *slack.Response
	if len(ud.Response.Beers.Items) != 0 {
		var count int
		sr = slack.NewResponse()
		for _, item := range ud.Response.Beers.Items {
			count++
			attachment := &slack.Attachment{
				Title: item.Title(),
				// Text:     item.Text(),
				// ImageURL: item.Beer.Label,
				CallbackID: callbackID,
				Actions:    []*slack.Action{slack.NewAction(item.Beer.ID)},
			}
			sr.Attachments = append(sr.Attachments, attachment)

			// TODO: add in actions to select from the list of options
			// for now break after the first attachment so we only return
			// a single result
			if count == 5 {
				break
			}
		}
	} else {
		sr = slack.NewEmptyResultsResponse()
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
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)

	http.Handle("/", isAuthorized(http.HandlerFunc(searchHandler)))
	http.Handle("/select", http.HandlerFunc(selectHandler))

	server := &http.Server{Addr: addr}

	go func() {
		log.Printf("slappd listening on 0.0.0.0%s", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stopChan

	log.Printf("slappd shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
