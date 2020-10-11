FROM golang:latest AS builder

ENV GOPATH=/go/
ENV CGO_ENABLED=0

COPY ./ /go/src/github.com/pthomison/yadr

WORKDIR /go/src/github.com/pthomison/yadr

RUN go build -o yadr main.go 

RUN chmod +x ./yadr

FROM scratch

COPY --from=builder /go/src/github.com/pthomison/yadr/yadr /yadr

ENTRYPOINT ["/yadr"]