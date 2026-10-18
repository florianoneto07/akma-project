package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func base() Config {
	return Config{Configuration: Configuration{
		FQDN: "af1.akma.example", UaStarID: "FF00000001", Listen: "192.168.16.129:8080",
		AanfURI: "http://192.168.16.130:8000/",
	}}
}

func TestDefaults(t *testing.T) {
	s, err := base().Settings()
	if err != nil {
		t.Fatal(err)
	}
	if s.AfID.Transport() != "af1.akma.example.FF00000001" || s.Realm != "3GPP-bootstrapping-akma@af1.akma.example" ||
		s.Path != DefaultPath || len(s.Qop) != 1 || s.Qop[0] != "auth" || s.NonceLifetime != DefaultNonceLifetime ||
		s.KafExpiry != KafExpiryRefetch || s.AanfURI != "http://192.168.16.130:8000" {
		t.Fatalf("settings %+v", s)
	}
}

func TestInvalid(t *testing.T) {
	for name, mutate := range map[string]func(*Configuration){
		"fqdn with port":      func(c *Configuration) { c.FQDN = "af1.akma.example:8080" },
		"fqdn empty label":    func(c *Configuration) { c.FQDN = "af1..example" },
		"fqdn non ascii":      func(c *Configuration) { c.FQDN = "äf1.example" },
		"ua* id short":        func(c *Configuration) { c.UaStarID = "FF0000" },
		"listen without port": func(c *Configuration) { c.Listen = "192.168.16.129" },
		"aanfUri":             func(c *Configuration) { c.AanfURI = "192.168.16.130:8000" },
		"qop":                 func(c *Configuration) { c.Qop = []string{"auth-conf"} },
		"nonceLifetime":       func(c *Configuration) { c.NonceLifetime = "-1s" },
		"kafExpiry":           func(c *Configuration) { c.KafExpiry = "renew" },
		"realmPrefix with @":  func(c *Configuration) { c.RealmPrefix = "a@b" },
		"path":                func(c *Configuration) { c.Path = "ua" },
		"clockOffset":         func(c *Configuration) { c.ClockOffset = "two hours" },
	} {
		t.Run(name, func(t *testing.T) {
			c := base()
			mutate(&c.Configuration)
			if _, err := c.Settings(); err == nil {
				t.Fatal("accepted")
			}
		})
	}
}

// TestSampleConfig loads the configuration of the laboratory.
func TestClockOffset(t *testing.T) {
	c := base()
	c.Configuration.ClockOffset = "-2h"
	s, err := c.Settings()
	if err != nil || s.ClockOffset != -2*time.Hour {
		t.Fatalf("settings %+v, %v", s, err)
	}
}

func TestSampleConfig(t *testing.T) {
	s, err := Load("../../../config/afcfg.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if s.AfID.Transport() != "af1.akma.example.FF00000001" || s.Listen != "192.168.16.129:8080" ||
		s.NonceLifetime != 300*time.Second || s.KafExpiry != KafExpiryRefetch || !s.Logger.Enable {
		t.Fatalf("settings %+v", s)
	}
}

func TestUnknownField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "afcfg.yaml")
	raw := "configuration:\n  fqdn: af1.akma.example\n  uaStarID: FF00000001\n"
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "uaStarID") {
		t.Fatalf("err = %v", err)
	}
}
