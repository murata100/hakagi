ARG PACKAGE_NAME=github.com/murata100/hakagi

FROM golang:1.13.15-alpine3.11
ARG PACKAGE_NAME

RUN apk --no-cache add make git curl
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/$PACKAGE_NAME
RUN cd /go/src/$PACKAGE_NAME && make dep && make
RUN cp /go/src/$PACKAGE_NAME/hakagi /bin/.

CMD hakagi

# sh -c "apk --no-cache add libc-dev gcc && cd /go/src/github.com/murata100/hakagi && make test"
