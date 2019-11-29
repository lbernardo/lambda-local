FROM golang:1.13
WORKDIR /var/app
COPY . /go/src/github.com/lbernardo/lambda-local
COPY docker/start.server.sh /bin/start.server.sh
RUN go get -u github.com/ghodss/yaml github.com/spf13/cobra github.com/docker/docker/client
RUN go install github.com/lbernardo/lambda-local
ENTRYPOINT [ "sh","/bin/start.server.sh" ]