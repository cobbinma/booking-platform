package auth0_test

import (
	"context"
	"fmt"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/gateway_api/internal/auth0"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	token = "test-token"
	email = "test@test.com"
	name  = "Test Test"
)

func Test_userService_GetUser(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer != token {
			t.Fatalf("bearer token = '%s', expected = '%s'", bearer, token)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"email":"%s","name":"%s"}`, email, name)))
	}))
	defer authServer.Close()
	ctx := auth0.AddTokenToCtx(context.Background(), token)
	us := auth0.NewUserService(authServer.URL)
	got, err := us.GetUser(ctx)
	if err != nil {
		t.Errorf("did not expect error, got '%s'", err)
		return
	}

	cupaloy.SnapshotT(t, got)
}
