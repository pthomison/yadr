package registry

import(
	"os"
	"io"
	// "io/ioutil"
	"fmt"

)

type Blob struct {
	digest string
	contentLength int64
	contentType string
	fileLocation string
	urlLocation string
	image string
}

type BlobUpload struct {
	digest string
	contentLength int64
	url string
	fileLocation string
	image string
	id string
}

func (r *Registry) BlobInit(digest string) (*Blob, error) {
	b := &Blob{
		digest: digest,
		fileLocation: fmt.Sprintf("%v%v", r.BlobFolder, digest),
	}

	exists, _ := r.checkForBlob(digest)

	if exists {
		f, err := os.Stat(b.fileLocation)
		if err != nil {
			return nil, err
		}

		b.contentLength = f.Size()
	}
	return b, nil
}

func (r *Registry) BlobUploadInit(image string, id string) *BlobUpload {
	b := &BlobUpload{
		image: image,
		id: id,
		fileLocation: r.UploadFolder + id,
		url: fmt.Sprintf("/v2/%s/blobs/uploads/%s", image, id),
	}

	return b
}


func (r *Registry) ProcessUpload(u *BlobUpload)  (*Blob, error) {
	digest, err := hashFile(u.fileLocation)
	if err != nil {
		return nil, err
	}

	b, err := r.BlobInit(digest)
	if err != nil {
		return nil, err
	}

	err = os.Rename(u.fileLocation, b.fileLocation)
	if err != nil {
		return nil, err
	}

	r.blobs[b.digest] = b

	return b, nil
}


func (b *BlobUpload) StoreUploadData(r io.Reader) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(b.fileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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


func (b *Blob) SendData(w io.Writer) error {
	f, err := os.Open(b.fileLocation)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(w, f)

	return err
}

func (r *Registry) Delete(b *Blob) error {
	err := os.Remove(b.fileLocation)
	if err != nil {
		return err
	}

	return nil
}