package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Configuration contains the global configuration of this CLI
type Configuration struct {
	Username            string     `json:"username"`
	OktaOrgURL          string     `json:"okta_org_url"`
	PreferredFactorType string     `json:"preferred_factor_type"`
	Profiles            []*Profile `json:"profiles"`
}

// Profile contains the configuration of each AWS profile
type Profile struct {
	Name    string `json:"name"`
	RoleArn string `json:"role_arn"`
}

// Config is the global configuration shared across the app
var Config = &Configuration{}

// ErrNotConfigured is thrown if this CLI hasn't been configured yet
var ErrNotConfigured = errors.New("aws-creds hasn't been configured yet")

var jsonMarshalIndent = json.MarshalIndent

// Load data from the config file into the Configuration struct
func (config *Configuration) Load(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return ErrNotConfigured
	}
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if string(raw) == "" {
		return ErrNotConfigured
	}

	return json.Unmarshal(raw, &config)
}

// Save data from the Configuration struct into the config file
func (config *Configuration) Save(configFile string) error {
	path := filepath.Dir(configFile)
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}

	raw, err := jsonMarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, raw, 0644)
}
