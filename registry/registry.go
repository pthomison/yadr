package registry

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Registry struct {
	StorageLocation string

}

func (r *Registry) Serve() {
    router := mux.NewRouter()
    router.HandleFunc(BaseAPI, BaseHandler)


    router.HandleFunc(BaseAPI + ManifestAPI, ManifestGetHandler).Methods("GET")


    router.HandleFunc(BaseAPI + BlobAPI, BlobGetHandler).Methods("GET")

    http.Handle("/", router)
    http.ListenAndServe(":5000", nil)
}
