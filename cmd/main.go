package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/umesshk/choose-your-adventure/internal"
)

var tpl = template.Must(template.ParseFiles("web/index.html"))

type HandlerOption func(*handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPath(fn func(*http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(story internal.Story, opts ...HandlerOption) http.Handler {
	h := handler{story, tpl, DefaultParseUrl}

	for _, opt := range opts {
		opt(&h)
	}

	return h

}

func pathFn(r *http.Request) string {

	url_path := r.URL.Path

	url_path = strings.TrimSpace(url_path)

	if url_path == "/story" || url_path == "/story/" {
		return "intro"
	}

	return strings.TrimPrefix(url_path, "/story/")

}

type handler struct {
	story  internal.Story
	t      *template.Template
	pathFn func(*http.Request) string
}

func DefaultParseUrl(r *http.Request) string {

	url_path := r.URL.Path

	url_path = strings.Trim(url_path, "/")

	if strings.HasPrefix(url_path, "css") || strings.HasPrefix(url_path, "favicon.ico") {
		return ""
	}

	return url_path

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	story_arc := h.pathFn(r)
	var err error
	if story_arc == "" {

		err = tpl.Execute(w, h.story["intro"])
	} else {

		if chapter, ok := h.story[story_arc]; ok {
			err = tpl.Execute(w, chapter)
		}
	}
	if err != nil {
		fmt.Println("Error Occured ", err)
		log.Fatal(err)
	}

}

func main() {

	file_name := flag.String("json", "gopher.json", "Provide any valid JSON file containing CYOA story ")
	port := flag.Int("port", 3000, "Provide port number for the application")
	flag.Parse()

	json_file, err := os.Open(*file_name)

	if err != nil {
		fmt.Println("Error Opening File ")
		os.Exit(1)
	}

	story, err := internal.MakeStory(json_file)

	if err != nil {
		fmt.Println("Error Occured : ", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	handler := NewHandler(story)

	h := NewHandler(story, WithPath(pathFn))

	mux.HandleFunc("/story", h.ServeHTTP)
	mux.HandleFunc("/", handler.ServeHTTP)

	log.Printf("Server Running on PORT %d", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)

}
