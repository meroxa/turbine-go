FROM golang:1.17 as build-env

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go install --tags=platform -mod=vendor ./examples/enrich/...

FROM gcr.io/distroless/static
USER nonroot:nonroot
WORKDIR /app
COPY --from=build-env /go/bin/* /app
COPY --from=build-env /go/src/app/examples/enrich/app.json /app
ENTRYPOINT ["/app/enrich", "--serve"]
