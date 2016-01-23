package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"text/template"
)

const (
	_targetURL  = "http://www.cloudflare.com:443"
	templateDir = "templates/"
)

type LookingGlass struct {
	proxy *httputil.ReverseProxy
}

func NewLookingGlass(target *url.URL) *LookingGlass {
	return &LookingGlass{
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

func (lg *LookingGlass) process(r *http.Request) {
	log.Print(fmt.Sprintf("lookingglass: logged %s\n", r.URL))
}

func (lg *LookingGlass) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lg.process(r)
	lg.proxy.ServeHTTP(w, r)
}

func main() {
	target, err := url.Parse(_targetURL)
	if err != nil {
		panic(err)
	}
	lg := NewLookingGlass(target)

	http.Handle("/", lg)
	http.HandleFunc("/stats", StatsHandler)
	log.Fatalln(http.ListenAndServe(":8182", nil))
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	p := map[string]interface{}{
		"Data": "asdf",
	}

	tmpl, _ := template.ParseFiles(templateDir + "index.html")
	tmpl.Execute(w, p)
}
