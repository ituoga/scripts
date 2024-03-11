FROM golang:1.21.6-bullseye AS build

ARG GITHUB_PAT
ENV GOPRIVATE=github.com

RUN apt update && apt install ca-certificates

WORKDIR /src
COPY go.* .
COPY gw/main.go main.go

RUN go mod download

RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /app main.go

FROM scratch AS bin
COPY --from=build /app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]