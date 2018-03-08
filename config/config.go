package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config contains the global configuration of this CLI
type Config struct {
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

// Conf is the global configuration shared across the app
var Conf = &Config{}

var errNotConfigured = errors.New("aws-creds hasn't been configured yet")

var jsonMarshalIndent = json.MarshalIndent

// Load data from the config file into the Config struct
func (conf *Config) Load(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return errNotConfigured
	}
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if string(raw) == "" {
		return errNotConfigured
	}

	return json.Unmarshal(raw, &conf)
}

// Save data from the Config struct into the config file
func (conf *Config) Save(configFile string) error {
	path := filepath.Dir(configFile)
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}

	raw, err := jsonMarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, raw, 0644)
}
