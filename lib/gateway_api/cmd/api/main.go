package main

import (
	mw "github.com/cobbinma/booking-platform/lib/gateway_api/cmd/api/middleware"
	"github.com/cobbinma/booking-platform/lib/gateway_api/internal/auth0"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
)

const defaultPort = "9999"

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	domain, present := os.LookupEnv("AUTH0_DOMAIN")
	if !present {
		panic("AUTH0_DOMAIN environment variable not set")
	}
	apiId, present := os.LookupEnv("AUTH0_API_IDENTIFIER")
	if !present {
		panic("AUTH0_API_IDENTIFIER environment variable not set")
	}
	venueURL, present := os.LookupEnv("VENUE_API_ROOT")
	if !present {
		panic("VENUE_API_ROOT environment variable not set")
	}

	conn, err := grpc.Dial(venueURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect : %s", err)
	}
	defer conn.Close()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{Resolvers: graph.NewResolver(auth0.NewUserService(domain), api.NewVenueAPIClient(conn))}))
	e := echo.New()
	e.Use(middleware.Logger())

	_, present = os.LookupEnv("ALLOW_CORS_URL")
	if present {
		e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	}

	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	e.POST("/query", echo.WrapHandler(srv), mw.Auth(domain, apiId))
	e.OPTIONS("/query", PreflightCors())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	e.Logger.Fatal(e.Start(":" + port))
}

func PreflightCors() echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := c.Request().Header
		for key, value := range headers {
			c.Response().Header().Set(key, value[0])
		}
		return c.NoContent(http.StatusOK)
	}
}
