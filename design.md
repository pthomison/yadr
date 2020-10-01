lightweight golang impmentation of the docker V2 registry API

https://github.com/opencontainers/distribution-spec/blob/master/spec.md
https://docs.docker.com/registry/spec/api/
https://ops4j.github.io/ramler/0.6.0/registry/

- unauthorized
- implements pull
- digest hash == sha256
- data storage will be in a flat folder with every digest having a folder with its data
- manifests and layers will be stored together
- manifests will be stored as json files
- layers will be stored as tar balls?




### Definitions

Several terms are used frequently in this document and warrant basic definitions:

- **Registry**: a service that handles the required APIs defined in this specification
- **Client**: a tool that communicates with registries
- **Push**: the act of uploading blobs and manifests to a registry
- **Pull**: the act of downloading blobs and manifests from a registry
- **Blob**: the binary form of content that is stored by a registry, addressable by a digest
- **Manifest**: a JSON document which defines an artifact. Manifests are defined under the [OCI Image Spec](https://github.com/opencontainers/image-spec/blob/master/manifest.md)
- **Config**: a section in the manifest (and associated blob) which contains artifact metadata
- **Artifact**: one conceptual piece of content stored as blobs with an accompanying manifest containing a config
- **Digest**: a unique identifier created from a cryptographic hash of a blob's content. Digests are defined under the [OCI Image Spec](https://github.com/opencontainers/image-spec/blob/master/descriptor.md)
- **Tag**: a custom, human-readable manifest identifier

API:

/v2/: base

	/{{repo}}/
	repo == repository == image

		/manifests/{{ref}}
		ref must be tag or manifest digest

		/blobs/{{digest}}

