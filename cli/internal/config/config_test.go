package config

import "testing"

func TestModelsRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Models == nil {
		t.Fatal("Models should be initialized")
	}
	for _, role := range Roles {
		if !ValidRole(role) {
			t.Errorf("role %q should be valid", role)
		}
	}
	if ValidRole("nope") {
		t.Error("nope should not be a valid role")
	}
	cfg.Models[RoleThinking] = "sonnet"
	if err := Save(dir, cfg); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.Models[RoleThinking] != "sonnet" {
		t.Fatalf("models = %v", got.Models)
	}
}
