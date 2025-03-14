FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /song
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /song/cmd/song
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o song/main .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder song/. .
WORKDIR /song/cmd/song
COPY --from=builder song/cmd/song .
EXPOSE 8080
CMD ["song/main"]