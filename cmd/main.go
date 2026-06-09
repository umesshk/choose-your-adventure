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

func NewHandler(story internal.Story, opts ...HandlerOption) http.Handler {
	h := handler{story, tpl}

	for _, opt := range opts {
		opt(&h)
	}

	return h

}

type handler struct {
	story internal.Story
	t     *template.Template
}

func (h handler) ParseUrl(url string) string {

	url = strings.Trim(url, "/")

	if strings.HasPrefix(url, "css") || strings.HasPrefix(url, "favicon.ico") {
		return ""
	}

	return url

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	url_path := r.URL.Path

	story_arc := h.ParseUrl(url_path)
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

	mux.HandleFunc("/", handler.ServeHTTP)

	log.Printf("Server Running on PORT %d", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)

}
