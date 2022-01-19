FROM golang:1.17 as build-env

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOSUMDB=off
ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/meroxa/valve/".insteadOf "https://github.com/meroxa/valve/"
RUN go get github.com/meroxa/valve
RUN go mod init github.com/meroxa/app
RUN go mod vendor
RUN go build -mod vendor -tags platform -o /go/bin/app

FROM gcr.io/distroless/static
USER nonroot:nonroot
COPY --from=build-env /go/bin/app /
ENTRYPOINT ["/app", "--serve"]

