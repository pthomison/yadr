package registry

import(
	"net/http"
	"fmt"

	"io"

	"github.com/gorilla/mux"
	"encoding/json"

)

func (r *Registry) PutManifest(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n --- Manifest Put Handler Called")

	var err error

	vars := mux.Vars(req)
	reference := vars["reference"]
	imageName := vars["image"]

	image, exists := r.images[imageName]

	if !exists {
		image, err = r.ImageInit(imageName)
		check(err)
		r.images[imageName] = image
	}


	m, err := image.AddManifest(reference, io.Reader(req.Body))
	check(err)

	w.Header().Set("Location", fmt.Sprintf("/v2/%s/manifests/%s", image, m.digest))
	w.Header().Set("Docker-Content-Digest", m.digest)
    w.WriteHeader(http.StatusCreated)
}

func (r *Registry) DeleteManifest(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n --- Manifest Put Handler Called")

	vars := mux.Vars(req)
	reference := vars["reference"]
	image := vars["image"]

	i, exists := r.images[image]

	if !exists {
	    w.WriteHeader(http.StatusNotFound)
	    return
	}

	err := i.DeleteReference(reference)
	check(err)

    w.WriteHeader(http.StatusAccepted)
}

func (r *Registry) GetManifest(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n --- Manifest Put Handler Called")

	var m *Manifest

	vars := mux.Vars(req)
	reference := vars["reference"]
	image := vars["image"]

	i, exists := r.images[image]

	if !exists {
	    w.WriteHeader(http.StatusNotFound)
	    return
	}

	if isHash(reference, hashType) {
		m, exists = i.manifests[reference]
	} else {
		m, exists = i.tags[reference]
	}

	if !exists {
	    w.WriteHeader(http.StatusNotFound)
	    return
	}

	w.Header().Set("Docker-Content-Digest", m.digest)
	w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")

	w.WriteHeader(http.StatusOK)

    m.SendData(w)
}

func (r *Registry) ListTags(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n --- List Tags Called")

	// var m *Manifest

	vars := mux.Vars(req)
	image := vars["image"]

	i, exists := r.images[image]

	if !exists {
	    w.WriteHeader(http.StatusNotFound)
	    return
	}

	tl := &TagList{
		Name: image,
	}

	for k, _ := range i.tags {
		tl.Tags = append(tl.Tags, k)
	}

	b, err := json.Marshal(tl)
	check(err)

	w.Write(b)
}

