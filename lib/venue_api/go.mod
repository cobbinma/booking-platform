module github.com/cobbinma/booking-platform/lib/venue_api

go 1.14

require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/auth0-community/go-auth0 v1.0.1-0.20191119091237-b9b0f95be568
	github.com/bradleyjkemp/cupaloy v1.3.0
	github.com/cobbinma/booking-platform/lib/protobuf v0.0.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/jmoiron/sqlx v1.3.1
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.9.0
	github.com/ory/dockertest/v3 v3.6.3
	github.com/stretchr/testify v1.5.1
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.35.0
	gopkg.in/square/go-jose.v2 v2.5.1
)

replace github.com/cobbinma/booking-platform/lib/protobuf v0.0.0 => ./.protobuf
