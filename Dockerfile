FROM golang:latest AS builder

ENV GOPATH=/go/

COPY ./ /go/src/github.com/pthomison/yadr

WORKDIR /go/src/github.com/pthomison/yadr

RUN go build -o yadr main.go 

FROM scratch

COPY --from=builder /go/src/github.com/pthomison/yadr/yadr /yadr

CMD /yadr