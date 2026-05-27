package generator

import "testing"

func TestNormalizeRichness(t *testing.T) {
	got, err := NormalizeRichness(" Rich ")
	if err != nil {
		t.Fatal(err)
	}
	if got != "rich" {
		t.Fatalf("richness = %q", got)
	}
	got, err = NormalizeRichness("")
	if err != nil {
		t.Fatal(err)
	}
	if got != DefaultRichness {
		t.Fatalf("default richness = %q", got)
	}
}

func TestNormalizeRichnessRejectsInvalidValue(t *testing.T) {
	if _, err := NormalizeRichness("busy"); err == nil {
		t.Fatal("expected error")
	}
}
