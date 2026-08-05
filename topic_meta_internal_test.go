package cfbase

import "testing"

func TestTopicMetaSetDirectory(t *testing.T) {
	tests := []struct {
		name      string
		dir       string
		wantValid bool
		wantIndex int
		wantID    string
	}{
		{
			name:      "valid",
			dir:       "001-getting-started",
			wantValid: true,
			wantIndex: 1,
			wantID:    "getting-started",
		},
		{
			name:      "valid underscore",
			dir:       "12-go_basics",
			wantValid: true,
			wantIndex: 12,
			wantID:    "go_basics",
		},
		{
			name: "missing index",
			dir:  "getting-started",
		},
		{
			name: "missing identifier",
			dir:  "1-",
		},
		{
			name: "invalid separator",
			dir:  "1_getting-started",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &TopicMeta{}
			if got := meta.setDirectory(tt.dir); got != tt.wantValid {
				t.Fatalf("setDirectory(%q) = %v, want %v", tt.dir, got, tt.wantValid)
			}
			if meta.dir != tt.dir {
				t.Errorf("dir = %q, want %q", meta.dir, tt.dir)
			}
			if meta.index != tt.wantIndex {
				t.Errorf("index = %d, want %d", meta.index, tt.wantIndex)
			}
			if meta.id != tt.wantID {
				t.Errorf("id = %q, want %q", meta.id, tt.wantID)
			}
		})
	}
}
