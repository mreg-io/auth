# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

ARG TARGETOS

ARG TARGETARCH

WORKDIR /go/src/app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./

RUN --mount=type=secret,id=netrc,target=/root/.netrc \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go mod download && \
    go mod verify

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -v -o /go/bin/auth-server gitlab.mreg.io/my-registry/auth/cmd/auth-server

FROM gcr.io/distroless/static-debian12:nonroot

EXPOSE 8080

COPY --from=builder --chown=nonroot:nonroot /go/bin/auth-server /

ENTRYPOINT ["/auth-server"]
