module github.com/cobbinma/booking-platform/lib/gateway_api

go 1.15

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/auth0-community/go-auth0 v1.0.0
	github.com/bradleyjkemp/cupaloy v2.3.0+incompatible
	github.com/cobbinma/booking-platform/lib/protobuf v0.0.0
	github.com/joho/godotenv v1.3.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo/v4 v4.1.17
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/vektah/gqlparser/v2 v2.1.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/grpc v1.35.0
	gopkg.in/square/go-jose.v2 v2.1.7
)

replace github.com/cobbinma/booking-platform/lib/protobuf v0.0.0 => ./.protobuf
