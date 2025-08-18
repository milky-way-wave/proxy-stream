FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) \
  CGO_ENABLED=0 \
  go build -a -installsuffix cgo \
  -ldflags '-extldflags "-static" -s -w' -o app .

FROM alpine:3.19
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata dumb-init && \
  addgroup -g 1001 -S go-app && \
  adduser -u 1001 -S go-app -G go-app && \
  chown -R go-app:go-app /app
COPY --from=builder --chown=1001:1001 /app/app .
USER 1001:1001
CMD ["./app", "serve"]
