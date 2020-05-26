ARG PACKAGE_NAME=github.com/syucream/hakagi

FROM golang:1.14.3-alpine3.11
ARG PACKAGE_NAME

RUN apk --no-cache add make git curl
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/$PACKAGE_NAME
RUN cd /go/src/$PACKAGE_NAME && make dep && make
RUN cp /go/src/$PACKAGE_NAME/hakagi /bin/.

CMD hakagi

# sh -c "apk --no-cache add libc-dev gcc && cd /go/src/github.com/syucream/hakagi && make test"
