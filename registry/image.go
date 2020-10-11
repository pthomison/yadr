package registry

import(
	"os"
	"io"
	"bytes"
	"io/ioutil"
	"fmt"
	"crypto/sha256"
	// "net/http"
	// "github.com/gorilla/mux"

)

type Image struct {
	// manifests stored under their descriptor
	manifests map[string]*Manifest
	// manifests stored under their tag
	tags map[string]*Manifest

	baseLocation string
	indexLocation string
	tagLocation string
	name string
}

type TagList struct {
	Name string
	Tags []string
}

func (r *Registry) ImageInit(name string) (*Image, error) {
	err := r.ensureImageFolder(name)
	if err != nil {
		return nil, err
	}

	basePath := r.ManifestFolder + name

	i := &Image{
		name: name,
		baseLocation: basePath,
		indexLocation: basePath + "/index/",
		tagLocation: basePath + "/tags/",
	}

    i.manifests = make(map[string]*Manifest)
    i.tags = make(map[string]*Manifest)

	err = i.Scan()
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Image) Scan() error {

	fmt.Printf("%+v\n", i)

	indexFiles, err := ioutil.ReadDir(i.indexLocation)
	if err != nil {
		return err
	}

	tagFiles, err := ioutil.ReadDir(i.tagLocation)
	if err != nil {
		return err
	}

	fmt.Printf("files grabbed\n")


	for _, f := range indexFiles {
		digest := f.Name()
		if isHash(digest, hashType) {
			i.manifests[digest] = &Manifest{
				image: i.name,
				fileLocation: i.indexLocation + digest,
				digest: digest,
			}
		}
	}

	fmt.Printf("index done\n")

	for _, f := range tagFiles {
		tag := f.Name()
		d, err := ioutil.ReadFile(i.tagLocation + tag)
		fmt.Printf("%+v\n", tag)

		if err != nil {
			return err
		}
		digest := string(d)
		if isHash(digest, hashType) {
			i.manifests[digest] = &Manifest{
				image: i.name,
				fileLocation: i.indexLocation + digest,
				digest: digest,
			}
		}
	}

	return nil
}

func (i *Image) DeleteReference(reference string) error {
	_, isTag := i.tags[reference]
	m, isDigest := i.manifests[reference]


	if isTag {
		delete(i.tags, reference)
		i.DeleteTag(reference)
	} else if isDigest {
		delete(i.manifests, reference)
		err := m.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Image) AddManifest(reference string, rd io.Reader) (*Manifest, error) {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, rd)
	if err != nil {
		return nil, err
	}

	computedHash := sha256.Sum256(buf.Bytes())
	computedDigest := fmt.Sprintf("sha256:%x", computedHash)

	m, exists := i.manifests[computedDigest]

	if !exists {
		m = &Manifest{
			image: i.name,
			fileLocation: i.indexLocation + computedDigest,
			digest: computedDigest,
		}
		err = m.WriteManifest(rd)
		if err != nil {
			return nil, err
		}
		i.manifests[computedDigest] = m
	}


	if ! isHash(reference, hashType) {
		err := ioutil.WriteFile(i.tagLocation + reference, []byte(computedDigest), 0644)
		if err != nil {
			return nil, err
		}
		i.tags[reference] = m

	}

	return m, nil
}

func (i *Image) DeleteTag(tag string) error {
	err := os.Remove(i.tagLocation + tag)
	if err != nil {
		return err
	}

	return nil
}