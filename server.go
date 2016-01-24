package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"text/template"
	"time"
)

const (
	_targetURL  = "https://www.cloudflare.com:443"
	templateDir = "templates/"

	// 500 requests in a 60 second window
	alertThreshold = 500
)

type RequestEvent struct {
	Req              *http.Request
	RecordedAt       time.Time
	UpstreamDuration time.Duration
}

type LookingGlass struct {
	proxy      *httputil.ReverseProxy
	requestLog []*RequestEvent
	mutex      *sync.RWMutex
}

func NewLookingGlass(target *url.URL) *LookingGlass {
	return &LookingGlass{
		requestLog: []*RequestEvent{},
		mutex:      &sync.RWMutex{},
		proxy: &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				// r.URL.Host = "cloudflare.com"
				// r.URL.Scheme = "http"
				// r.Host = "cloudflare.com"

				r.URL.Host = target.Host
				r.URL.Scheme = target.Scheme
				r.Host = target.Host
				// something here
			},
		},
	}
}

func (lg *LookingGlass) logEvent(e *RequestEvent) {
	lg.mutex.Lock()
	lg.requestLog = append(lg.requestLog, e)
	lg.mutex.Unlock()
}

func (lg *LookingGlass) Requests() []*RequestEvent {
	lg.mutex.RLock()
	reqs := lg.requestLog
	lg.mutex.RUnlock()

	return reqs
}

func (lg *LookingGlass) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := &RequestEvent{
		Req:        r,
		RecordedAt: time.Now(),
	}
	lg.logEvent(event)

	ts := time.Now()
	lg.proxy.ServeHTTP(w, r)
	event.UpstreamDuration = time.Now().Sub(ts)

	log.Print(fmt.Sprintf("lookingglass: logged %s\n", r.URL))
}

var lg *LookingGlass

func main() {
	target, err := url.Parse(_targetURL)
	if err != nil {
		panic(err)
	}
	lg = NewLookingGlass(target)

	http.Handle("/", lg)
	http.HandleFunc("/stats", StatsHandler)
	log.Fatalln(http.ListenAndServe(":8182", nil))
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	reqs := lg.Requests()

	tf := &TimeFilter{from: time.Now().Add(-time.Hour), to: time.Now()}
	reqs = tf.Filter(reqs)

	rg := &RouteGrouper{}
	gr := rg.Group(reqs)

	p := map[string]interface{}{
		"Requests":             reqs,
		"RouteGroups":          gr,
		"ResponseTimeAnalysis": analyzeResponseTimes(reqs),
	}

	tmpl, _ := template.ParseFiles(templateDir + "index.html")
	tmpl.Execute(w, p)
}
