FROM golang:alpine
WORKDIR /go/src/github.com/cksharma11/guessing/

ADD . .
RUN apk update && apk add --no-cache git ca-certificates make && go get ./... && make guessing

FROM golang:alpine
WORKDIR /app
COPY --from=0 /go/src/github.com/cksharma11/guessing/bin/guessing ./guessing
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
EXPOSE 8080
ENTRYPOINT [ "/app/guessing" ]