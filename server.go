package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"text/template"
)

const (
	templateDir = "templates/"
)

type LookingGlass struct {
	proxy *httputil.ReverseProxy
}

func NewLookingGlass() *LookingGlass {
	return &LookingGlass{
		proxy: &httputil.ReverseProxy{
			Director: func(r *http.Request) {
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
	lg := NewLookingGlass()

	http.Handle("/", lg)
	http.HandleFunc("/stats", StatsHandler)
	log.Fatalln(http.ListenAndServe(":8181", nil))
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	p := map[string]interface{}{
		"Data": "asdf",
	}

	tmpl, _ := template.ParseFiles(templateDir + "index.html")
	tmpl.Execute(w, p)
}
