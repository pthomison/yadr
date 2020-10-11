package registry

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func BaseGetHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Base Get Handler Called")
	w.WriteHeader(http.StatusOK)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Error("Default Handler Called: ", r.URL.String())

	w.WriteHeader(http.StatusNotFound)
}
