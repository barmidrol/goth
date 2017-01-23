// Package heroku implements the OAuth2 protocol for authenticating users through heroku.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package procore

import (
	// "encoding/json"
	"io"
	"net/http"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

const (
	authURL         string = "https://app.procore.com/oauth/authorize"
	tokenURL        string = "https://app.procore.com/oauth/token"
)

// Provider is the implementation of `goth.Provider` for accessing Heroku.
type Provider struct {
	ClientKey          string
	ClientSecret      string
	CallbackURL string
	HTTPClient  *http.Client
	config      *oauth2.Config
}

func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:   clientKey,
		ClientSecret:  secret,
		CallbackURL: callbackURL,
	}
	p.config = newConfig(p, scopes)
	return p
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return "procore"
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

func (p *Provider) Debug(debug bool) {}

func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	return &Session{
		AuthURL: p.config.AuthCodeURL(state),
	}, nil
}

// FetchUser will go to Heroku and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	s := session.(*Session)
	user := goth.User{
		AccessToken:  s.AccessToken,
		Provider:     p.Name(),
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
	}
	// req, err := http.NewRequest("GET", endpointProfile, nil)
	// if err != nil {
	// 	return user, err
	// }
  //
	// req.Header.Set("Authorization", "Bearer "+s.AccessToken)
  //
	// resp, err := p.Client().Do(req)
	// if err != nil {
	// 	if resp != nil {
	// 		resp.Body.Close()
	// 	}
	// 	return user, err
	// }
	// defer resp.Body.Close()
  //
	// err = userFromReader(resp.Body, &user)
	return user, nil
}

func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	if len(scopes) > 0 {
		for _, scope := range scopes {
			c.Scopes = append(c.Scopes, scope)
		}
	}
	return c
}

func userFromReader(r io.Reader, user *goth.User) error {
	// u := struct {
	// 	Name  string `json:"name"`
	// 	Email string `json:"email"`
	// 	ID    string `json:"id"`
	// }{}
	// err := json.NewDecoder(r).Decode(&u)
	// if err != nil {
	// 	return err
	// }
	// user.Email = u.Email
	// user.Name = u.Name
	// user.UserID = u.ID
	return nil
}

//RefreshTokenAvailable refresh token is provided by auth provider or not
func (p *Provider) RefreshTokenAvailable() bool {
	return true
}

//RefreshToken get new access token based on the refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.config.TokenSource(goth.ContextForClient(p.Client()), token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}
