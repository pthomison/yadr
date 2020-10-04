package registry

import(
	"io/ioutil"
	"fmt"
	"os"

	"crypto/sha256"

)

const(
	blobFolder = "/blobs/"
	manifestFolder = "/manifests/"
	uploadFolder = "/uploads/"
)

func hashFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", nil
	}

	hash := sha256.Sum256(content)

	return fmt.Sprintf("sha256:%x", hash), nil
}

func (r *Registry) checkForBlob(descriptor string) bool {
	files, err := ioutil.ReadDir(r.StorageLocation + blobFolder)
	check(err)

	for _, file := range files {
		fmt.Printf(file.Name())
		if file.Name() == descriptor {
			return true
		}
	}	

	return false
}

func (r *Registry) moveBlobUpload(uuid string, descriptor string) error {
	uploadLocation := r.StorageLocation + uploadFolder + uuid
	blobLocation := r.StorageLocation + blobFolder + descriptor

	return os.Rename(uploadLocation, blobLocation)
} 