package module_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	cfbase "github.com/btnguyen2k/cf-base"
)

func TestDocumentMetaTextMaps(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
		get  func(*cfbase.DocumentMeta) map[string]string
	}{
		{
			name: "single title",
			want: map[string]string{"en": "What it is"},
			get: func(meta *cfbase.DocumentMeta) map[string]string {
				meta.Title = "What it is"
				return meta.GetTitleMap()
			},
		},
		{
			name: "localized title",
			want: map[string]string{"en": "What it is", "vi": "La gi"},
			get: func(meta *cfbase.DocumentMeta) map[string]string {
				meta.Title = map[string]any{"en": "What it is", "vi": "La gi"}
				return meta.GetTitleMap()
			},
		},
		{
			name: "single summary",
			want: map[string]string{"en": "Introduction"},
			get: func(meta *cfbase.DocumentMeta) map[string]string {
				meta.Summary = "Introduction"
				return meta.GetSummaryMap()
			},
		},
		{
			name: "localized content file",
			want: map[string]string{"en": "index-en.md", "vi": "index-vi.md"},
			get: func(meta *cfbase.DocumentMeta) map[string]string {
				meta.ContentFile = map[string]string{"en": "index-en.md", "vi": "index-vi.md"}
				return meta.GetContentFileMap()
			},
		},
		{
			name: "unsupported value",
			want: map[string]string{},
			get: func(meta *cfbase.DocumentMeta) map[string]string {
				meta.Title = 42
				return meta.GetTitleMap()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &cfbase.DocumentMeta{DefaultLanguage: "en"}
			if got := tt.get(meta); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("text map = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDocumentMetaGetTagsMap(t *testing.T) {
	tests := []struct {
		name string
		tags any
		want map[string][]string
	}{
		{
			name: "flat tags",
			tags: []any{"search", "cms", "ci/cd"},
			want: map[string][]string{"en": {"ci/cd", "cms", "search"}},
		},
		{
			name: "localized tags",
			tags: map[string]any{
				"en": []any{"search", "cms"},
				"vi": []string{"tìm kiếm", "nội dung"},
			},
			want: map[string][]string{
				"en": {"cms", "search"},
				"vi": {"nội dung", "tìm kiếm"},
			},
		},
		{
			name: "nil tags",
			tags: nil,
			want: map[string][]string{},
		},
		{
			name: "unsupported tags",
			tags: "cms",
			want: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &cfbase.DocumentMeta{
				DefaultLanguage: "en",
				Tags:            tt.tags,
			}
			if got := meta.GetTagsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetTagsMap() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDocumentMetaToMap(t *testing.T) {
	author := &cfbase.Author{Name: "Alice"}
	meta := &cfbase.DocumentMeta{
		DefaultLanguage: "en",
		Title:           "What it is",
		Summary:         "Introduction",
		Icon:            "document.svg",
		ContentFile:     "index.md",
		Tags:            []string{"cms"},
		EntryImage:      "entry.png",
		DocPage:         "about",
		DocStyle:        "wide",
		TimestampCreate: 100,
		TimestampUpdate: 200,
		Author:          author,
	}
	want := map[string]any{
		"id":      "",
		"icon":    "document.svg",
		"title":   map[string]string{"en": "What it is"},
		"summary": map[string]string{"en": "Introduction"},
		"file":    map[string]string{"en": "index.md"},
		"tags":    map[string][]string{"en": {"cms"}},
		"img":     "entry.png",
		"page":    "about",
		"style":   "wide",
		"tc":      int64(100),
		"tu":      int64(200),
		"author":  author,
	}

	if got := meta.ToMap(); !reflect.DeepEqual(got, want) {
		t.Fatalf("ToMap() = %#v, want %#v", got, want)
	}
}

func TestLoadDocumentMetaFromTestData(t *testing.T) {
	meta, err := cfbase.LoadDocumentMetaFromYaml(
		filepath.Join("test_data", "01-intro", "01-whatitis", "meta.yaml"),
	)
	if err != nil {
		t.Fatalf("LoadDocumentMetaFromYaml() returned error: %v", err)
	}

	wantTitles := map[string]string{
		"en": "What it is",
		"vi": "ContentFlow là gì",
	}
	if got := meta.GetTitleMap(); !reflect.DeepEqual(got, wantTitles) {
		t.Errorf("GetTitleMap() = %#v, want %#v", got, wantTitles)
	}

	wantSummaries := map[string]string{
		"en": "ContentFlow is a Content Management System that helps authors publish website content through a CI/CD flow. Unlike other CMS, there is no UI to create, update and publish content in ContentFlow. Instead, website content is built and published via CI/CD pipelines.",
		"vi": "ContentFlow là Hệ thống quản lý nội dung giúp tác giả xuất bản nội dung trang web thông qua luồng CI/CD. Sẽ không có giao diện để người dùng tạo, cập nhật và xuất bản nội dung lên trang web. Thay vào đó, nội dung của trang web sẽ được xây dựng và xuất bản thông qua qui trình CI/CD.",
	}
	if got := meta.GetSummaryMap(); !reflect.DeepEqual(got, wantSummaries) {
		t.Errorf("GetSummaryMap() = %#v, want %#v", got, wantSummaries)
	}

	wantFiles := map[string]string{
		"en": "index-en.md",
		"vi": "index-vi.md",
	}
	if got := meta.GetContentFileMap(); !reflect.DeepEqual(got, wantFiles) {
		t.Errorf("GetContentFileMap() = %#v, want %#v", got, wantFiles)
	}

	wantTags := map[string][]string{
		"en": {"ci/cd", "cms", "content management"},
		"vi": {"ci/cd", "cms", "quản lý nội dung", "quản trị nội dung"},
	}
	if got := meta.GetTagsMap(); !reflect.DeepEqual(got, wantTags) {
		t.Errorf("GetTagsMap() = %#v, want %#v", got, wantTags)
	}

	if meta.Icon != "" || meta.EntryImage != "" || meta.DocPage != "" || meta.DocStyle != "" {
		t.Errorf("optional presentation fields = %#v", meta)
	}
	if meta.TimestampCreate != 0 || meta.TimestampUpdate != 0 || meta.Author != nil {
		t.Errorf("optional metadata fields = %#v", meta)
	}

	asMap := meta.ToMap()
	if got := asMap["title"]; !reflect.DeepEqual(got, wantTitles) {
		t.Errorf("ToMap()[title] = %#v, want %#v", got, wantTitles)
	}
	if got := asMap["summary"]; !reflect.DeepEqual(got, wantSummaries) {
		t.Errorf("ToMap()[summary] = %#v, want %#v", got, wantSummaries)
	}
	if got := asMap["file"]; !reflect.DeepEqual(got, wantFiles) {
		t.Errorf("ToMap()[file] = %#v, want %#v", got, wantFiles)
	}
	if got := asMap["tags"]; !reflect.DeepEqual(got, wantTags) {
		t.Errorf("ToMap()[tags] = %#v, want %#v", got, wantTags)
	}
}

func TestLoadDocumentMetaAutoFromTestData(t *testing.T) {
	meta, err := cfbase.LoadDocumentMetaAuto(
		filepath.Join("test_data", "01-intro", "01-whatitis"),
	)
	if err != nil {
		t.Fatalf("LoadDocumentMetaAuto() returned error: %v", err)
	}
	if got := meta.GetTitleMap()["en"]; got != "What it is" {
		t.Fatalf("English title = %q, want %q", got, "What it is")
	}
}

func TestLoadDocumentMetaErrors(t *testing.T) {
	tests := []struct {
		name    string
		ext     string
		content string
		load    func(string) (*cfbase.DocumentMeta, error)
		wantErr string
	}{
		{
			name:    "empty YAML",
			ext:     ".yaml",
			load:    cfbase.LoadDocumentMetaFromYaml,
			wantErr: "empty or null",
		},
		{
			name:    "null JSON",
			ext:     ".json",
			content: "null",
			load:    cfbase.LoadDocumentMetaFromJson,
			wantErr: "empty or null",
		},
		{
			name:    "malformed YAML",
			ext:     ".yaml",
			content: "title: [",
			load:    cfbase.LoadDocumentMetaFromYaml,
			wantErr: "did not find expected node content",
		},
		{
			name:    "malformed JSON",
			ext:     ".json",
			content: "{",
			load:    cfbase.LoadDocumentMetaFromJson,
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

func TestLoadDocumentMetaMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	for _, load := range []func(string) (*cfbase.DocumentMeta, error){
		cfbase.LoadDocumentMetaFromYaml,
		cfbase.LoadDocumentMetaFromJson,
	} {
		if meta, err := load(missing); !os.IsNotExist(err) || meta != nil {
			t.Errorf("load(%q) = (%#v, %v), want (nil, not-exist error)", missing, meta, err)
		}
	}
}

func TestLoadDocumentMetaAuto(t *testing.T) {
	t.Run("prefers YAML over JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.yaml"), "title: YAML")
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{"title":"JSON"}`)

		meta, err := cfbase.LoadDocumentMetaAuto(dir)
		if err != nil {
			t.Fatalf("LoadDocumentMetaAuto() returned error: %v", err)
		}
		if meta.Title != "YAML" {
			t.Fatalf("Title = %#v, want YAML", meta.Title)
		}
	})

	t.Run("falls back to JSON", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{"title":"JSON"}`)

		meta, err := cfbase.LoadDocumentMetaAuto(dir)
		if err != nil {
			t.Fatalf("LoadDocumentMetaAuto() returned error: %v", err)
		}
		if meta.Title != "JSON" {
			t.Fatalf("Title = %#v, want JSON", meta.Title)
		}
	})

	t.Run("returns invalid preferred file error", func(t *testing.T) {
		dir := t.TempDir()
		writeFileAt(t, filepath.Join(dir, "meta.yaml"), "title: [")
		writeFileAt(t, filepath.Join(dir, "meta.json"), `{"title":"JSON"}`)

		if meta, err := cfbase.LoadDocumentMetaAuto(dir); err == nil || meta != nil {
			t.Fatalf("LoadDocumentMetaAuto() = (%#v, %v), want YAML error", meta, err)
		}
	})

	t.Run("no metadata file", func(t *testing.T) {
		if meta, err := cfbase.LoadDocumentMetaAuto(t.TempDir()); err == nil || meta != nil {
			t.Fatalf("LoadDocumentMetaAuto() = (%#v, %v), want not-found error", meta, err)
		}
	})
}
