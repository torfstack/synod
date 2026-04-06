package http

import "github.com/torfstack/synod/backend/config"

// testConfig returns a minimal config suitable for unit tests.
func testConfig() config.Config {
	return config.Config{
		Auth: config.AuthConfig{
			Issuer:       "https://auth.example.com",
			ClientID:     "test-client",
			ClientSecret: "test-secret",
			RedirectURL:  "https://app.example.com/callback",
		},
		Server: config.ServerConfig{
			Port:    8080,
			BaseURL: "https://app.example.com",
		},
	}
}
