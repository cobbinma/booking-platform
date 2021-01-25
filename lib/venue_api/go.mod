module github.com/cobbinma/booking-platform/lib/venue_api

go 1.14

require (
	github.com/cobbinma/booking-platform/lib/protobuf v0.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/sirupsen/logrus v1.6.0
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	google.golang.org/grpc v1.35.0
)

replace github.com/cobbinma/booking-platform/lib/protobuf v0.0.0 => ./.protobuf
