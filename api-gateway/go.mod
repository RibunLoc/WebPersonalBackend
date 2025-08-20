module github.com/RibunLoc/WebPersonalBackend/api-gateway

go 1.23.4

require (
	github.com/RibunLoc/WebPersonalBackend/gen v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/cors v1.2.1
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.75.0
)

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)

replace github.com/RibunLoc/WebPersonalBackend/gen => ../gen
