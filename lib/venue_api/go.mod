module github.com/cobbinma/booking-platform/lib/venue_api

go 1.14

require (
	github.com/cobbinma/booking-platform/lib/protobuf v0.0.0 // indirect
	github.com/labstack/echo/v4 v4.1.17
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	google.golang.org/grpc v1.35.0
)

replace github.com/cobbinma/booking-platform/lib/protobuf v0.0.0 => ./.protobuf
