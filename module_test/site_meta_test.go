package module_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	cfbase "github.com/btnguyen2k/cf-base"
)

func TestGetDescriptionMap(t *testing.T) {
	tests := []struct {
		name        string
		description any
		want        map[string]string
	}{
		{
			name:        "nil",
			description: nil,
			want:        map[string]string{},
		},
		{
			name:        "single string",
			description: "A website",
			want:        map[string]string{"en": "A website"},
		},
		{
			name: "localized map",
			description: map[string]any{
				"en": "A website",
				"vi": "Mot trang web",
			},
			want: map[string]string{
				"en": "A website",
				"vi": "Mot trang web",
			},
		},
		{
			name:        "unsupported type",
			description: 42,
			want:        map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &cfbase.SiteMeta{
				DefaultLanguage: "en",
				Description:     tt.description,
			}
			if got := meta.GetDescriptionMap(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetDescriptionMap() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestGetTagAliasMap(t *testing.T) {
	tests := []struct {
		name    string
		aliases any
		want    map[string]map[string][]string
	}{
		{
			name:    "nil",
			aliases: nil,
			want:    map[string]map[string][]string{},
		},
		{
			name: "flat aliases",
			aliases: map[string]any{
				"go":   []any{"golang", "go-lang"},
				"docs": []string{"documentation"},
			},
			want: map[string]map[string][]string{
				"en": {
					"go":   {"golang", "go-lang"},
					"docs": {"documentation"},
				},
			},
		},
		{
			name: "localized aliases",
			aliases: map[string]any{
				"en": map[string]any{"go": []any{"golang"}},
				"vi": map[string][]string{"go": {"go-lang"}},
			},
			want: map[string]map[string][]string{
				"en": {"go": {"golang"}},
				"vi": {"go": {"go-lang"}},
			},
		},
		{
			name:    "unsupported type",
			aliases: map[string]any{"go": "golang"},
			want:    map[string]map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &cfbase.SiteMeta{
				DefaultLanguage: "en",
				TagsAlias:       tt.aliases,
			}
			if got := meta.GetTagAliasMap(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetTagAliasMap() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestGetTagAliasMapSkipsNilWhenDetectingShape(t *testing.T) {
	meta := &cfbase.SiteMeta{
		DefaultLanguage: "en",
		TagsAlias: map[string]any{
			"unset": nil,
			"go":    []string{"golang"},
		},
	}

	for i := 0; i < 10; i++ {
		got := meta.GetTagAliasMap()
		if aliases := got["en"]["go"]; !reflect.DeepEqual(aliases, []string{"golang"}) {
			t.Fatalf("GetTagAliasMap() iteration %d returned %#v", i, got)
		}
	}
}

func TestSiteMetaToMap(t *testing.T) {
	author := &cfbase.Author{Name: "Alice", Email: "alice@example.com"}
	meta := &cfbase.SiteMeta{
		Name:            "Example",
		Languages:       map[string]string{"en": "English"},
		DefaultLanguage: "en",
		Icon:            "icon.png",
		Description:     "A website",
		Contacts:        map[string]string{"email": "contact@example.com"},
		Tags:            map[string]any{"topic": "go"},
		Mode:            cfbase.SiteModeBlog,
		Author:          author,
	}
	want := map[string]any{
		"name":            "Example",
		"languages":       meta.Languages,
		"defaultLanguage": "en",
		"icon":            "icon.png",
		"description":     map[string]string{"en": "A website"},
		"contacts":        meta.Contacts,
		"tags":            meta.Tags,
		"mode":            cfbase.SiteModeBlog,
		"author":          author,
	}

	if got := meta.ToMap(); !reflect.DeepEqual(got, want) {
		t.Fatalf("ToMap() = %#v, want %#v", got, want)
	}
}

func TestLoadSiteMetaFromTestData(t *testing.T) {
	meta, err := cfbase.LoadSiteMetaFromYaml(filepath.Join("test_data", "meta.yaml"))
	if err != nil {
		t.Fatalf("LoadSiteMetaFromYaml() returned error: %v", err)
	}

	if meta.Name != "ContentFlow" {
		t.Errorf("Name = %q, want ContentFlow", meta.Name)
	}
	if meta.Icon != "fas fa-code" {
		t.Errorf("Icon = %q, want %q", meta.Icon, "fas fa-code")
	}
	if meta.DefaultLanguage != "en" {
		t.Errorf("DefaultLanguage = %q, want en", meta.DefaultLanguage)
	}
	if meta.Mode != cfbase.DefaultSiteMode {
		t.Errorf("Mode = %q, want %q", meta.Mode, cfbase.DefaultSiteMode)
	}

	wantLanguages := map[string]string{
		"default": "en",
		"en":      "English",
		"vi":      "Tiếng Việt",
	}
	if !reflect.DeepEqual(meta.Languages, wantLanguages) {
		t.Errorf("Languages = %#v, want %#v", meta.Languages, wantLanguages)
	}

	wantDescriptions := map[string]string{
		"en": "Content Management System where its content is built through CI/CD pipeline",
		"vi": "Hệ thống Quản trị nội dung với dữ liệu được xây dựng thông qua qui trình CI/CD",
	}
	if got := meta.GetDescriptionMap(); !reflect.DeepEqual(got, wantDescriptions) {
		t.Errorf("GetDescriptionMap() = %#v, want %#v", got, wantDescriptions)
	}

	wantContacts := map[string]string{
		"website":  "https://github.com/btnguyen2k/docms",
		"email":    "btnguyen2k (at) gmail (dot) com",
		"github":   "https://github.com/btnguyen2k/",
		"linkedin": "https://www.linkedin.com/in/btnguyen2k/",
	}
	if !reflect.DeepEqual(meta.Contacts, wantContacts) {
		t.Errorf("Contacts = %#v, want %#v", meta.Contacts, wantContacts)
	}

	if got := meta.Tags["build"]; got != "${build_datetime}" {
		t.Errorf("Tags[build] = %#v, want %q", got, "${build_datetime}")
	}
	demo, ok := meta.Tags["demo"].(map[string]interface{})
	if !ok {
		t.Fatalf("Tags[demo] has type %T, want map[string]interface{}", meta.Tags["demo"])
	}
	if got := demo["tag3"]; got != "this tag _content_ has **markdown**" {
		t.Errorf("Tags[demo][tag3] = %#v", got)
	}

	aliases := meta.GetTagAliasMap()
	if got := aliases["en"]["cms"]; !reflect.DeepEqual(got, []string{
		"content management",
		"content management system",
		"contentflow",
	}) {
		t.Errorf("English cms aliases = %#v", got)
	}
	if got := aliases["en"]["contentflow runtime"]; !reflect.DeepEqual(got, []string{"runtime"}) {
		t.Errorf("English ContentFlow runtime aliases = %#v", got)
	}
	if got := aliases["vi"]["cms"]; !reflect.DeepEqual(got, []string{
		"quản trị nội dung",
		"hệ thống quản trị nội dung",
		"quản lý nội dung",
		"hệ thống quản lý nội dung",
		"contentflow",
	}) {
		t.Errorf("Vietnamese cms aliases = %#v", got)
	}
	if got := aliases["vi"]["contentflow runtime"]; !reflect.DeepEqual(got, []string{"runtime"}) {
		t.Errorf("Vietnamese ContentFlow runtime aliases = %#v", got)
	}
	if got := aliases["vi"]["toàn văn"]; !reflect.DeepEqual(got, []string{
		"fti",
		"fulltext index",
		"full text index",
		"index",
		"tìm kiếm",
		"tìm kiếm toàn văn",
		"toàn văn",
		"chỉ mục",
		"chỉ mục toàn văn",
		"toàn văn chỉ mục",
	}) {
		t.Errorf("Vietnamese full-text aliases = %#v", got)
	}

	asMap := meta.ToMap()
	if got := asMap["description"]; !reflect.DeepEqual(got, wantDescriptions) {
		t.Errorf("ToMap()[description] = %#v, want %#v", got, wantDescriptions)
	}
	author, ok := asMap["author"].(*cfbase.Author)
	if !ok || author != nil {
		t.Errorf("ToMap()[author] = %#v, want a nil *Author", asMap["author"])
	}
}

func TestLoadSiteMetaAutoFromTestData(t *testing.T) {
	meta, err := cfbase.LoadSiteMetaAuto("test_data")
	if err != nil {
		t.Fatalf("LoadSiteMetaAuto() returned error: %v", err)
	}
	if meta.Name != "ContentFlow" {
		t.Fatalf("Name = %q, want ContentFlow", meta.Name)
	}
}

func TestLoadSiteMeta(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		content  string
		load     func(string) (*cfbase.SiteMeta, error)
		wantName string
		wantLang string
		wantMode string
	}{
		{
			name: "yaml",
			ext:  ".yaml",
			content: `
name: Example YAML
description: A YAML website
languages:
  en: English
mode: invalid
contacts:
  email: contact@example.com
  phone: ""
`,
			load:     cfbase.LoadSiteMetaFromYaml,
			wantName: "Example YAML",
			wantLang: "en",
			wantMode: cfbase.DefaultSiteMode,
		},
		{
			name: "json",
			ext:  ".json",
			content: `{
				"name": "Example JSON",
				"description": "A JSON website",
				"languages": {
					"default": "vi",
					"en": "English",
					"vi": "Vietnamese"
				},
				"mode": "blog",
				"contacts": {
					"email": "contact@example.com",
					"phone": ""
				}
			}`,
			load:     cfbase.LoadSiteMetaFromJson,
			wantName: "Example JSON",
			wantLang: "vi",
			wantMode: cfbase.SiteModeBlog,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTestFile(t, "meta"+tt.ext, tt.content)
			meta, err := tt.load(path)
			if err != nil {
				t.Fatalf("load returned error: %v", err)
			}

			if meta.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", meta.Name, tt.wantName)
			}
			if meta.DefaultLanguage != tt.wantLang {
				t.Errorf("DefaultLanguage = %q, want %q", meta.DefaultLanguage, tt.wantLang)
			}
			if meta.Mode != tt.wantMode {
				t.Errorf("Mode = %q, want %q", meta.Mode, tt.wantMode)
			}
			if !reflect.DeepEqual(meta.Contacts, map[string]string{"email": "contact@example.com"}) {
				t.Errorf("Contacts = %#v", meta.Contacts)
			}
			if meta.Tags == nil || len(meta.Tags) != 0 {
				t.Errorf("Tags = %#v, want initialized empty map", meta.Tags)
			}
		})
	}
}

func TestLoadSiteMetaValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		ext     string
		content string
		load    func(string) (*cfbase.SiteMeta, error)
		wantErr string
	}{
		{
			name:    "empty YAML",
			ext:     ".yaml",
			content: "",
			load:    cfbase.LoadSiteMetaFromYaml,
			wantErr: "empty or null",
		},
		{
			name:    "null JSON",
			ext:     ".json",
			content: "null",
			load:    cfbase.LoadSiteMetaFromJson,
			wantErr: "empty or null",
		},
		{
			name:    "malformed YAML",
			ext:     ".yaml",
			content: "languages: [",
			load:    cfbase.LoadSiteMetaFromYaml,
			wantErr: "did not find expected node content",
		},
		{
			name:    "malformed JSON",
			ext:     ".json",
			content: "{",
			load:    cfbase.LoadSiteMetaFromJson,
			wantErr: "unexpected end of JSON input",
		},
		{
			name:    "no languages",
			ext:     ".yaml",
			content: "name: Example",
			load:    cfbase.LoadSiteMetaFromYaml,
			wantErr: "at least one language must be configured",
		},
		{
			name: "multiple languages without default",
			ext:  ".yaml",
			content: `
languages:
  en: English
  vi: Vietnamese
`,
			load:    cfbase.LoadSiteMetaFromYaml,
			wantErr: "languages.default is required",
		},
		{
			name: "unknown default language",
			ext:  ".json",
			content: `{
				"languages": {
					"default": "fr",
					"en": "English"
				}
			}`,
			load:    cfbase.LoadSiteMetaFromJson,
			wantErr: `default language "fr" is not configured`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTestFile(t, "meta"+tt.ext, tt.content)
			meta, err := tt.load(path)
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

func TestLoadSiteMetaMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	for _, load := range []func(string) (*cfbase.SiteMeta, error){
		cfbase.LoadSiteMetaFromYaml,
		cfbase.LoadSiteMetaFromJson,
	} {
		if meta, err := load(missing); !os.IsNotExist(err) || meta != nil {
			t.Errorf("load(%q) = (%#v, %v), want (nil, not-exist error)", missing, meta, err)
		}
	}
}

func TestLoadSiteMetaAuto(t *testing.T) {
	t.Run("prefers YAML over JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.yaml"), `
name: YAML
languages:
  en: English
`)
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{
			"name": "JSON",
			"languages": {"en": "English"}
		}`)

		meta, err := cfbase.LoadSiteMetaAuto(dir)
		if err != nil {
			t.Fatalf("LoadSiteMetaAuto() returned error: %v", err)
		}
		if meta.Name != "YAML" {
			t.Fatalf("Name = %q, want YAML", meta.Name)
		}
	})

	t.Run("falls back to JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{
			"name": "JSON",
			"languages": {"en": "English"}
		}`)

		meta, err := cfbase.LoadSiteMetaAuto(dir)
		if err != nil {
			t.Fatalf("LoadSiteMetaAuto() returned error: %v", err)
		}
		if meta.Name != "JSON" {
			t.Fatalf("Name = %q, want JSON", meta.Name)
		}
	})

	t.Run("returns invalid preferred file error", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.yaml"), "languages: [")
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{
			"name": "JSON",
			"languages": {"en": "English"}
		}`)

		if meta, err := cfbase.LoadSiteMetaAuto(dir); err == nil || meta != nil {
			t.Fatalf("LoadSiteMetaAuto() = (%#v, %v), want YAML error", meta, err)
		}
	})

	t.Run("no metadata file", func(t *testing.T) {
		if meta, err := cfbase.LoadSiteMetaAuto(t.TempDir()); err == nil || meta != nil {
			t.Fatalf("LoadSiteMetaAuto() = (%#v, %v), want not-found error", meta, err)
		}
	})
}

func writeTestFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	writeFileAt(t, path, content)
	return path
}

func writeFileAt(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}
