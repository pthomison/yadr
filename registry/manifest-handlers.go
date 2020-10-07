package registry

import(
	"net/http"
	"fmt"

	"io"

	"github.com/gorilla/mux"

)

func (r *Registry) ManifestPutHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("\n\n --- Manifest Put Handler")

		vars := mux.Vars(req)
	
		manifest, size, err := r.writeManifestTag(vars["image"], vars["reference"], io.Reader(req.Body))
		check(err)

		w.Header().Set("Content-Length", fmt.Sprintf("%v", size))
		w.Header().Set("Docker-Content-Digest", manifest)

	    w.Header().Set("Accept-Encoding", "gzip")

	    w.WriteHeader(http.StatusCreated)

	}
}

func (r *Registry) ManifestGetHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("\n\n --- Manifest Get Handler")

		vars := mux.Vars(req)
		image := vars["image"]
		reference := vars["reference"]

		var manifest []byte
		var err error
		var hash string

		fmt.Println("Sending Manifest")

		if isHash(reference, hashType) {
			fmt.Println("Hash Requested")
			hash = reference
		} else {
			fmt.Println("Tag Requested")
			hash, err = r.lookupTag(image, reference)
			check(err)
		}

		manifest, err = r.readManifest(image, hash)
		check(err)

		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		w.Header().Set("Docker-Content-Digest", hash)

		fmt.Println("Write Header")

	    w.WriteHeader(http.StatusOK)

	    w.Write(manifest)
	}
}