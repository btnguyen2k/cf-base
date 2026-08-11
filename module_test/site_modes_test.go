package module_test

import (
	"reflect"
	"testing"

	cfbase "github.com/btnguyen2k/cf-base"
)

func TestIsValidSiteMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want bool
	}{
		{name: "document", mode: cfbase.SiteModeDocument, want: true},
		{name: "blog", mode: cfbase.SiteModeBlog, want: true},
		{name: "empty", mode: "", want: false},
		{name: "unknown", mode: "unknown", want: false},
		{name: "case sensitive", mode: "Document", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cfbase.IsValidSiteMode(tt.mode); got != tt.want {
				t.Fatalf("IsValidSiteMode(%q) = %v, want %v", tt.mode, got, tt.want)
			}
		})
	}
}

func TestSiteModes(t *testing.T) {
	want := []string{cfbase.SiteModeDocument, cfbase.SiteModeBlog}
	if got := cfbase.SiteModes(); !reflect.DeepEqual(got, want) {
		t.Fatalf("SiteModes() = %#v, want %#v", got, want)
	}
}

func TestSiteModesReturnsIndependentSlice(t *testing.T) {
	modes := cfbase.SiteModes()
	modes[0] = "changed"

	want := []string{cfbase.SiteModeDocument, cfbase.SiteModeBlog}
	if got := cfbase.SiteModes(); !reflect.DeepEqual(got, want) {
		t.Fatalf("SiteModes() after caller mutation = %#v, want %#v", got, want)
	}
}

func TestDefaultSiteMode(t *testing.T) {
	if cfbase.DefaultSiteMode != cfbase.SiteModeDocument {
		t.Fatalf("DefaultSiteMode = %q, want %q", cfbase.DefaultSiteMode, cfbase.SiteModeDocument)
	}
	if !cfbase.IsValidSiteMode(cfbase.DefaultSiteMode) {
		t.Fatalf("DefaultSiteMode %q is not valid", cfbase.DefaultSiteMode)
	}
}
