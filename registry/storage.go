package registry

import(
	"io/ioutil"
	"fmt"
	"os"

	"crypto/sha256"

	// "bytes"
	// "io"

)

const(
	blobFolder = "/blobs/"
	manifestFolder = "/manifests/"
	uploadFolder = "/uploads/"

	storageFolderPerms = 0777
	storageFilePerms = 0777
)

func hashFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", nil
	}

	hash := sha256.Sum256(content)

	return fmt.Sprintf("sha256:%x", hash), nil
}

func (r *Registry) checkForBlob(descriptor string) (bool, int64) {
	files, err := ioutil.ReadDir(r.StorageLocation + blobFolder)
	check(err)

	for _, file := range files {
		fmt.Printf(file.Name())
		if file.Name() == descriptor {
			length := file.Size()
			return true, length
		}
	}	

	return false, 0
}

func (r *Registry) createImageFolder(image string) error {
	err := os.MkdirAll(r.StorageLocation + manifestFolder + image + "/tags", storageFolderPerms)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.StorageLocation + manifestFolder + image + "/index", storageFolderPerms)
	if err != nil {
		return err
	}

	return nil
} 