package config

import (
	"fmt"
)

type Config struct {
	ServerURL     string
	Realm         string
	ClientID      string
	AdminClientID string
	ClientSecret  string
	AdminUser     string
	AdminPass     string
	RealmURL      string
	AdminRealmURL string
	TokenURL      string
	UserURL       string
}

func GetConfig() *Config {
	config := &Config{
		ServerURL:     "http://conducirya.com.ar:18080",
		Realm:         "DriveFluency",
		ClientID:      "drivefluency",
		ClientSecret:  "083E22w85Iw9T2vctotLkT3ZAEDaqXsA",
		AdminClientID: "admin-cli",
		AdminUser:     "drivefluency@gmail.com",
		AdminPass:     "admin",
	}

	// URLs de Keycloack
	config.RealmURL = fmt.Sprintf("%s/realms/%s", config.ServerURL, config.Realm)
	config.AdminRealmURL = fmt.Sprintf("%s/admin/realms/%s", config.ServerURL, config.Realm)
	config.TokenURL = fmt.Sprintf("%s/protocol/openid-connect/token", config.RealmURL)
	config.UserURL = fmt.Sprintf("%s/users", config.AdminRealmURL)

	return config
}
