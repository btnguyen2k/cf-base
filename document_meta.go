package cfbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/btnguyen2k/consu/reddo"
	"gopkg.in/yaml.v3"
)

// DocumentMeta describes a document and its presentation metadata.
type DocumentMeta struct {
	fileInfo os.FileInfo `json:"-" yaml:"-"`
	index    int         `json:"-" yaml:"-"`
	id       string      `json:"-" yaml:"-"`
	dir      string      `json:"-" yaml:"-"`

	// DefaultLanguage is the default language inherited from the containing site.
	DefaultLanguage string `json:"def_lang" yaml:"def_lang"`
	// Title is either a string or a map of language codes to titles.
	Title any `json:"title" yaml:"title"`
	// Summary is either a string or a map of language codes to summaries.
	Summary any `json:"summary,omitempty" yaml:"summary,omitempty"`
	// Icon is the location of the document icon.
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty"`
	// ContentFile is either a filename or a map of language codes to filenames.
	ContentFile any `json:"file" yaml:"file"`
	// Tags is either a tag list or a map of language codes to tag lists.
	Tags any `json:"tags,omitempty" yaml:"tags,omitempty"`
	// EntryImage is the location of the document's entry image.
	EntryImage string `json:"img,omitempty" yaml:"img,omitempty"`
	// DocPage identifies a special site page, such as contact or about.
	DocPage string `json:"page,omitempty" yaml:"page,omitempty"`
	// DocStyle identifies the document's presentation style.
	DocStyle string `json:"style,omitempty" yaml:"style,omitempty"`
	// TimestampCreate is the Unix timestamp when the document was created.
	TimestampCreate int64 `json:"tc" yaml:"tc"`
	// TimestampUpdate is the Unix timestamp when the document was last updated.
	TimestampUpdate int64 `json:"tu" yaml:"tu"`
	// Author is the document's author.
	Author *Author `json:"author,omitempty" yaml:"author,omitempty"`
}

// func (dm *DocumentMeta) setDirectory(dir string) bool {
// 	dm.dir = dir
// 	if !_rexpContentDir.MatchString(dir) {
// 		return false
// 	}
// 	matches := _rexpContentDir.FindStringSubmatch(dir)
// 	dm.index, _ = strconv.Atoi(matches[1])
// 	dm.id = matches[2]
// 	return true
// }

// ToMap returns a map representation of the document metadata.
func (dm *DocumentMeta) ToMap() map[string]any {
	return map[string]any{
		"id":      dm.id,
		"icon":    dm.Icon,
		"title":   dm.GetTitleMap(),
		"summary": dm.GetSummaryMap(),
		"file":    dm.GetContentFileMap(),
		"tags":    dm.GetTagsMap(),
		"img":     dm.EntryImage,
		"page":    dm.DocPage,
		"style":   dm.DocStyle,
		"tc":      dm.TimestampCreate,
		"tu":      dm.TimestampUpdate,
		"author":  dm.Author,
	}
}

// GetSummaryMap returns the document summary keyed by language code.
func (dm *DocumentMeta) GetSummaryMap() map[string]string {
	summary := make(map[string]string)
	if dm.Summary != nil {
		switch reflect.TypeOf(dm.Summary).Kind() {
		case reflect.String:
			summary[dm.DefaultLanguage] = fmt.Sprintf("%s", dm.Summary)
		case reflect.Map:
			temp, err := reddo.Convert(dm.Summary, _typMapString)
			if err == nil && temp != nil {
				summary = temp.(map[string]string)
			}
		}
	}
	return summary
}

// GetTitleMap returns the document title keyed by language code.
func (dm *DocumentMeta) GetTitleMap() map[string]string {
	title := make(map[string]string)
	if dm.Title != nil {
		switch reflect.TypeOf(dm.Title).Kind() {
		case reflect.String:
			title[dm.DefaultLanguage] = fmt.Sprintf("%s", dm.Title)
		case reflect.Map:
			temp, err := reddo.Convert(dm.Title, _typMapString)
			if err == nil && temp != nil {
				title = temp.(map[string]string)
			}
		}
	}
	return title
}

// GetContentFileMap returns the document content filename keyed by language code.
func (dm *DocumentMeta) GetContentFileMap() map[string]string {
	files := make(map[string]string)
	if dm.ContentFile != nil {
		switch reflect.TypeOf(dm.ContentFile).Kind() {
		case reflect.String:
			files[dm.DefaultLanguage] = fmt.Sprintf("%s", dm.ContentFile)
		case reflect.Map:
			temp, err := reddo.Convert(dm.ContentFile, _typMapString)
			if err == nil && temp != nil {
				files = temp.(map[string]string)
			}
		}
	}
	return files
}

// GetTagsMap returns sorted document tags keyed by language code.
func (dm *DocumentMeta) GetTagsMap() map[string][]string {
	tagsMap := make(map[string][]string)
	if dm.Tags != nil {
		switch reflect.TypeOf(dm.Tags).Kind() {
		case reflect.Array, reflect.Slice:
			temp, err := reddo.Convert(dm.Tags, reflect.TypeOf([]string{}))
			if err == nil && temp != nil {
				tagsMap[dm.DefaultLanguage] = temp.([]string)
			}
		case reflect.Map:
			temp, err := reddo.Convert(dm.Tags, reflect.TypeOf(map[string][]string{}))
			if err == nil && temp != nil {
				tagsMap = temp.(map[string][]string)
			}
		}
	}
	for k := range tagsMap {
		sort.Strings(tagsMap[k])
	}
	return tagsMap
}

// LoadDocumentMetaAuto loads the first available meta.yaml, meta.yml, or meta.json file in dir.
func LoadDocumentMetaAuto(dir string) (*DocumentMeta, error) {
	yamlFiles := []string{dir + "/meta.yaml", dir + "/meta.yml"}
	for _, yamlFilePath := range yamlFiles {
		documentMeta, err := LoadDocumentMetaFromYaml(yamlFilePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return documentMeta, nil
	}

	jsonFiles := []string{dir + "/meta.json"}
	for _, jsonFilePath := range jsonFiles {
		documentMeta, err := LoadDocumentMetaFromJson(jsonFilePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return documentMeta, nil
	}

	return nil, fmt.Errorf("no meta file found at '%s', '%s' or '%s'", yamlFiles[0], yamlFiles[1], jsonFiles[0])
}

// LoadDocumentMetaFromYaml loads document metadata from a YAML file.
func LoadDocumentMetaFromYaml(filePath string) (*DocumentMeta, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var metadata *DocumentMeta
	if err := yaml.Unmarshal(buf, &metadata); err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, fmt.Errorf("document metadata file %q is empty or null", filePath)
	}

	metadata.fileInfo = fi
	return metadata, nil
}

// LoadDocumentMetaFromJson loads document metadata from a JSON file.
func LoadDocumentMetaFromJson(filePath string) (*DocumentMeta, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var metadata *DocumentMeta
	if err := json.Unmarshal(buf, &metadata); err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, fmt.Errorf("document metadata file %q is empty or null", filePath)
	}

	metadata.fileInfo = fi
	return metadata, nil
}
