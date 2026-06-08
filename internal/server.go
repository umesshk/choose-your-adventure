package internal

import (
	"log"
	"net/http"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func ServerPage() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", ServeHTTP)

	log.Println("Server Running on PORT 3000")
	http.ListenAndServe(":3000", mux)

}
