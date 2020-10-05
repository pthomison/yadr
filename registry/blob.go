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
		fmt.Println("\n\n --- Blob Upload Request Post Handler")

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
		fmt.Println("\n\n --- Blob Upload Patch Handler")

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
		fmt.Println("\n\n --- Blob Upload Complete Post Handler")

		id := filepath.Base(req.URL.Path)

		uploadFile := r.StorageLocation + "/uploads/" + id

		vars := mux.Vars(req)
		
		fmt.Println(vars)

		hash, err := hashFile(uploadFile)
		check(err)

	    w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/%s", vars["image"], id))

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
		fmt.Println("\n\n --- Blob Head Handler Called")

		vars := mux.Vars(req)
		fmt.Printf("%+v\n", vars)

		// fmt.Printf("\nCheck for blobs:%+v\n", r.checkForBlob(vars["digest"]))

		exists, size := r.checkForBlob(vars["digest"])

		if exists {
			w.Header().Set("Content-Length", fmt.Sprintf("%v", size))
			w.Header().Set("Docker-Content-Digest", vars["digest"])

			w.WriteHeader(http.StatusOK)
		} else {
	    	w.WriteHeader(http.StatusNotFound)
		}

		fmt.Printf("\n%+v\n\n", w)
	}
}

func (r *Registry) BlobGetHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("\n\n --- Blob Get Handler Called")

		vars := mux.Vars(req)
		fmt.Printf("%+v\n", vars)

		exists, size := r.checkForBlob(vars["digest"])

		if exists {

			blobLocation := r.StorageLocation + blobFolder + vars["digest"]	

			w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
			w.Header().Set("Content-Length", fmt.Sprintf("%v", size))
			w.Header().Set("Docker-Content-Digest", vars["digest"])

			w.WriteHeader(http.StatusOK)

			f, err := os.Open(blobLocation)
			check(err)
			defer f.Close()


			_, err = io.Copy(w, f)
			check(err)

		} else {
	    	w.WriteHeader(http.StatusNotFound)
		}

		fmt.Printf("\n%+v\n\n", w)
	}
}