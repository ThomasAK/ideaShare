package auth

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"ideashare/config"
)

const IdeaShareIDToken = "ideaShare_id_token"

type OidcUserClaims struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}

type OidcTokenVerifier struct {
	*oidc.IDTokenVerifier
}

func (v *OidcTokenVerifier) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	return v.IDTokenVerifier.Verify(ctx, rawIDToken)
}

func SetUpOIDC() (*oidc.Provider, *oauth2.Config, config.IDTokenVerifier) {
	provider, err := oidc.NewProvider(context.Background(), config.GetStringOr(config.OIDCProviderUrl, "http://localhost:8747/realms/master"))
	if err != nil {
		panic(err)
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     config.GetStringOr(config.OIDCClientId, "IdeaShare"),
		ClientSecret: config.GetStringOr(config.OIDCClientSecret, "8pWPZek1fRltycAoV6ITqZw3FGbB8L74"),
		RedirectURL:  config.GetStringOr(config.OIDCCallbackUrl, "http://localhost:5173/api/auth/authorize"),

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	return provider, &oauth2Config, &OidcTokenVerifier{provider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID})}
}
