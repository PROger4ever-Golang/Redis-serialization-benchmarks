# "ARG before FROM" available for recent docker releases only
ARG GOLANG_TAG

FROM golang:${GOLANG_TAG} as build-env
ARG PROJECT_PATH

COPY . /go/src/${PROJECT_PATH}

WORKDIR /go/src/${PROJECT_PATH}
RUN go get

ENTRYPOINT ["/usr/local/go/bin/go"]