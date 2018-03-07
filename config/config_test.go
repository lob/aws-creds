package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testConfig struct {
	label     string
	setupFunc func(*Config, string)
	shouldErr bool
}

var (
	s       = rand.NewSource(time.Now().UnixNano())
	r       = rand.New(s)
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func TestLoad(t *testing.T) {
	tests := []testConfig{
		{
			"Success",
			func(cfg *Config, configFile string) {
				path := filepath.Dir(configFile)
				_ = os.MkdirAll(path, 0700)
				_ = ioutil.WriteFile(configFile, []byte("{}"), 0644)
			},
			false,
		},
		{
			"NotExistError",
			func(cfg *Config, configFile string) {},
			true,
		},
		{
			"NoReadError",
			func(cfg *Config, configFile string) {
				path := filepath.Dir(configFile)
				_ = os.MkdirAll(path, 0700)
				_ = ioutil.WriteFile(configFile, []byte("{}"), 0222)
			},
			true,
		},
		{
			"EmptyFileError",
			func(cfg *Config, configFile string) {
				path := filepath.Dir(configFile)
				_ = os.MkdirAll(path, 0700)
				_ = ioutil.WriteFile(configFile, []byte(""), 0644)
			},
			true,
		},
	}

	for _, tt := range tests {
		cfg := &Config{}
		configFile := fmt.Sprintf("/tmp/aws-creds-%s/config", randStr())
		tt.setupFunc(cfg, configFile)
		err := cfg.Load(configFile)
		if tt.shouldErr && err == nil {
			t.Errorf("%s: expected error", tt.label)
		} else if !tt.shouldErr && err != nil {
			t.Errorf("%s: unexpected error: %s", tt.label, err)
		}
		path := filepath.Dir(configFile)
		_ = os.RemoveAll(path)
	}
}

func TestSave(t *testing.T) {
	origMarshal := jsonMarshalIndent
	defer func() {
		jsonMarshalIndent = origMarshal
	}()

	tests := []testConfig{
		{
			"Success",
			func(cfg *Config, configFile string) {
				jsonMarshalIndent = origMarshal
			},
			false,
		},
		{
			"NoPermissionsError",
			func(cfg *Config, configFile string) {
				jsonMarshalIndent = origMarshal
				path := filepath.Dir(configFile)
				parentPath := filepath.Dir(path)
				_ = os.MkdirAll(parentPath, 0200)
			},
			true,
		},
		{
			"JSONMarshalError",
			func(cfg *Config, configFile string) {
				jsonMarshalIndent = func(i interface{}, a, b string) ([]byte, error) {
					return nil, errors.New("err")
				}
			},
			true,
		},
	}

	for _, tt := range tests {
		cfg := &Config{}
		configFile := fmt.Sprintf("/tmp/aws-creds-%s/aws-creds/config", randStr())
		tt.setupFunc(cfg, configFile)
		err := cfg.Save(configFile)
		if tt.shouldErr && err == nil {
			t.Errorf("%s: expected error", tt.label)
		} else if !tt.shouldErr && err != nil {
			t.Errorf("%s: unexpected error: %s", tt.label, err)
		}
		path := filepath.Dir(configFile)
		parentPath := filepath.Dir(path)
		_ = os.RemoveAll(parentPath)
	}
}

func randStr() string {
	n := 16
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
