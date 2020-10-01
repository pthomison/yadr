package registry

import(
	"crypto/sha256"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"encoding/json"

)

const(
	ManifestAPI = "/{image}/manifests/{reference}"
	BlobAPI = "/{image}/blobs/{digest}"
)

type manifest struct {
	SchemaVersion int `json:"schemaVersion"`
	MediaType string `json:"mediaType"`
	Config descriptor `json:"config"`
	Layers []layer `json:"layers"` 
}

type layer struct {
	descriptor
}

type descriptor struct {
	content []byte
	Size int64 `json:"size"`
	MediaType string `json:"mediaType"`
	digest
}

type digest struct {
	hashType string
	verified bool
	Digest string `json:"digest"`
}

func ManifestGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%+v\n", vars)
    w.WriteHeader(http.StatusOK)
}

func BlobGetHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func (m *manifest) marshal() []byte {
	j, err := json.Marshal(m)
	if err != nil {
		fmt.Println("error:", err)
	}

	return j
}

func (l *layer) blobify() {
	
}

func (d *descriptor) hash() {
	d.hashType = "sha256"
	d.Digest = fmt.Sprintf("%x", sha256.Sum256(d.content))
}

// func (b *descriptor) store() {
// 	b.Digest.HashType = "sha256"
// 	b.Digest.Verified = false
// 	b.Size = len(b.content)
// 	b.Digest.Digest = fmt.Sprintf("%x", sha256.Sum256(b.content))
// }

// func (b *descriptor) load() {
// 	b.HashType = "sha256"
// 	b.Verified = false
// 	b.Size = len(b.content)
// 	b.Digest.Digest = fmt.Sprintf("%x", sha256.Sum256(b.content))
// }

func newLayer(content []byte) (*layer) {
	l := &layer{ 
		descriptor{
			content: content,
			MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
		},
	}

	l.Size = int64(len(l.content))
	l.verified = false

	return l;
}

func newManifest(layers []layer) (*manifest) {
	m := &manifest{ 
		SchemaVersion: 2,
		MediaType: "application/vnd.docker.distribution.manifest.v2+json",
		Config: descriptor{
			content: content,
			MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
		},
		Layers: layers,
	}

	m.Size = len(l.content)
	m.Verified = false

	return m
}