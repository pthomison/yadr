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

type digest struct {
	hashType string
	verified bool
	Digest string `json:"digest"`
}

type blob struct {

}

func (r *Registry) BlobUploadRequestPostHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Blob Upload Request Post Handler")

		id := uuid.New().String()

		fmt.Printf("%+v\n", id)


		vars := mux.Vars(req)
	 

	    w.Header().Set("Accept-Encoding", "gzip")
	    w.Header().Set("Content-Length", "0")
	    w.Header().Set("Range", "0-0")

	    w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", vars["image"], id))
	    w.Header().Set("Docker-Upload-UUID", id)
	    w.WriteHeader(http.StatusAccepted)

	    fmt.Printf("%+v\n", w.Header())
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

		fmt.Printf("%+v\n", rng)

	
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

		fmt.Println(hash)

		hashCheck := (hash == vars["digest"])

		fmt.Println(hashCheck)


		check(r.moveBlobUpload(id, hash))	

		w.WriteHeader(http.StatusCreated)
	}
}




func BlobGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Blob Get Handler Called")

    w.WriteHeader(http.StatusOK)
}

func BlobPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Blob Post Handler Called")


    w.WriteHeader(http.StatusOK)
}


func (r *Registry) BlobHeadHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Blob Head Handler Called")

		vars := mux.Vars(req)
		fmt.Printf("%+v\n", vars)

		r.checkForBlob(vars["digest"])


	    w.WriteHeader(http.StatusNotFound)
	}
}

func (r *Registry) BlobPostHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Blob Post Handler Called")

		vars := mux.Vars(req)
		fmt.Printf("%+v\n", vars)



		// r.checkForBlob(vars["digest"])


	    w.WriteHeader(http.StatusNotFound)
	}
}



func check(e error) {
	if e != nil {
		fmt.Printf("%+v\n", e)
		panic(e)
	}
}