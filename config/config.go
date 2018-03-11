package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config contains the configuration of this CLI.
type Config struct {
	Filepath            string
	Username            string     `json:"username"`
	OktaOrgURL          string     `json:"okta_org_url"`
	PreferredFactorType string     `json:"preferred_factor_type"`
	Profiles            []*Profile `json:"profiles"`
}

// Profile contains the configuration of each AWS profile.
type Profile struct {
	Name    string `json:"name"`
	RoleARN string `json:"role_arn"`
}

const (
	directoryPermissions = 0700
	filePermissions      = 0644
)

var errNotConfigured = errors.New("aws-creds hasn't been configured yet")

// New creates a new Config reference with the given filepath.
func New(path string) *Config {
	return &Config{Filepath: path}
}

// Load loads data from the config file into the Config struct.
func (c *Config) Load() error {
	if _, err := os.Stat(c.Filepath); os.IsNotExist(err) {
		return errNotConfigured
	}
	raw, err := ioutil.ReadFile(c.Filepath)
	if err != nil {
		return err
	}
	if string(raw) == "" {
		return errNotConfigured
	}

	return json.Unmarshal(raw, &c)
}

// Save saves data from the Config struct into the config file.
func (c *Config) Save() error {
	path := filepath.Dir(c.Filepath)
	err := os.MkdirAll(path, directoryPermissions)
	if err != nil {
		return err
	}

	raw, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Filepath, raw, filePermissions)
}
