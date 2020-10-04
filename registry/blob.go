package registry

import(
	"os"
	"io"
	"net/http"
	"fmt"
	// "time"
	"github.com/gorilla/mux"

	"github.com/google/uuid"

	"path/filepath"

)

// type digest struct {
// 	hashType string
// 	verified bool
// 	Digest string `json:"digest"`
// }

// type blob struct {

// }

func (r *Registry) BlobUploadRequestPostHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Blob Upload Request Post Handler")

		id := uuid.New().String()

		vars := mux.Vars(req)
	 

	    w.Header().Set("Accept-Encoding", "gzip")
	    w.Header().Set("Content-Length", "0")
	    w.Header().Set("Range", "0-0")

	    w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", vars["image"], id))
	    w.Header().Set("Docker-Upload-UUID", id)
	    w.WriteHeader(http.StatusAccepted)
	}
}

func (r *Registry) BlobUploadPatchFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Blob Upload Patch Handler")

		id := filepath.Base(req.URL.Path)

		uploadFile := r.StorageLocation + "/uploads/" + id

		vars := mux.Vars(req)

		f, err := os.Create(uploadFile)
		check(err)


		c, err := io.Copy(f, req.Body)
		check(err)

		rng := fmt.Sprintf("0-%v", c)
	
	    w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", vars["image"], id))
	    w.Header().Set("Range", rng)
		w.WriteHeader(http.StatusAccepted)
	}
}


func (r *Registry) BlobUploadCompletePostFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Blob Upload Complete Post Handler")

		id := filepath.Base(req.URL.Path)

		uploadFile := r.StorageLocation + "/uploads/" + id

		vars := mux.Vars(req)
		
		fmt.Println(vars)

		hash, err := hashFile(uploadFile)
		check(err)

		if hash == vars["digest"] {
			check(r.moveBlobUpload(id, hash))
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}



func (r *Registry) BlobHeadHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Blob Head Handler Called")

		vars := mux.Vars(req)
		fmt.Printf("%+v\n", vars)

		fmt.Printf("\nCheck for blobs:%+v\n", r.checkForBlob(vars["digest"]))

		if r.checkForBlob(vars["digest"]) {
			w.WriteHeader(http.StatusOK)
		} else {
	    	w.WriteHeader(http.StatusNotFound)
		}
	}
}