# syntax = docker/dockerfile:1-experimental

FROM golang:1.17.3-alpine3.14 as builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY confs/dev/live.json /root/.gocryptotrader/config.json
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gct .

WORKDIR /app/cmd/dbseed
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dbseed .

WORKDIR /app/cmd/dbmigrate
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dbmigrate .

WORKDIR /app
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates postgresql-client
#
# WORKDIR /root/
#
# # Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
# COPY confs/dev/live.json /root/.gocryptotrader/config.json
# COPY --from=builder /app/gct .
# COPY --from=builder /app/database ./database
# COPY --from=builder /app/cmd/dbseed/dbseed .
# COPY --from=builder /app/cmd/dbmigrate/dbmigrate .
# COPY --from=builder /app/confs ./confs
# # COPY --from=builder /app/.env .
# EXPOSE 8080
# # ENTRYPOINT [""]
