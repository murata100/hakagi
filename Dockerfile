ARG PACKAGE_NAME=github.com/syucream/hakagi.git

FROM golang:1.14.3-alpine3.11
ARG PACKAGE_NAME

RUN apk --no-cache add make git curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . /go/src/$PACKAGE_NAME
RUN cd /go/src/$PACKAGE_NAME && make dep && make
RUN cp /go/src/$PACKAGE_NAME/hakagi /bin/.

CMD hakagi
