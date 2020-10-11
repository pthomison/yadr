package registry

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
    "io/ioutil"

    "github.com/sirupsen/logrus"
)

const(
	hashType = "sha256"
)

type Registry struct {
	StorageRoot string
    BlobFolder string
    UploadFolder string
    ManifestFolder string

    // key should be blob hash
    blobs map[string]*Blob

    // key should be image name
    images map[string]*Image
}

func New(storagePath string) (*Registry, error) {
	r := &Registry{
		StorageRoot: storagePath,
        BlobFolder: storagePath + blobFolder,
        UploadFolder: storagePath + uploadFolder,
        ManifestFolder: storagePath + manifestFolder,
	}

	err := r.InitializeStorage()
	if err != nil {
		return nil, err
	}

    r.blobs = make(map[string]*Blob)
    err = r.ScanBlobs()
    if err != nil {
        return nil, err
    }

    r.images = make(map[string]*Image)
    err = r.ScanManifests()
    if err != nil {
        return nil, err
    }    

	r.SetAPI()

	return r, nil
}



func (r *Registry) Serve() {
    http.ListenAndServe(":5000", nil)
}

func (r *Registry) SetAPI() {
    // instantiate router
    router := mux.NewRouter()

    // Serve 200 for v2 api endpoint
    router.HandleFunc(BaseAPI, BaseGetHandler).
    	Methods(http.MethodGet)

    // MANIFEST HANDLERS
    router.HandleFunc(ManifestAPI, r.PutManifest).
    	Methods(http.MethodPut)

    router.HandleFunc(ManifestAPI, r.GetManifest).
        Methods(http.MethodGet)

    router.HandleFunc(ManifestAPI, r.DeleteManifest).
        Methods(http.MethodDelete)

    router.HandleFunc(TagAPI, r.ListTags).
        Methods(http.MethodGet)


    // BLOB HANDLERS
    // Head/Check
    router.HandleFunc(BlobAPI, r.HeadBlob).
        Methods(http.MethodHead)

    // Get
    router.HandleFunc(BlobAPI, r.GetBlob).
    	Methods(http.MethodGet)

    // Delete
    router.HandleFunc(BlobAPI, r.DeleteBlob).
    	Methods(http.MethodDelete)

    // Monolithic, 1 step Upload
    router.HandleFunc(BlobUploadRequestAPI, r.CompleteBlobUpload).
        Methods(http.MethodPost).    
        Queries("digest", "{digest}")

    // Monolithic, 2 step Upload && Chunked, 3 step Upload
    router.HandleFunc(BlobUploadRequestAPI, r.RequestBlobUpload).
        Methods(http.MethodPost)

    router.HandleFunc(BlobUploadAPI, r.PatchBlobUpload).
        Methods(http.MethodPatch)

    router.HandleFunc(BlobUploadAPI, r.CompleteBlobUpload).
        Methods(http.MethodPut).
        Queries("digest", "{digest}")


    // Serve 404 to all other requests
    router.PathPrefix("/").HandlerFunc(DefaultHandler)

    // attach gorilla router
    http.Handle("/", router)
}

func (r *Registry) InitializeStorage() error {

    for _, folder := range []string{r.StorageRoot, r.BlobFolder, r.UploadFolder, r.ManifestFolder} {
        err := os.MkdirAll(folder, storageFolderPerms)
        if err != nil {
            return err
        }        
    }

	return nil
} 

func check(e error) {
	if e != nil {
        logrus.Panic(e)        
	}
}

func (r *Registry) ScanBlobs() error {
    blobFiles, err := ioutil.ReadDir(r.BlobFolder)
    if err != nil {
        return err
    }

    for _, file := range blobFiles {
        hash := file.Name()

        b, err := r.BlobInit(hash)
        if err != nil {
            return err
        }
        r.blobs[hash] = b
    }   

    return nil
}

func (r *Registry) ScanManifests() error {
    imageFolders, err := ioutil.ReadDir(r.ManifestFolder)
    if err != nil {
        return err
    }

    for _, f := range imageFolders {
        image := f.Name()

        i, err := r.ImageInit(image)
        if err != nil {
            return err
        }

        r.images[image] = i 
    }   

    return nil
}
