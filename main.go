package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("[warn] Unable to collect hostname: %s\n", err.Error())
		hostname = "unknown"
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		errors := []string{}

		headers := r.Header
		u := r.URL

		var body string
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			body = string(buf)
		}

		host := r.Host
		proto := r.Proto
		method := r.Method

		r.ParseForm()
		form := r.Form

		rMsg := struct {
			Headers  http.Header `json:"headers"`
			URL      url.URL     `json:"url"`
			Body     string      `json:"body"`
			Host     string      `json:"host"`
			Protocol string      `json:"proto"`
			Method   string      `json:"method"`
			Form     url.Values  `json:"form"`
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
			Hostname  string      `json:"host"`
			Timestamp time.Time   `json:"ts"`
			Request   interface{} `json:"request"`
			Errors    []string    `json:"errors"`
		}{
			Hostname:  hostname,
			Timestamp: time.Now().UTC(),
			Request:   rMsg,
			Errors:    errors,
		}

		rjson, err := json.MarshalIndent(msg, "", "\t")
		if err != nil {
			return
		}
		fmt.Printf("[info] request received\n%s\n", string(rjson))

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", "echosrv@latest")
		fmt.Fprintf(w, "%s\n", string(rjson))
	}

	bind := ":8889"
	s := &http.Server{
		Addr:           bind,
		Handler:        http.HandlerFunc(h),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("[info] 🚀 Listening for traffic on: %s\n", bind)
	log.Fatal(s.ListenAndServe())
}
