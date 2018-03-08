package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	test "github.com/lob/aws-creds/testing"
)

type testConfig struct {
	setupFunc func(*testing.T, *Config, string) func()
	shouldErr bool
}

func TestLoad(t *testing.T) {
	cases := map[string]testConfig{
		"Success": {
			func(t *testing.T, conf *Config, configFile string) func() {
				path := filepath.Dir(configFile)
				if err := os.MkdirAll(path, 0700); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if err := ioutil.WriteFile(configFile, []byte("{}"), 0644); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return func() {
					if err := os.RemoveAll(path); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
				}
			},
			false,
		},
		"NotExistError": {
			func(t *testing.T, conf *Config, configFile string) func() {
				return func() {}
			},
			true,
		},
		"NoReadError": {
			func(t *testing.T, conf *Config, configFile string) func() {
				path := filepath.Dir(configFile)
				if err := os.MkdirAll(path, 0700); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if err := ioutil.WriteFile(configFile, []byte("{}"), 0222); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return func() {
					if err := os.RemoveAll(path); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
				}
			},
			true,
		},
		"EmptyFileError": {
			func(t *testing.T, conf *Config, configFile string) func() {
				path := filepath.Dir(configFile)
				if err := os.MkdirAll(path, 0700); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if err := ioutil.WriteFile(configFile, []byte(""), 0644); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return func() {
					if err := os.RemoveAll(path); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
				}
			},
			true,
		},
	}

	for key, tc := range cases {
		conf := &Config{}
		configFile := fmt.Sprintf("/tmp/aws-creds-%s/config", test.RandStr(16))
		defer tc.setupFunc(t, conf, configFile)()
		err := conf.Load(configFile)
		if tc.shouldErr && err == nil {
			t.Errorf("%s: expected error", key)
		} else if !tc.shouldErr && err != nil {
			t.Errorf("%s: unexpected error: %s", key, err)
		}
	}
}

func TestSave(t *testing.T) {
	origMarshal := jsonMarshalIndent
	defer func() {
		jsonMarshalIndent = origMarshal
	}()

	cases := map[string]testConfig{
		"Success": {
			func(t *testing.T, conf *Config, configFile string) func() {
				jsonMarshalIndent = origMarshal
				return func() {
					path := filepath.Dir(configFile)
					parentPath := filepath.Dir(path)
					if err := os.RemoveAll(parentPath); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
				}
			},
			false,
		},
		"NoPermissionsError": {
			func(t *testing.T, conf *Config, configFile string) func() {
				jsonMarshalIndent = origMarshal
				path := filepath.Dir(configFile)
				parentPath := filepath.Dir(path)
				if err := os.MkdirAll(parentPath, 0200); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return func() {
					path := filepath.Dir(configFile)
					parentPath := filepath.Dir(path)
					if err := os.RemoveAll(parentPath); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
				}
			},
			true,
		},
		"JSONMarshalError": {
			func(t *testing.T, conf *Config, configFile string) func() {
				jsonMarshalIndent = func(i interface{}, a, b string) ([]byte, error) {
					return nil, errors.New("err")
				}
				return func() {
					path := filepath.Dir(configFile)
					parentPath := filepath.Dir(path)
					if err := os.RemoveAll(parentPath); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
				}
			},
			true,
		},
	}

	for key, tc := range cases {
		conf := &Config{}
		configFile := fmt.Sprintf("/tmp/aws-creds-%s/aws-creds/config", test.RandStr(16))
		defer tc.setupFunc(t, conf, configFile)()
		err := conf.Save(configFile)
		if tc.shouldErr && err == nil {
			t.Errorf("%s: expected error", key)
		} else if !tc.shouldErr && err != nil {
			t.Errorf("%s: unexpected error: %s", key, err)
		}
	}
}
