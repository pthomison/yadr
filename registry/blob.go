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
	exists bool
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


func (r *Registry) BlobUploadInit(image string, uploadID string) (*BlobUpload, error) {
	b := &BlobUpload{
		image: image,
		id: uploadID,
		fileLocation: r.StorageLocation + "/uploads/" + uploadID,
		url: fmt.Sprintf("/v2/%s/blobs/uploads/%s", image, uploadID),
	}

	return b, nil
}


func (r *Registry) ProcessUpload(u *BlobUpload)  (*Blob, error) {

	digest, err := hashFile(u.fileLocation)
	if err != nil {
		return nil, err
	}

	b, err := r.BlobDigestInit(digest, u.image)
	if err != nil {
		return nil, err
	}

	err = os.Rename(u.fileLocation, b.fileLocation)
	if err != nil {
		return nil, err
	}

	return b, nil
}


func (b *BlobUpload) StoreUploadData(r io.Reader) error {

	f, err := os.Create(b.fileLocation)
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