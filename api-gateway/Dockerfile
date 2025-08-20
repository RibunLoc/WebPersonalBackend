FROM golang:1.23 AS build

WORKDIR /src 

# Copy gen
COPY gen /src/gen

# Tối ưu cache module
COPY api-gateway/go.mod api-gateway/go.sum /src/api-gateway/
WORKDIR /src/api-gateway
RUN go mod download 

# Copy toàn bộ mã nguồn 
COPY api-gateway /src/api-gateway/ 

# Build static binary 
ARG VERSION=dev
ARG COMMIT=none
ARG DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/api-gateway \
        -ldflags="-s -w \
        -X main.version=${VERSION} \
        -X main.commit=${COMMIT} \
        -X main.buildDate=${DATE} " \
        .
# --- Runtime stage ----
# Debian 12 distroless, user non-root
FROM gcr.io/distroless/base-debian12:nonroot

# Copy binary 
COPY --from=build /out/api-gateway /api-gateway

USER nonroot:nonroot 

EXPOSE 3000

ENTRYPOINT [ "/api-gateway" ]
