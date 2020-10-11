package registry

import(
	"net/http"
	"fmt"
	"github.com/gorilla/mux"

	"github.com/google/uuid"

	"path/filepath"


    "github.com/sirupsen/logrus"
)

func (r *Registry) HeadBlob(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Blob Head Handler Called")

	vars := mux.Vars(req)
	digest := vars["digest"]

	b, exists := r.blobs[digest]

	if exists {
		w.Header().Set("Content-Length", fmt.Sprintf("%v", b.contentLength))
		w.Header().Set("Docker-Content-Digest", b.digest)

		w.WriteHeader(http.StatusOK)
	} else {
    	w.WriteHeader(http.StatusNotFound)
	}
}

func (r *Registry) GetBlob(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Blob Get Handler Called")


	vars := mux.Vars(req)
	digest := vars["digest"]


	b, exists := r.blobs[digest]

	if exists {
		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		w.Header().Set("Content-Length", fmt.Sprintf("%v", b.contentLength))
		w.Header().Set("Docker-Content-Digest", b.digest)

		w.WriteHeader(http.StatusOK)

		err := b.SendData(w)
		check(err)


	} else {
		logrus.Debug("Blob Not Found")

    	w.WriteHeader(http.StatusNotFound)
	}
}

func (r *Registry) DeleteBlob(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Blob Delete Handler Called")

	vars := mux.Vars(req)
	digest := vars["digest"]

	b, exists := r.blobs[digest]

	if exists {
		err := r.Delete(b)
		check(err)

		delete(r.blobs, b.digest)

		w.WriteHeader(http.StatusAccepted)

	} else {
    	w.WriteHeader(http.StatusNotFound)
	}
}

// end of all upload types
func (r *Registry) CompleteBlobUpload(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Complete Blob Upload")

	vars := mux.Vars(req)
	image := vars["image"]
	var digest string 
	var id string

	if val, ok := vars["digest"]; ok {
	    digest = val
	} else {
		id = uuid.New().String()
		w.WriteHeader(http.StatusBadRequest)
	}


	if val, ok := vars["sessionID"]; ok {
	    id = val
	} else {
		id = uuid.New().String()
	}

	u := r.BlobUploadInit(image, id)

	if req.Header.Get("Content-Length") != "0" {
		err := u.StoreUploadData(req.Body)
		check(err)
	} 


	b, err := r.ProcessUpload(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		check(err)
		return
	}

	if b.digest != digest {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/%s", image, b.digest))
    w.Header().Set("Range", fmt.Sprintf("0-%v", b.contentLength))
	w.WriteHeader(http.StatusCreated)

	logrus.Info("Blob Uploaded: ", b.digest)
}


// always accept
func (r *Registry) RequestBlobUpload(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Blob Upload Request")

	vars := mux.Vars(req)
	image := vars["image"]
	id := uuid.New().String()

	// start of upload, more requests coming
	if req.Header.Get("Content-Length") == "0" {
		w.Header().Set("Content-Length", "0")
		w.Header().Set("Range", "0-0")

		w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", image, id))
		w.Header().Set("Docker-Upload-UUID", id)

		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(UserErrorResponse)
	}
}

// chunked upload
func (r *Registry) PatchBlobUpload(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Blob Upload Patch Handler")

	vars := mux.Vars(req)
	image := vars["image"]

	id := filepath.Base(req.URL.Path)

	u := r.BlobUploadInit(image, id)

	if req.Header.Get("Content-Length") != "0" {
		u.StoreUploadData(req.Body)
	} else {

	}

	err := u.StoreUploadData(req.Body)
	check(err)

    w.Header().Set("Location", u.url)
    w.Header().Set("Range", fmt.Sprintf("0-%v", u.contentLength))
	w.WriteHeader(http.StatusAccepted)
}