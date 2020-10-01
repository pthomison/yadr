package registry

import(
	"net/http"
)

const(
	BaseAPI = "/v2"
)

func BaseHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}