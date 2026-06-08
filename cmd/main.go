package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/umesshk/choose-your-adventure/internal"
)

func NewHandler(story internal.Story) http.Handler {
	return handler{
		story: story,
	}
}

type handler struct {
	story internal.Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("web/index.html"))

	err := tpl.Execute(w, h.story["intro"])

	if err != nil {
		fmt.Println("Error Occured ", err)
		log.Fatal(err)
	}

}

func main() {

	file_name := flag.String("json", "gopher.json", "Provide any valid JSON file containing CYOA story ")
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
	log.Println("Server Running on PORT 3000")
	http.ListenAndServe(":3000", mux)

}
