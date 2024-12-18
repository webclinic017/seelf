FROM node:18-alpine AS front_builder
WORKDIR /app
COPY ./cmd/serve/front/package*.json .
RUN npm ci
COPY ./cmd/serve/front .
RUN npm run build

FROM golang:1.23-alpine AS builder
# build-base needed to compile the sqlite3 dependency
RUN apk add --update-cache build-base
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=front_builder /app/build ./cmd/serve/front/build
RUN make build-back

FROM alpine:3.16
LABEL org.opencontainers.image.authors="julien@leicher.me" \
    org.opencontainers.image.source="https://github.com/YuukanOO/seelf"
RUN apk add --update-cache openssh-client && \
    rm -rf /var/cache/apk/*
ENV DATA_PATH=/seelf/data
WORKDIR /app
COPY --from=builder /app/seelf ./
EXPOSE 8080
CMD ["./seelf", "-c", "/seelf/data/conf.yml", "serve"]