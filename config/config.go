package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config contains the configuration of this CLI.
type Config struct {
	Username            string     `json:"username"`
	OktaHost            string     `json:"okta_host,omitempty"`
	OktaAppPath         string     `json:"okta_app_path,omitempty"`
	PreferredFactorType string     `json:"preferred_factor_type,omitempty"`
	Profiles            []*Profile `json:"profiles"`
	filepath            string
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
	return &Config{filepath: path}
}

// Load loads data from the config file into the Config struct.
func (c *Config) Load() error {
	if _, err := os.Stat(c.filepath); os.IsNotExist(err) {
		return errNotConfigured
	}
	raw, err := ioutil.ReadFile(c.filepath)
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
	path := filepath.Dir(c.filepath)
	err := os.MkdirAll(path, directoryPermissions)
	if err != nil {
		return err
	}

	raw, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.filepath, raw, filePermissions)
	if err != nil {
		return err
	}
	fmt.Println("Configuration saved!")
	return nil
}
