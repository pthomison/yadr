package registry

import(
	"crypto/sha256"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"

)

const(
	ManifestAPI = "/{image}/manifests/{reference}"
	BlobAPI = "/{image}/blobs/{digest}"
)

type manifest struct {
	SchemaVersion int
	MediaType string
	Config struct {
		MediaType string
	}
	Layers []layer
	Blob blob
}

type layer struct {
	Blob blob
}

type blob struct {
	Digest digest
	Content []byte
	Size int
	MediaType string
}

type digest struct {
	HashType string
	Verified bool
	Digest string
}

func ManifestGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%+v\n", vars)
    w.WriteHeader(http.StatusOK)
}

func BlobGetHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func (m *manifest) blobify() {

}

func (l *layer) blobify() {
	
}

func (b *blob) hash() {
	b.Digest.HashType = "sha256"
	b.Digest.Verified = false
	b.Size = len(b.Content)
	b.Digest.Digest = fmt.Sprintf("%x", sha256.Sum256(b.Content))
}

func (b *blob) store() {
	b.Digest.HashType = "sha256"
	b.Digest.Verified = false
	b.Size = len(b.Content)
	b.Digest.Digest = fmt.Sprintf("%x", sha256.Sum256(b.Content))
}

func (b *blob) load() {
	b.Digest.HashType = "sha256"
	b.Digest.Verified = false
	b.Size = len(b.Content)
	b.Digest.Digest = fmt.Sprintf("%x", sha256.Sum256(b.Content))
}