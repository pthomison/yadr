package registry

import(
	// "crypto/sha256"
	"net/http"
	"fmt"
	// "encoding/json"
	// "os"
	"io"
	"io/ioutil"

	"github.com/gorilla/mux"
	"strings"

)


// type manifest struct {
// 	SchemaVersion int `json:"schemaVersion"`
// 	MediaType string `json:"mediaType"`
// 	Config descriptor `json:"config"`
// 	Layers []layer `json:"layers"` 
// }

// type layer struct {
// 	descriptor
// }

// type descriptor struct {
// 	content []byte
// 	Size int64 `json:"size"`
// 	MediaType string `json:"mediaType"`
// 	digest
// }


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

func (r *Registry) readManifest(image string, hash string) ([]byte, error) {
	manifestLocation := r.StorageLocation + manifestFolder + image + "/index/" + hash

	content, err := ioutil.ReadFile(manifestLocation)
	if err != nil {
		return nil, err
	} else {
		return content, nil
	}
}

func (r *Registry) lookupTag(image string, tag string) (string, error) {
	tagLocation := r.StorageLocation + manifestFolder + image + "/tags/" + tag + "/link"
	content, err := ioutil.ReadFile(tagLocation)
	if err != nil {
		return "", err
	} else {
		return string(content), nil
	}
}

func isHash(reference string, hashType string) bool {
	if strings.Contains(reference, hashType) {
		return true
	} else {
		return false
	}
}