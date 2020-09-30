package registry

import (
	"fmt"
	// "os"
	"github.com/gorilla/mux"
	"net/http"
)

func Serve() {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    // r.HandleFunc("/products", ProductsHandler)
    // r.HandleFunc("/articles", ArticlesHandler)
    http.Handle("/", r)
    http.ListenAndServe(":5000", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // vars := mux.Vars(r)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Category: %v\n", "tod")
}