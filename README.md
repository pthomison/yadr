### YADR: Yet Another Docker Registry

YADR is an attempt to build an OCI compliant docker registry with as little complexity as possible

No authentication, no encryption, no problem? This is a learning project that should not run anywhere close to a production environment

#### Implemtation Status:
| Workflow           | Status                                                                                           |
|--------------------|--------------------------------------------------------------------------------------------------|
| Pull               | Done, but needs better error checking                                                            |
| Push               | Done, but needs better error checking                                                            |
| Content Discovery  | Can list tags, but does not implement the catalog endpoint                                       |
| Content Management | Content can be deleted explicitly, but there is no garbage collection for orphaned blogs or tags |


#### Go Packages Used:

Gorilla Mux: github.com/gorilla/mux
Logrus: github.com/sirupsen/logrus
UUID: github.com/google/uuid
Cobra: github.com/spf13/cobra

#### Standards/Documentation Used:
OCI Spec: https://github.com/opencontainers/distribution-spec/blob/master/spec.md
Docker Registry API (V2) Spec: https://github.com/docker/distribution/blob/5cb406d511b7b9163bff9b6439072e4892e5ae3b/docs/spec/api.md


#### Do the thing:

```make runimage```

Seperate term:

```make compliance```
or
```make test```

The compliance task will render an html document in the project directory with the test results

TODO List:
- Need to check manifest type (ie parse the json content)
- Check & correctly deal with Range headers on chunked upload
