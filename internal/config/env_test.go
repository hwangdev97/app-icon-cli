package config

import "testing"

func TestResolveAPIKeyFromEnvReference(t *testing.T) {
	t.Setenv("APPICON_TEST_KEY", "secret")
	got, err := ResolveAPIKey("env:APPICON_TEST_KEY")
	if err != nil {
		t.Fatal(err)
	}
	if got != "secret" {
		t.Fatalf("key = %q", got)
	}
}

func TestResolveAPIKeyMissingEnv(t *testing.T) {
	_, err := ResolveAPIKey("env:APPICON_MISSING_KEY")
	if err == nil {
		t.Fatal("expected error")
	}
}
