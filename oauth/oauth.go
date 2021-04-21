package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const modeDefaultPerm = 0600

func LoadConfig(path string, scope ...string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	config, err := google.ConfigFromJSON(b, scope...)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	return config, nil
}

func LoadTokenFromWeb(config oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Authorize this app by visiting this url:", authURL)
	fmt.Print("Enter the code from that page here: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("authorization code: %w", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("web token: %w", err)
	}

	return tok, nil
}

func LoadToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	defer f.Close()
	tok := &oauth2.Token{}
	return tok, json.NewDecoder(f).Decode(tok)
}

func LoadScopes(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	defer f.Close()
	v := make([]string, 0)
	if err := json.NewDecoder(f).Decode(&v); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return v, nil
}

func WriteConfig(filename, text string) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return fmt.Errorf("write directory: %w", err)
	}
	if err := ioutil.WriteFile(filename, []byte(text), modeDefaultPerm); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func WriteToken(filename string, tok *oauth2.Token) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return fmt.Errorf("write directory: %w", err)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, modeDefaultPerm)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(tok)
}

func WriteScopes(filename string, scopes ...string) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return fmt.Errorf("write directory: %w", err)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, modeDefaultPerm)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(scopes)
}
