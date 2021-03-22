package auth0_test

import (
	"context"
	"fmt"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/gateway_api/internal/auth0"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	token      = "test-token"
	email      = "test@test.com"
	familyName = "Test"
	givenName  = "Test"
)

func Test_userService_GetUser(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer != token {
			t.Fatalf("bearer token = '%s', expected = '%s'", bearer, token)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"email":"%s","given_name":"%s","family_name":"%s"}`, email, givenName, familyName)))
	}))
	defer authServer.Close()
	ctx := models.AddTokenToCtx(context.Background(), token)
	us := auth0.NewUserService(authServer.URL)
	got, err := us.GetUser(ctx)
	if err != nil {
		t.Errorf("did not expect error, got '%s'", err)
		return
	}

	cupaloy.SnapshotT(t, got)
}
