package registry

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"fmt"
)

type Registry struct {
	StorageLocation string
}

func New(storagePath string) (*Registry, error) {
	r := &Registry{
		StorageLocation: storagePath,
	}

	err := r.InitializeStorage()
	if err != nil {
		return nil, err
	}

	r.SetAPI()

	return r, nil
}



func (r *Registry) Serve() {
    http.ListenAndServe(":5000", nil)
}

func (r *Registry) SetAPI() {
    router := mux.NewRouter()
    router.HandleFunc(BaseAPI, BaseGetHandler).
    	Methods("GET")

    // MANIFEST HANDLERS
    router.HandleFunc(ManifestAPI, r.ManifestPutHandlerFactory()).
    	Methods("PUT")

    // BLOB HANDLERS
    router.HandleFunc(BlobAPI, r.BlobHeadHandlerFactory()).
    	Methods("HEAD")

    router.HandleFunc(BlobUploadRequestAPI, r.BlobUploadRequestPostHandlerFactory()).
    	Methods("POST")
    
    router.HandleFunc(BlobUploadAPI, r.BlobUploadPatchFactory()).
    	Methods("PATCH")

    router.HandleFunc(BlobUploadAPI, r.BlobUploadCompletePostFactory()).
    	Methods("PUT").
    	Queries("digest", "{digest}")


    router.PathPrefix("/").HandlerFunc(DefaultHandler)

    http.Handle("/", router)
}

func (r *Registry) InitializeStorage() error {
	err := os.MkdirAll(r.StorageLocation, storageFolderPerms)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.StorageLocation + blobFolder, storageFolderPerms)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.StorageLocation + manifestFolder, storageFolderPerms)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.StorageLocation + uploadFolder, storageFolderPerms)
	if err != nil {
		return err
	}

	return nil
} 

func check(e error) {
	if e != nil {
		fmt.Printf("%+v\n", e)
		panic(e)
	}
}