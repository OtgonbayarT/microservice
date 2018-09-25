FROM golang:1.8 as buildstage
ENV SRC=/go/src/github.com/OtgonbayarT/
RUN mkdir -p /go/src/github.com/OtgonbayarT/
WORKDIR /go/src/github.com/OtgonbayarT/microservice
RUN go get github.com/rapidloop/skv github.com/prometheus/client_golang/prometheus
RUN git clone -b master https://github.com/OtgonbayarT/microservice.git /go/src/github.com/OtgonbayarT/microservice/ \
&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -o bin/go_docker

FROM alpine:3.6 as baseimagealp
RUN apk add --no-cache bash
ENV WORK_DIR=/docker/bin
WORKDIR $WORK_DIR
COPY --from=buildstage /go/src/github.com/OtgonbayarT/microservice/bin/ ./
ENTRYPOINT /docker/bin/go_docker
EXPOSE 8080