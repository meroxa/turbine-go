FROM golang:1.17 as build-env

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -mod vendor -tags platform -o /go/bin/app

FROM gcr.io/distroless/static
USER nonroot:nonroot
COPY --from=build-env /go/bin/app /
ENTRYPOINT ["/app", "--serve"]
