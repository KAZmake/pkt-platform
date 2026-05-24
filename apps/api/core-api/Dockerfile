FROM golang:1.25-alpine AS builder
WORKDIR /workspace

# Download deps first (cache layer)
COPY go.work go.work.sum ./
COPY core-api/go.mod core-api/go.sum ./core-api/
COPY expertise-svc/go.mod expertise-svc/go.sum ./expertise-svc/
COPY sync-svc/go.mod sync-svc/go.sum ./sync-svc/
COPY assistant-svc/go.mod assistant-svc/go.sum ./assistant-svc/
RUN go work sync

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/core-api ./core-api/cmd/server

# ─── Final image ──────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/core-api /core-api
EXPOSE 8080
ENTRYPOINT ["/core-api"]
