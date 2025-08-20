FROM golang:1.23 AS build 
WORKDIR /src

# Tối ưu cache module
COPY go.mod go.sum ./
RUN go mod download 

# Copy toàn bộ mã nguồn 
COPY . .

# Build static binary (không ccaanf CGO)
ARG VERSION=dev
ARG COMMIT=none
ARG DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/contact-service \
        -ldflags="-s -w \
        -X main.version=${VERSION} \
        -X main.commit=${COMMIT} \
        -X main.buildDate=${DATE}" \
    .

# ---- Runtime stage -----
# Debian 12 distroless (có CA cert), user non-root
FROM gcr.io/distroless/base-debian12:nonroot

# Copy binary 
COPY --from=build /out/contact-service /contact-service

# Chạy với user non root 
USER nonroot:nonroot 

EXPOSE 8082 50052

ENTRYPOINT [ "/contact-service" ]

