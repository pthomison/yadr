package registry

import(
	// "crypto/sha256"
	"net/http"
	"fmt"
	// "encoding/json"
	// "os"
	"io"
	"github.com/gorilla/mux"

)


// type manifest struct {
// 	SchemaVersion int `json:"schemaVersion"`
// 	MediaType string `json:"mediaType"`
// 	Config descriptor `json:"config"`
// 	Layers []layer `json:"layers"` 
// }

// type layer struct {
// 	descriptor
// }

// type descriptor struct {
// 	content []byte
// 	Size int64 `json:"size"`
// 	MediaType string `json:"mediaType"`
// 	digest
// }


func (r *Registry) ManifestPutHandlerFactory() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("\n\n\n\n\n\n--------Manifest Put Handler")

		// id := uuid.New().String()


		// c, err := io.Copy(os.Stdout, req.Body)
		// check(err)

		vars := mux.Vars(req)

		// fmt.Printf("C: %+v\n", c)
		fmt.Printf("Vars: %+v\n", vars)
	
		r.writeManifestTag(vars["image"], vars["reference"], io.Reader(req.Body))

	    w.Header().Set("Accept-Encoding", "gzip")
	    w.Header().Set("Content-Length", "0")
	    w.Header().Set("Range", "0-0")

	    // w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", vars["image"], id))
	    // w.Header().Set("Docker-Upload-UUID", id)
	    w.WriteHeader(http.StatusAccepted)
	}
}




// func ManifestGetHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Manifest Get Handler Called")
//     w.WriteHeader(http.StatusOK)
// }

// func ManifestPostHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Manifest Post Handler Called")
//     w.WriteHeader(http.StatusOK)
// }

// func ManifestHeadHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Manifest Head Handler Called")
//     w.WriteHeader(http.StatusNotFound)
// }



















































// func (m *manifest) marshal() []byte {
// 	j, err := json.Marshal(m)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}

// 	return j
// }

// func (l *layer) blobify() {
	
// }

// func (d *descriptor) hash() {
// 	d.hashType = "sha256"
// 	d.Digest = fmt.Sprintf("%x", sha256.Sum256(d.content))
// }

// // func (b *descriptor) store() {
// // 	b.Digest.HashType = "sha256"
// // 	b.Digest.Verified = false
// // 	b.Size = len(b.content)
// // 	b.Digest.Digest = fmt.Sprintf("%x", sha256.Sum256(b.content))
// // }

// // func (b *descriptor) load() {
// // 	b.HashType = "sha256"
// // 	b.Verified = false
// // 	b.Size = len(b.content)
// // 	b.Digest.Digest = fmt.Sprintf("%x", sha256.Sum256(b.content))
// // }

// // func newLayer(content []byte) (*layer) {
// // 	l := &layer{ 
// // 		descriptor{
// // 			content: content,
// // 			MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
// // 		},
// // 	}

// // 	l.Size = int64(len(l.content))
// // 	l.verified = false

// // 	return l;
// // }

// // func newManifest(layers []layer) (*manifest) {
// // 	m := &manifest{ 
// // 		SchemaVersion: 2,
// // 		MediaType: "application/vnd.docker.distribution.manifest.v2+json",
// // 		Config: descriptor{
// // 			content: content,
// // 			MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
// // 		},
// // 		Layers: layers,
// // 	}

// // 	m.Size = len(l.content)
// // 	m.Verified = false

// // 	return m
// // }