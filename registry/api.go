package registry

const (
	BaseAPI = "/v2/"

	BlobUploadRequestAPI = BaseAPI + "{image}/blobs/uploads/"

	BlobUploadAPI = BaseAPI + "{image}/blobs/uploads/{sessionID}"

	BlobAPI = BaseAPI + "{image}/blobs/{digest}"

	ManifestAPI = BaseAPI + "{image}/manifests/{reference}"

	TagAPI = BaseAPI + "{image}/tags/list"
)
