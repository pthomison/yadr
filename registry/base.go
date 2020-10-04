package registry

import(
	"net/http"
	"fmt"
)

func BaseGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Base Get Handler Called")
    w.WriteHeader(http.StatusOK)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler Called")
	fmt.Printf("%+v\n", r)

    w.WriteHeader(http.StatusNotFound)
}