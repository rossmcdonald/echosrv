package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "echo_request_count",
		Help: "The total number of received requests",
	})
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("[warn] Unable to collect hostname: %s\n", err.Error())
		hostname = "unknown"
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" && r.Method == http.MethodGet {
			promhttp.Handler().ServeHTTP(w, r)
			return
		}

		errors := []string{}

		headers := r.Header
		u := r.URL

		requestCount.Inc()

		var jsonBody interface{}
		var strBody string
		buf, err := ioutil.ReadAll(io.LimitReader(r.Body, 10000))
		if err != nil {
			errors = append(errors, fmt.Sprintf("Encountered error ready request body: %s", err.Error()))
		} else {
			reqTooLarge := false
			if len(buf) > 10000 {
				errors = append(errors, fmt.Sprintf("Request size %d too large to process", len(buf)))
				reqTooLarge = true
			}

			if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") && !reqTooLarge {
				err = json.Unmarshal(buf, &jsonBody)
				if err != nil {
					errors = append(errors, fmt.Sprintf("Encountered parsing JSON request body: %s", err.Error()))
					strBody = string(buf)
				}
			} else if strings.HasPrefix(r.Header.Get("Content-Type"), "application/octet-stream") {
				strBody = "<some binary data>"
			} else {
				if reqTooLarge {
					strBody = string(buf[:10000])
				} else {
					strBody = string(buf)
				}
			}

		}

		host := r.Host
		proto := r.Proto
		method := r.Method

		r.ParseForm()
		form := r.Form

		rMsg := struct {
			Headers  http.Header `json:"headers"`
			URL      url.URL     `json:"url"`
			Body     interface{} `json:"body"`
			Host     string      `json:"host"`
			Protocol string      `json:"proto"`
			Method   string      `json:"method"`
			Form     url.Values  `json:"form"`
		}{
			Headers:  headers,
			URL:      *u,
			Body:     strBody,
			Host:     host,
			Protocol: proto,
			Method:   method,
			Form:     form,
		}

		if jsonBody != nil {
			rMsg.Body = jsonBody
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
