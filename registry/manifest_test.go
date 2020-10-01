package registry

import "testing"
import "fmt"

func TestLoadImage(t *testing.T) {
	var m := &manifest{
		SchemaVersion: 2,
		MediaType: "application/vnd.oci.image.config.v1+json",
		Config: {
			MediaType: "application/vnd.oci.image.config.v1+json",
		},
	}

	var layerLocation


}