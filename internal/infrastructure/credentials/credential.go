package credentials

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

func ValidateToken(token string, ctx context.Context) error {
	token = strings.Replace(token, "Bearer ", "", 1)
	provider, err := oidc.NewProvider(ctx, os.Getenv("KEYCLOAK"))
	if err != nil {
		return errors.New("error to connect the provider")
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: "emailn"})
	_, err = verifier.Verify(ctx, token)
	if err != nil {
		return errors.New("invalid token")
	}
	return nil
}
