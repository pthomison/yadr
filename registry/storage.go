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

// func (r *Registry) sendBlob(descriptor string) error {
// 	blobLocation := r.StorageLocation + blobFolder + descriptor

// 	check(err)

// 	for _, file := range files {
// 		fmt.Printf(file.Name())
// 		if file.Name() == descriptor {
// 			length := file.Size()
// 			return true, length
// 		}
// 	}	

// 	return false, 0
// }

func (r *Registry) moveBlobUpload(uuid string, descriptor string) error {
	uploadLocation := r.StorageLocation + uploadFolder + uuid
	blobLocation := r.StorageLocation + blobFolder + descriptor

	return os.Rename(uploadLocation, blobLocation)
} 

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

// func (r *Registry) ReadManifestTag(image string, tag string, rd io.Reader) (string, int64, error) {

// }