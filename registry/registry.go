package registry

import (
	"github.com/gorilla/mux"
	"net/http"
)

type registry struct {
	StorageLocation string

}

func (r *registry) Serve() {
    r := mux.NewRouter()
    r.HandleFunc(BaseAPI, BaseHandler)


    r.Methods("GET").HandleFunc(BaseAPI + ManifestAPI, ManifestGetHandler)


    r.Methods("GET").HandleFunc(BaseAPI + BlobAPI, BlobGetHandler)

    http.Handle("/", r)
    http.ListenAndServe(":5000", nil)
}
