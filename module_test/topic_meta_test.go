package module_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	cfbase "github.com/btnguyen2k/cf-base"
)

func TestTopicMetaTextMaps(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		want        map[string]string
		description bool
	}{
		{
			name:  "nil title",
			value: nil,
			want:  map[string]string{},
		},
		{
			name:  "single title",
			value: "Getting started",
			want:  map[string]string{"en": "Getting started"},
		},
		{
			name:  "localized title",
			value: map[string]any{"en": "Getting started", "vi": "Bat dau"},
			want:  map[string]string{"en": "Getting started", "vi": "Bat dau"},
		},
		{
			name:  "unsupported title",
			value: 42,
			want:  map[string]string{},
		},
		{
			name:        "single description",
			value:       "Introduction",
			want:        map[string]string{"en": "Introduction"},
			description: true,
		},
		{
			name:        "localized description",
			value:       map[string]string{"en": "Introduction", "vi": "Gioi thieu"},
			want:        map[string]string{"en": "Introduction", "vi": "Gioi thieu"},
			description: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &cfbase.TopicMeta{DefaultLanguage: "en"}
			var got map[string]string
			if tt.description {
				meta.Description = tt.value
				got = meta.GetDescriptionMap()
			} else {
				meta.Title = tt.value
				got = meta.GetTitleMap()
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("text map = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestTopicMetaToMap(t *testing.T) {
	meta := &cfbase.TopicMeta{
		DefaultLanguage: "en",
		Title:           "Getting started",
		Description:     map[string]string{"en": "Introduction"},
		Icon:            "topic.svg",
		EntryImage:      "entry.png",
		Hidden:          true,
	}
	want := map[string]interface{}{
		"id":          "",
		"num_docs":    0,
		"icon":        "topic.svg",
		"title":       map[string]string{"en": "Getting started"},
		"description": map[string]string{"en": "Introduction"},
		"img":         "entry.png",
		"hidden":      true,
	}

	if got := meta.ToMap(); !reflect.DeepEqual(got, want) {
		t.Fatalf("ToMap() = %#v, want %#v", got, want)
	}
}

func TestLoadTopicMetaFromTestData(t *testing.T) {
	meta, err := cfbase.LoadTopicMetaFromYaml(filepath.Join("test_data", "01-intro", "meta.yaml"))
	if err != nil {
		t.Fatalf("LoadTopicMetaFromYaml() returned error: %v", err)
	}

	wantTitles := map[string]string{
		"en": "Introduction",
		"vi": "Giới thiệu",
	}
	if got := meta.GetTitleMap(); !reflect.DeepEqual(got, wantTitles) {
		t.Errorf("GetTitleMap() = %#v, want %#v", got, wantTitles)
	}

	wantDescriptions := map[string]string{
		"en": "An introduction of DO CMS: what it is and how it work.",
		"vi": "Giới thiệu về DO CMS: tổng quan và cách thức hoạt động.",
	}
	if got := meta.GetDescriptionMap(); !reflect.DeepEqual(got, wantDescriptions) {
		t.Errorf("GetDescriptionMap() = %#v, want %#v", got, wantDescriptions)
	}

	if meta.Icon != "fas fa-lightbulb" {
		t.Errorf("Icon = %q, want %q", meta.Icon, "fas fa-lightbulb")
	}
	if meta.EntryImage != "" {
		t.Errorf("EntryImage = %q, want empty", meta.EntryImage)
	}
	if meta.Hidden {
		t.Error("Hidden = true, want false")
	}

	wantMap := map[string]interface{}{
		"id":          "",
		"num_docs":    0,
		"icon":        "fas fa-lightbulb",
		"title":       wantTitles,
		"description": wantDescriptions,
		"img":         "",
		"hidden":      false,
	}
	if got := meta.ToMap(); !reflect.DeepEqual(got, wantMap) {
		t.Errorf("ToMap() = %#v, want %#v", got, wantMap)
	}
}

func TestLoadTopicMetaAutoFromTestData(t *testing.T) {
	meta, err := cfbase.LoadTopicMetaAuto(filepath.Join("test_data", "01-intro"))
	if err != nil {
		t.Fatalf("LoadTopicMetaAuto() returned error: %v", err)
	}
	if got := meta.GetTitleMap()["en"]; got != "Introduction" {
		t.Fatalf("English title = %q, want Introduction", got)
	}
}

func TestLoadTopicMeta(t *testing.T) {
	tests := []struct {
		name    string
		ext     string
		content string
		load    func(string) (*cfbase.TopicMeta, error)
	}{
		{
			name: "YAML",
			ext:  ".yaml",
			content: `
def_lang: en
title: Getting started
description:
  en: Introduction
icon: topic.svg
img: entry.png
hidden: true
`,
			load: cfbase.LoadTopicMetaFromYaml,
		},
		{
			name: "JSON",
			ext:  ".json",
			content: `{
				"def_lang": "en",
				"title": "Getting started",
				"description": {"en": "Introduction"},
				"icon": "topic.svg",
				"img": "entry.png",
				"hidden": true
			}`,
			load: cfbase.LoadTopicMetaFromJson,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := tt.load(writeTestFile(t, "meta"+tt.ext, tt.content))
			if err != nil {
				t.Fatalf("load returned error: %v", err)
			}
			if meta.DefaultLanguage != "en" {
				t.Errorf("DefaultLanguage = %q, want en", meta.DefaultLanguage)
			}
			if !reflect.DeepEqual(meta.GetTitleMap(), map[string]string{"en": "Getting started"}) {
				t.Errorf("title = %#v", meta.GetTitleMap())
			}
			if !reflect.DeepEqual(meta.GetDescriptionMap(), map[string]string{"en": "Introduction"}) {
				t.Errorf("description = %#v", meta.GetDescriptionMap())
			}
			if meta.Icon != "topic.svg" || meta.EntryImage != "entry.png" || !meta.Hidden {
				t.Errorf("metadata = %#v", meta)
			}
		})
	}
}

func TestLoadTopicMetaErrors(t *testing.T) {
	tests := []struct {
		name    string
		ext     string
		content string
		load    func(string) (*cfbase.TopicMeta, error)
		wantErr string
	}{
		{
			name:    "empty YAML",
			ext:     ".yaml",
			load:    cfbase.LoadTopicMetaFromYaml,
			wantErr: "empty or null",
		},
		{
			name:    "null JSON",
			ext:     ".json",
			content: "null",
			load:    cfbase.LoadTopicMetaFromJson,
			wantErr: "empty or null",
		},
		{
			name:    "malformed YAML",
			ext:     ".yaml",
			content: "title: [",
			load:    cfbase.LoadTopicMetaFromYaml,
			wantErr: "did not find expected node content",
		},
		{
			name:    "malformed JSON",
			ext:     ".json",
			content: "{",
			load:    cfbase.LoadTopicMetaFromJson,
			wantErr: "unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := tt.load(writeTestFile(t, "meta"+tt.ext, tt.content))
			if err == nil {
				t.Fatalf("load returned metadata %#v, want error containing %q", meta, tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want substring %q", err, tt.wantErr)
			}
			if meta != nil {
				t.Fatalf("metadata = %#v, want nil", meta)
			}
		})
	}
}

func TestLoadTopicMetaMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	for _, load := range []func(string) (*cfbase.TopicMeta, error){
		cfbase.LoadTopicMetaFromYaml,
		cfbase.LoadTopicMetaFromJson,
	} {
		if meta, err := load(missing); !os.IsNotExist(err) || meta != nil {
			t.Errorf("load(%q) = (%#v, %v), want (nil, not-exist error)", missing, meta, err)
		}
	}
}

func TestLoadTopicMetaAuto(t *testing.T) {
	t.Run("prefers YAML over JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.yaml"), "title: YAML")
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{"title":"JSON"}`)

		meta, err := cfbase.LoadTopicMetaAuto(dir)
		if err != nil {
			t.Fatalf("LoadTopicMetaAuto() returned error: %v", err)
		}
		if meta.Title != "YAML" {
			t.Fatalf("Title = %#v, want YAML", meta.Title)
		}
	})

	t.Run("falls back to JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{"title":"JSON"}`)

		meta, err := cfbase.LoadTopicMetaAuto(dir)
		if err != nil {
			t.Fatalf("LoadTopicMetaAuto() returned error: %v", err)
		}
		if meta.Title != "JSON" {
			t.Fatalf("Title = %#v, want JSON", meta.Title)
		}
	})

	t.Run("returns invalid preferred file error", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.yaml"), "title: [")
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{"title":"JSON"}`)

		if meta, err := cfbase.LoadTopicMetaAuto(dir); err == nil || meta != nil {
			t.Fatalf("LoadTopicMetaAuto() = (%#v, %v), want YAML error", meta, err)
		}
	})

	t.Run("no metadata file", func(t *testing.T) {
		if meta, err := cfbase.LoadTopicMetaAuto(t.TempDir()); err == nil || meta != nil {
			t.Fatalf("LoadTopicMetaAuto() = (%#v, %v), want not-found error", meta, err)
		}
	})
}
