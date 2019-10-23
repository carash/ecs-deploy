FROM golang:1.13.1-alpine as builder

RUN apk add -qq --no-cache \
    git \
    upx \
    ca-certificates

WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod verify

ENV GOOS=linux GOARCH=386
COPY . /src/
RUN go build -ldflags "-s -w" -o /app ./cmd/ecs-deploy && \
    upx --brute /app

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

CMD ["/app"]
