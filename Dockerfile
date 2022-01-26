FROM golang:1.17 as build-env

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN #go build -mod vendor -tags platform -o /go/bin/app/example/simple/...
RUN --mount=type=cache,target=/root/.cache/go-build go install --tags=platform -mod=vendor ./examples/simple/...

#FROM gcr.io/distroless/static
#USER nonroot:nonroot
#WORKDIR /app
#COPY --from=build-env /go/bin/* /app
#COPY --from=build-env /go/src/app/examples/simple/app.json /app
#ENTRYPOINT ["/app/app", "--serve"]

FROM ubuntu
WORKDIR /app
ENV PATH="/app:${PATH}"

COPY --from=build-env /go/bin/* /app
COPY --from=build-env /go/src/app/examples/simple/app.json /app

ENTRYPOINT ["/app/simple", "--serve"]