FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
      -o /src/app ./cmd/app

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /src/app /app

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
