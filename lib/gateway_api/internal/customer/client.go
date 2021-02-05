package customer

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/customer/api"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func NewCustomerClient(url string, log *zap.SugaredLogger, token *oauth2.Token) (graph.CustomerService, func(log *zap.SugaredLogger), error) {
	creds, err := credentials.NewClientTLSFromFile("localhost.crt", "localhost")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load credentials : %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(token)),
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect : %s", err)
	}

	return &customerClient{
			client: api.NewCustomerAPIClient(conn),
			log:    log,
		}, func(log *zap.SugaredLogger) {
			if err := conn.Close(); err != nil {
				log.Error("could not close connection : %s", err)
			}
		}, nil
}

type customerClient struct {
	client api.CustomerAPIClient
	log    *zap.SugaredLogger
}

func (c customerClient) IsAdmin(ctx context.Context, venueID string, email string) (bool, error) {
	resp, err := c.client.IsAdmin(ctx, &api.IsAdminRequest{
		VenueId: venueID,
		Email:   email,
	})
	if err != nil {
		return false, fmt.Errorf("could not get is admin from customer client : %w", err)
	}

	return resp.IsAdmin, nil
}
