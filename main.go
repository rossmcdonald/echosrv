package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	h := func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header
		u := r.URL
		body := r.Body
		host := r.Host
		proto := r.Proto
		method := r.Method

		r.ParseForm()
		form := r.Form

		rMsg := struct {
			Headers  http.Header   `json:"headers"`
			URL      url.URL       `json:"url"`
			Body     io.ReadCloser `json:"body"`
			Host     string        `json:"host"`
			Protocol string        `json:"proto"`
			Method   string        `json:"method"`
			Form     url.Values    `json:"form"`
		}{
			Headers:  headers,
			URL:      *u,
			Body:     body,
			Host:     host,
			Protocol: proto,
			Method:   method,
			Form:     form,
		}

		msg := struct {
			Timestamp time.Time   `json:"ts"`
			Request   interface{} `json:"request"`
		}{
			Timestamp: time.Now().UTC(),
			Request:   rMsg,
		}

		rjson, err := json.MarshalIndent(msg, "", "\t")
		if err != nil {
			return
		}
		fmt.Printf("%s\n", string(rjson))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", string(rjson))
	}

	s := &http.Server{
		Addr:           ":8888",
		Handler:        http.HandlerFunc(h),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
