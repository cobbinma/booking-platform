module github.com/cobbinma/booking-platform/lib/venue_api

go 1.14

require (
	github.com/auth0-community/go-auth0 v1.0.1-0.20191119091237-b9b0f95be568
	github.com/cobbinma/booking-platform/lib/protobuf v0.0.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/joho/godotenv v1.3.0
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	google.golang.org/grpc v1.35.0
	gopkg.in/square/go-jose.v2 v2.5.1
)

replace github.com/cobbinma/booking-platform/lib/protobuf v0.0.0 => ./.protobuf
