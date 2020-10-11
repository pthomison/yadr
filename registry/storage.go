package registry

import (
	"fmt"
	"io/ioutil"
	"os"

	"crypto/sha256"
	"strings"
)

const (
	blobFolder     = "/blobs/"
	manifestFolder = "/manifests/"
	uploadFolder   = "/uploads/"

	storageFolderPerms = 0755
	storageFilePerms   = 0644
)

func hashFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(content)

	return fmt.Sprintf("sha256:%x", hash), nil
}

func (r *Registry) checkForBlob(descriptor string) (bool, int64) {
	files, err := ioutil.ReadDir(r.BlobFolder)
	check(err)

	for _, file := range files {
		if file.Name() == descriptor {
			length := file.Size()
			return true, length
		}
	}

	return false, 0
}

func (r *Registry) ensureImageFolder(image string) error {
	err := os.MkdirAll(r.ManifestFolder+image+"/tags", storageFolderPerms)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.ManifestFolder+image+"/index", storageFolderPerms)
	if err != nil {
		return err
	}

	return nil
}

func isHash(reference string, hashType string) bool {
	if strings.Contains(reference, hashType) {
		return true
	} else {
		return false
	}
}
