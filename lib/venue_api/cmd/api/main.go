package main

import (
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("could not listen : %s", err)
	}

	s := grpc.NewServer()
	api.RegisterVenueAPIServer(s, &api.UnimplementedVenueAPIServer{})

	log.Infof("starting gRPC listener on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %s", err)
	}
}
