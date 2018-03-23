FROM golang:latest
MAINTAINER Kris Nova "kris@nivenly.com"
ADD . /go/src/github.com/kris-nova/kale
WORKDIR /go/src/github.com/kris-nova/kale
RUN make compile
CMD ["./bin/kale"]