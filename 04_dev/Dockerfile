# syntax=docker/dockerfile:1.7

# ---------- builder ----------
    FROM golang:1.22-bookworm AS builder
    WORKDIR /src
    
    RUN --mount=type=cache,target=/go/pkg/mod true
    RUN --mount=type=cache,target=/root/.cache/go-build true
    
    COPY go.mod go.sum ./
    
    RUN --mount=type=cache,target=/go/pkg/mod go mod download
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build \
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -trimpath -buildvcs=false -ldflags="-s -w" \
        -o /out/jenkins-cicd-go-app .
    
    # ---------- runtime (alpine) ----------
    FROM alpine:3.20
    WORKDIR /app
    RUN apk add --no-cache ca-certificates
    RUN adduser -D -u 65532 appuser
    COPY --chown=appuser:appuser --from=builder /out/jenkins-cicd-go-app /app/jenkins-cicd-go-app
    COPY --chown=appuser:appuser index.html /app/index.html
    
    EXPOSE 9000
    USER 65532:65532
    ENV SERVICE_NAME=calculator
    ENV HOST=0.0.0.0
    ENV PORT=9000
    ENTRYPOINT ["/app/jenkins-cicd-go-app"]
    
    