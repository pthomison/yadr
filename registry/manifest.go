package registry

import(
	"io/ioutil"
	"strings"

	// "io/ioutil"
	"fmt"
	"os"

	"crypto/sha256"

	"bytes"
	"io"
)

type Manifest struct {
	blob *Blob
	tag string
	digest string
	image string
}

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


// func (r *Registry) manifestInit(image string, reference string) (*Manifest, error) {

// }

func (r *Registry) writeManifestDigest(image string, rd io.Reader) (string, int64, error) {
	err := r.createImageFolder(image)
	if err != nil {
		return "", 0, err
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, rd)
	if err != nil {
		return "", 0, err
	}

	hash := sha256.Sum256(buf.Bytes())
	digest := fmt.Sprintf("sha256:%x", hash)

	// fmt.Printf("Buff: %+v\n", string(buf.Bytes()))
	// fmt.Printf("Hash: %+v\n", digest)

    f, err := os.Create(r.StorageLocation + manifestFolder + image + "/index/" + digest )

    if err != nil {
        return "", 0, err
    }

    defer f.Close()

	c, err := io.Copy(f, &buf)

	return digest, c, nil
} 

func (r *Registry) writeManifestTag(image string, tag string, rd io.Reader) (string, int64, error) {
	digest, size, err := r.writeManifestDigest(image, rd)
	if err != nil {
		return "", 0, err
	}

	err = os.MkdirAll(r.StorageLocation + manifestFolder + image + "/tags/" + tag, storageFolderPerms)
	if err != nil {
		return "", 0, err
	}

    f, err := os.Create(r.StorageLocation + manifestFolder + image + "/tags/" + tag + "/link" )

    if err != nil {
		return "", 0, err
    }

    defer f.Close()

	_, err = f.WriteString(digest)

    if err != nil {
		return "", 0, err
    }

	return digest, size, nil
} 


func (r *Registry) readManifest(image string, hash string) ([]byte, error) {
	manifestLocation := r.StorageLocation + manifestFolder + image + "/index/" + hash

	content, err := ioutil.ReadFile(manifestLocation)
	if err != nil {
		return nil, err
	} else {
		return content, nil
	}
}

func (r *Registry) lookupTag(image string, tag string) (string, error) {
	tagLocation := r.StorageLocation + manifestFolder + image + "/tags/" + tag + "/link"
	content, err := ioutil.ReadFile(tagLocation)
	if err != nil {
		return "", err
	} else {
		return string(content), nil
	}
}

func isHash(reference string, hashType string) bool {
	if strings.Contains(reference, hashType) {
		return true
	} else {
		return false
	}
}