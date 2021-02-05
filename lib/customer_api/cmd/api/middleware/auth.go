package middleware

import (
	"context"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"os"
	"strings"
)

func EnsureValidToken() (grpc.UnaryServerInterceptor, error) {
	baseURL := os.Getenv("AUTH0_DOMAIN")
	if baseURL == "" {
		return nil, fmt.Errorf("missing auth0 domain environment variable")
	}

	if len(baseURL) > 0 && baseURL[len(baseURL)-1] != '/' {
		baseURL = baseURL + "/"
	}
	apiIdentifier := os.Getenv("AUTH0_API_IDENTIFIER")
	if apiIdentifier == "" {
		return nil, fmt.Errorf("missing api identifier environment variable")
	}
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: fmt.Sprintf("%s.well-known/jwks.json", baseURL)}, nil)
	validator := auth0.NewValidator(auth0.NewConfiguration(client, []string{apiIdentifier}, baseURL, jose.RS256), nil)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		if len(md["authorization"]) != 1 {
			return nil, status.Errorf(codes.InvalidArgument, "auth token has incorrect length")
		}

		token, err := jwt.ParseSigned(strings.TrimPrefix(md["authorization"][0], "Bearer"))
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "could not parse token : %s", err)
		}
		if err = validator.ValidateToken(token); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "could not validate token : %s", err)
		}

		return handler(ctx, req)
	}, nil
}
