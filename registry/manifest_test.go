package registry

import(
	"io/ioutil"
	"testing"
	"fmt"
	"os"
)

const(
	testData = "/hacking/testing-data"
	storeDir = "/hacking/registry-data"
	layerLocation = testData + "/fedora-core-layer"

)

func TestStoreImage(t *testing.T) {
	layerContent, err := ioutil.ReadFile(layerLocation)
	if err != nil {
		t.Fatal(err)
	}

	l := NewLayer(layerContent)
	m := NewManifest([]Layer{l})
	
	fmt.Printf("%+v\n", string(m.marshal()))

	err = createFolders()
	if err != nil {
		t.Fatal(err)
	}





	_ = m






	// err = ioutil.WriteFile(storeDir + "/image", content, 0644)

	// if err != nil {
	// 	t.Fatal(err)
	// }

}

func createFolders() error {
	err := os.Mkdir(storeDir, 0644)
	if err != nil {
		return err
	}

	err = os.Mkdir(storeDir + "/blobs", 0644)
	if err != nil {
		return err
	}

	err = os.MkdirAll(storeDir + "/manifests/fedora/tags/latest", 0644)
	if err != nil {
		return err
	}

	return nil
}