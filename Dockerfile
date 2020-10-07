FROM fedora:32

RUN mkdir /go/

RUN dnf update -y && dnf install golang -y

ENV GOPATH=/go/

RUN chown -R 1000:1000 /go/

RUN go get github.com/go-delve/delve/cmd/dlv && go install github.com/go-delve/delve/cmd/dlv