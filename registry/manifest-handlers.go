package registry

import (
	"fmt"
	"net/http"

	"io"

	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

func (r *Registry) PutManifest(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Manifest Put Handler Called")

	var err error

	vars := mux.Vars(req)
	reference := vars["reference"]
	imageName := vars["image"]

	image, exists := r.images[imageName]

	if !exists {
		image, err = r.ImageInit(imageName)
		check(err)
		r.images[imageName] = image
		logrus.Info("Image Added: ", imageName)
	}

	m, err := image.AddManifest(reference, io.Reader(req.Body))
	check(err)
	logrus.Info("Manifest Added: ", reference)

	w.Header().Set("Location", fmt.Sprintf("/v2/%s/manifests/%s", image, m.digest))
	w.Header().Set("Docker-Content-Digest", m.digest)
	w.WriteHeader(http.StatusCreated)
}

func (r *Registry) DeleteManifest(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Manifest Delete Handler Called")

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
	logrus.Info("Manifest Deleted: ", reference)

	w.WriteHeader(http.StatusAccepted)
}

func (r *Registry) GetManifest(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("Manifest Get Handler Called")

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
	// probably shouldn't be hardcoded.... TODO
	w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")

	w.WriteHeader(http.StatusOK)

	m.SendData(w)
}

func (r *Registry) ListTags(w http.ResponseWriter, req *http.Request) {
	logrus.Debug("List Tags Called")

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
