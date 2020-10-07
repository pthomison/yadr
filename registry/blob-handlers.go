package registry

import(
	"net/http"
	"fmt"
	"github.com/gorilla/mux"

	"github.com/google/uuid"

	"path/filepath"

)

func (r *Registry) RequireBlobDigest(fn func(w http.ResponseWriter, req *http.Request, b *Blob)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		digest := vars["digest"]
		image := vars["image"]

		blob, err := r.BlobDigestInit(digest, image)
		check(err)

		fn(w, req, blob)
	}
}


func (r *Registry) HeadBlob(w http.ResponseWriter, req *http.Request, b *Blob) {
	fmt.Println("\n\n --- Blob Head Handler Called")

	if b.exists {
		w.Header().Set("Content-Length", fmt.Sprintf("%v", b.contentLength))
		w.Header().Set("Docker-Content-Digest", b.digest)

		w.WriteHeader(http.StatusOK)
	} else {
    	w.WriteHeader(http.StatusNotFound)
	}
}


func (r *Registry) GetBlob(w http.ResponseWriter, req *http.Request, b *Blob) {
	fmt.Println("\n\n --- Blob Get Handler Called")

	if b.exists {
		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		w.Header().Set("Content-Length", fmt.Sprintf("%v", b.contentLength))
		w.Header().Set("Docker-Content-Digest", b.digest)

		w.WriteHeader(http.StatusOK)

		err := b.SendData(w)
		check(err)


	} else {
    	w.WriteHeader(http.StatusNotFound)
	}
}


// always accept
func (r *Registry) PostBlobUploadRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n --- Blob Upload Request Post Handler")

	id := uuid.New().String()

	vars := mux.Vars(req)
	image := vars["image"]

	w.Header().Set("Accept-Encoding", "gzip")
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Range", "0-0")

	w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", image, id))
	w.Header().Set("Docker-Upload-UUID", id)
	w.WriteHeader(http.StatusAccepted)
}


func (r *Registry) PatchBlobUpload(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n ----- Blob Upload Patch Handler")

	id := filepath.Base(req.URL.Path)

	vars := mux.Vars(req)
	image := vars["image"]

	b, err := r.BlobUploadInit(image, id)
	check(err)

	err = b.StoreUploadData(req.Body)
	check(err)

	fmt.Printf("%+v\n", b)


    w.Header().Set("Location", b.url)
    w.Header().Set("Range", fmt.Sprintf("0-%v", b.contentLength))
	w.WriteHeader(http.StatusAccepted)
}


func (r *Registry) PostBlobUploadComplete(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n --- Blob Upload Complete Post Handler")

	id := filepath.Base(req.URL.Path)

	vars := mux.Vars(req)
	image := vars["image"]

	u, err := r.BlobUploadInit(image, id)
	check(err)

	b, err := r.ProcessUpload(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Location", b.urlLocation)
	w.WriteHeader(http.StatusCreated)
}