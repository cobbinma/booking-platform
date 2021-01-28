package main

import (
	"crypto/tls"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/venue_api/cmd/api/middleware"
	"github.com/cobbinma/booking-platform/lib/venue_api/internal/postgres"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
)

func main() {
	_ = godotenv.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("could not start logger : %s", err)
		os.Exit(-1)
	}
	defer func(log *zap.Logger) {
		if err := log.Sync(); err != nil {
			log.Error("could not sync logger", zap.Error(err))
		}
	}(logger)
	log := logger.Sugar()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	db, closeDB, err := postgres.NewPostgres(log)
	if err != nil {
		log.Fatalf("could not construct postgres : %s", err)
	}
	defer closeDB(log)

	ensureValidToken, err := middleware.EnsureValidToken()
	if err != nil {
		log.Fatalf("could not construct auth server inceptor : %s", err)
	}

	cert, err := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
	if err != nil {
		log.Fatalf("failed to load cert : %s", err)
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("could not listen : %s", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc_middleware.WithUnaryServerChain(grpc_zap.UnaryServerInterceptor(logger), ensureValidToken),
	}

	s := grpc.NewServer(opts...)
	api.RegisterVenueAPIServer(s, db)
	api.RegisterTableAPIServer(s, db)

	log.Infof("starting gRPC listener on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %s", err)
	}
}
