FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY swaggers/externalapi/models ./swaggers/externalapi/models
RUN CGO_ENABLED=0 go build -trimpath -o /out/users ./cmd/externalApi

FROM alpine:3.21
COPY --from=build /out/users /usr/local/bin/users
USER 65532:65532
EXPOSE 8083 9091
ENTRYPOINT ["/usr/local/bin/users"]
