package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"

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

type RequestObj struct {
	Headers  CHeader     `json:"headers" xml:"-"`
	URL      url.URL     `json:"url" xml:"-"`
	Body     interface{} `json:"body" xml:"-"`
	Host     string      `json:"host" xml:"host"`
	Protocol string      `json:"proto" xml:"proto"`
	Method   string      `json:"method" xml:"method"`
	Form     url.Values  `json:"form" xml:"-"`
}

type ResponseMsg struct {
	Hostname  string      `json:"host" xml:"host"`
	Timestamp string      `json:"ts" xml:"ts"`
	Request   *RequestObj `json:"request" xml:"request"`
	Errors    []string    `json:"errors" xml:"errors"`
}

type CHeader struct {
	*http.Header
}

func (h *CHeader) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return nil
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Unable to collect hostname: %s\n", err.Error())
		hostname = "unknown"
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" && r.Method == http.MethodGet {
			promhttp.Handler().ServeHTTP(w, r)
			return
		}

		errors := []string{}

		headers := CHeader{
			&r.Header,
		}
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
					log.Errorf("Encountered parsing JSON request body: %s\n", err.Error())
					errors = append(errors, fmt.Sprintf("Encountered parsing JSON request body: %s", err.Error()))
					strBody = string(buf)
				}
			} else if strings.HasPrefix(r.Header.Get("Content-Type"), "application/xml") && !reqTooLarge {
				err = xml.Unmarshal(buf, &jsonBody)
				if err != nil {
					log.Errorf("Encountered parsing XML request body: %s\n", err.Error())
					errors = append(errors, fmt.Sprintf("Encountered parsing XML request body: %s", err.Error()))
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

		rMsg := RequestObj{
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

		msg := ResponseMsg{
			Hostname:  hostname,
			Timestamp: time.Now().UTC().String(),
			Request:   &rMsg,
			Errors:    errors,
		}
		w.Header().Add("Server", "echosrv@latest")

		if strings.HasPrefix(r.Header.Get("Accept"), "application/xml") {
			rxml, err := xml.MarshalIndent(msg, "", "\t")
			if err != nil {
				log.Errorf("Encountered outputting XML request body: %s\n", err.Error())
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("<error>Error encountered while processing your request</error>"))
				return
			}
			w.Header().Add("Content-Type", "application/xml")
			fmt.Fprintf(w, "%s\n", string(rxml))
		} else {
			rjson, err := json.MarshalIndent(msg, "", "\t")
			if err != nil {
				log.Errorf("Encountered outputting JSON response: %s\n", err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{\"error\":\"Error encountered while processing your request\"}"))
				return
			}

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, "%s\n", string(rjson))
		}
	}

	bind := ":8889"
	s := &http.Server{
		Addr:           bind,
		Handler:        handlers.CombinedLoggingHandler(os.Stdout, handlers.RecoveryHandler()(http.HandlerFunc(h))),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Infof("🚀 Listening for traffic on: %s\n", bind)
	log.Fatal(s.ListenAndServe())
}
