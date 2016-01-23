package main

import (
	"net/http"
	"net/http/httputil"
	"text/template"
)

const (
	templateDir = "templates/"
)

type LookingGlass struct {
	proxy httputil.ReverseProxy
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

func (lg *LookingGlass) log(r *http.Request) {
	log.Sprintf("lookingglass: logged %s", r.URL)

}

func (lg *LookingGlass) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lg.log(w, r)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.Handle("/", &LookingGlass{})
	http.HandleFunc("/stats", StatsHandler)
	log.Fatalln(http.ListenAndServe(":8181", nil))
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	p := map[string]interface{}{
		"Data": "asdf",
	}

	tmpl, err := template.ParseFiles(templateDir + "index.html")
	tmpl.ExecuteTemplate(w, p)
}
