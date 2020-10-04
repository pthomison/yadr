package registry

import(
	"io/ioutil"
	"fmt"
	"os"

	"crypto/sha256"

	"bytes"
	"io"

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

func (r *Registry) writeManifestDigest(image string, rd io.Reader) (string, error) {
	err := r.createImageFolder(image)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, rd)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(buf.Bytes())
	digest := fmt.Sprintf("sha256:%x", hash)

	fmt.Printf("Buff: %+v\n", string(buf.Bytes()))
	fmt.Printf("Hash: %+v\n", digest)

    f, err := os.Create(r.StorageLocation + manifestFolder + image + "/index/" + digest )

    if err != nil {
        return "", err
    }

    defer f.Close()

	_, err = io.Copy(f, &buf)

	return digest, nil
} 

func (r *Registry) writeManifestTag(image string, tag string, rd io.Reader) error {
	digest, err := r.writeManifestDigest(image, rd)
	if err != nil {
		return err
	}

    f, err := os.Create(r.StorageLocation + manifestFolder + image + "/tags/" + tag + "/link" )

    if err != nil {
        return err
    }

    defer f.Close()

	_, err = f.WriteString(digest)

	return nil
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