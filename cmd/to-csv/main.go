package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/mail"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ConfigGroup struct {
	ID         string             `yaml:"id"`
	Name       string             `yaml:"name"`
	Recipients []*ConfigRecipient `yaml:"recipients"`
}

type ConfigRecipient struct {
	Email       string   `yaml:"email"`
	CC          []string `yaml:"cc,omitempty"`
	OpeningLine string   `yaml:"opening_line"`
}

type FileConfig struct {
	SecretKey      string         `yaml:"secret_key"`
	PublicHost     string         `yaml:"public_host"`
	HTTPOnly       bool           `yaml:"http_only"`
	GoogleClientID string         `yaml:"google_client_id"`
	GoogleSecret   string         `yaml:"google_secret"`
	Groups         []*ConfigGroup `yaml:"groups"`
	Port           *int           `yaml:"port"`
	Title          string         `yaml:"title"`

	// For development; ignore Google authentication.
	NoGoogleAuth bool `yaml:"no_google_auth"`

	// For TLS configuration.
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`

	// Should be a string like "google4f9d0c78202b2454.html". If non-empty and
	// not starting with "google", it will be prepended. If it does not end with
	// ".html", ".html" will be appended.
	GoogleSiteVerification string `yaml:"google_site_verification"`
}

func loadConfig(filename string) (*FileConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := new(FileConfig)
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{
		"ID", "Name", "Email", "Address",
	})
	for _, group := range cfg.Groups {
		for _, recipient := range group.Recipients {
			addr, err := mail.ParseAddress(recipient.Email)
			if err != nil {
				log.Fatal(err)
			}
			if err := w.Write([]string{
				group.ID,
				addr.Name,
				addr.Address,
				recipient.OpeningLine,
			}); err != nil {
				log.Fatal(err)
			}
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
