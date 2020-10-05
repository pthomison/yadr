package registry

import(
	"os"
	"io"
	// "io/ioutil"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"

	"github.com/google/uuid"

	"path/filepath"

)

type Blob struct {
	digest string
	contentLength int64
	contentType string
	fileLocation string
	uploadUrl string
	uploadFileLocation string
	urlLocation string
	exists bool
	image string
	uploadID string
}



type Digest string


// takes in a digest, return a blob with all fields besides content, contentType, & upload location loaded
func (r *Registry) BlobDigestInit(digest string, image string) (*Blob, error) {
	b := &Blob{
		digest: digest,
		image: image,
		fileLocation: fmt.Sprintf("%v%v%v", r.StorageLocation, blobFolder, digest),
		urlLocation: fmt.Sprintf("/v2/%s/blobs/%s", image, digest),
	}

	exists, _ := r.checkForBlob(digest)

	b.exists = exists


	if b.exists {
		f, err := os.Stat(b.fileLocation)
		if err != nil {
			return nil, err
		}

		b.contentLength = f.Size()
	}
	return b, nil
}





func (r *Registry) BlobUploadInit(image string, uploadID string) (*Blob, error) {
	b := &Blob{
		image: image,
		exists: false,
		uploadID: uploadID,
		uploadFileLocation: r.StorageLocation + "/uploads/" + uploadID,
		uploadUrl: fmt.Sprintf("/v2/%s/blobs/uploads/%s", image, uploadID),
	}

	return b, nil
}


















func (b *Blob) ConvertUpload() error {

	b.digest = 

	f, err := os.Create(b.uploadFileLocation)
	if err != nil {
		return err
	}

	defer f.Close()

	c, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	b.contentLength = c

	return nil
}









func (b *Blob) WriteUploadData(r io.Reader) error {

	f, err := os.Create(b.uploadFileLocation)
	if err != nil {
		return err
	}

	defer f.Close()

	c, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	b.contentLength = c

	return nil
}

func (b *Blob) ReadData(w io.Writer) error {
	f, err := os.Open(b.fileLocation)
	if err != nil {
		return err
	}

	defer f.Close()


	_, err = io.Copy(w, f)

	return err
}

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

		err := b.ReadData(w)
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

	err = b.WriteUploadData(req.Body)
	check(err)

	fmt.Printf("%+v\n", b)


    w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", vars["image"], id))
    w.Header().Set("Range", string(b.contentLength))
	w.WriteHeader(http.StatusAccepted)
}

func (r *Registry) PostBlobUploadComplete(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n --- Blob Upload Complete Post Handler")

	id := filepath.Base(req.URL.Path)

	vars := mux.Vars(req)
	image := vars["image"]

	b, err := r.BlobUploadInit(image, id)
	check(err)



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



