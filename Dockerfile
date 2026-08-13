FROM golang:1.25 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /fluxkv .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /fluxkv /fluxkv

WORKDIR /app

EXPOSE 8000 8001 9000 9001

ENTRYPOINT ["/fluxkv"]
