package akamai

import "net/http"

type Session struct {
	apiKey string
	client *http.Client
}

// NewSession creates a new Session that can be used to make requests to the FRIS Solutions API.
func NewSession(apiKey string) *Session {
	return &Session{
		apiKey: apiKey,
		client: http.DefaultClient,
	}
}

// WithClient sets a new client that will be used to make requests to the Hyper Solutions API.
func (s *Session) WithClient(client *http.Client) *Session {
	s.client = client
	return s
}
