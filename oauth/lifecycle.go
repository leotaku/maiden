package oauth

import (
	"fmt"
	"path"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Lifecycle struct {
	ConfigFile string
	TokenFile  string
	ScopesFile string
}

func NewLifecycle(directory string) *Lifecycle {
	return &Lifecycle{
		ConfigFile: path.Join(directory, "credentials.json"),
		TokenFile:  path.Join(directory, "token.json"),
		ScopesFile: path.Join(directory, "scopes.json"),
	}
}

func (l *Lifecycle) Init(credentials string, scopes ...string) (*oauth2.Config, *oauth2.Token, error) {
	config, err := google.ConfigFromJSON([]byte(credentials), scopes...)
	if err != nil {
		return nil, nil, fmt.Errorf("config: %w", err)
	}
	token, err := LoadTokenFromWeb(*config)
	if err != nil {
		return nil, nil, fmt.Errorf("web auth: %w", err)
	}

	// Only write after success
	if err := WriteToken(l.TokenFile, token); err != nil {
		return nil, nil, fmt.Errorf("write token: %w", err)
	}
	if err := WriteConfig(l.ConfigFile, credentials); err != nil {
		return nil, nil, fmt.Errorf("write config: %w", err)
	}
	if err := WriteScopes(l.ScopesFile, scopes...); err != nil {
		return nil, nil, fmt.Errorf("write scopes: %w", err)
	}

	return config, token, nil
}

func (l *Lifecycle) Load() (*oauth2.Config, *oauth2.Token, error) {
	scopes, err := LoadScopes(l.ScopesFile)
	if err != nil {
		return nil, nil, fmt.Errorf("scopes: %w", err)
	}
	config, err := LoadConfig(l.ConfigFile, scopes...)
	if err != nil {
		return nil, nil, fmt.Errorf("config: %w", err)
	}
	token, err := LoadToken(l.TokenFile)
	if err != nil {
		return nil, nil, fmt.Errorf("token: %w", err)
	}

	return config, token, nil
}
