package cfbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/btnguyen2k/consu/reddo"
	"gopkg.in/yaml.v3"
)

// SiteMeta describes a website.
type SiteMeta struct {
	fileInfo os.FileInfo `json:"-" yaml:"-"`

	// Name is the website name.
	Name string `json:"name" yaml:"name"`
	// Description is either a string or a map of language codes to descriptions.
	Description any `json:"description" yaml:"description"`
	// Languages maps available language codes to their display names.
	Languages map[string]string `json:"languages" yaml:"languages"`
	// DefaultLanguage is the resolved default language code.
	DefaultLanguage string `json:"-" yaml:"-"`
	// Icon is the location of the website icon.
	Icon string `json:"icon" yaml:"icon"`
	// Contacts contains the website's contact details.
	Contacts map[string]string `json:"contacts,omitempty" yaml:"contacts,omitempty"`
	// Tags contains arbitrary website metadata.
	Tags map[string]any `json:"tags,omitempty" yaml:"tags,omitempty"`
	// TagsAlias contains aliases by tag, optionally grouped by language code.
	TagsAlias any `json:"tagalias,omitempty" yaml:"tagalias,omitempty"`
	// Mode is the website mode, such as document or blog.
	Mode string `json:"mode" yaml:"mode"`
	// Author is the website's author and the default author for its documents.
	Author *Author `json:"author,omitempty" yaml:"author,omitempty"`
}

func (sm *SiteMeta) init() error {
	// init field "default language"
	defaultLang := sm.Languages["default"]
	languageCount := len(sm.Languages)
	if _, exists := sm.Languages["default"]; exists {
		languageCount--
	}
	if languageCount == 0 {
		return fmt.Errorf("at least one language must be configured")
	}
	if defaultLang == "" {
		if languageCount != 1 {
			return fmt.Errorf("languages.default is required when multiple languages are configured")
		}
		for lang := range sm.Languages {
			if lang != "default" {
				defaultLang = lang
				break
			}
		}
	} else if _, exists := sm.Languages[defaultLang]; !exists {
		return fmt.Errorf("default language %q is not configured", defaultLang)
	}
	sm.DefaultLanguage = defaultLang

	// verify "site-mode"
	ok := false
	for _, mode := range SiteModes() {
		if mode == sm.Mode {
			ok = true
			break
		}
	}
	if !ok {
		sm.Mode = DefaultSiteMode
	}

	// normalize field "contacts"
	for k, v := range sm.Contacts {
		if v == "" {
			delete(sm.Contacts, k)
		}
	}

	// normalize field "tags"
	if sm.Tags == nil {
		sm.Tags = make(map[string]any)
	}

	return nil
}

// ToMap returns a map representation of the site metadata.
func (sm *SiteMeta) ToMap() map[string]any {
	return map[string]any{
		"name":            sm.Name,
		"languages":       sm.Languages,
		"defaultLanguage": sm.DefaultLanguage,
		"icon":            sm.Icon,
		"description":     sm.GetDescriptionMap(),
		"contacts":        sm.Contacts,
		"tags":            sm.Tags,
		"mode":            sm.Mode,
		"author":          sm.Author,
	}
}

// GetDescriptionMap returns the site description keyed by language code.
func (sm *SiteMeta) GetDescriptionMap() map[string]string {
	desc := make(map[string]string)
	if sm.Description != nil {
		switch reflect.TypeOf(sm.Description).Kind() {
		case reflect.String:
			desc[sm.DefaultLanguage] = fmt.Sprintf("%s", sm.Description)
		case reflect.Map:
			temp, err := reddo.Convert(sm.Description, _typMapString)
			if err == nil && temp != nil {
				desc = temp.(map[string]string)
			}
		}
	}
	return desc
}

// GetTagAliasMap returns tag aliases grouped by language code and tag.
func (sm *SiteMeta) GetTagAliasMap() map[string]map[string][]string {
	empty := make(map[string]map[string][]string)
	outer, err := reddo.Convert(sm.TagsAlias, _typMapAny)
	if err != nil || outer == nil {
		// the top level must be a map
		return empty
	}

	var nextLevelKind reflect.Kind
	for _, v := range outer.(map[string]any) {
		if v != nil {
			nextLevelKind = reflect.TypeOf(v).Kind()
			break
		}
	}
	switch nextLevelKind {
	case reflect.Array, reflect.Slice:
		// then the next level must be either array/slice
		if result, err := reddo.Convert(sm.TagsAlias, reflect.TypeOf(make(map[string][]string))); err == nil && result != nil {
			return map[string]map[string][]string{
				sm.DefaultLanguage: result.(map[string][]string),
			}
		}
	case reflect.Map:
		// or a map
		if result, err := reddo.Convert(sm.TagsAlias, reflect.TypeOf(empty)); err == nil && result != nil {
			return result.(map[string]map[string][]string)
		}
	}

	return empty
}

// LoadSiteMetaAuto loads the first available meta.yaml, meta.yml, or meta.json file in dir.
func LoadSiteMetaAuto(dir string) (*SiteMeta, error) {
	yamlFiles := []string{dir + "/meta.yaml", dir + "/meta.yml"}
	for _, yamlFilePath := range yamlFiles {
		siteMeta, err := LoadSiteMetaFromYaml(yamlFilePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return siteMeta, nil
	}

	jsonFiles := []string{dir + "/meta.json"}
	for _, jsonFilePath := range jsonFiles {
		siteMeta, err := LoadSiteMetaFromJson(jsonFilePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return siteMeta, nil
	}

	return nil, fmt.Errorf("no meta file found")
}

// LoadSiteMetaFromYaml loads site metadata from a YAML file.
func LoadSiteMetaFromYaml(filePath string) (*SiteMeta, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var metadata *SiteMeta
	if err := yaml.Unmarshal(buf, &metadata); err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, fmt.Errorf("site metadata file %q is empty or null", filePath)
	}

	metadata.fileInfo = fi
	if err := metadata.init(); err != nil {
		return nil, err
	}
	return metadata, nil
}

// LoadSiteMetaFromJson loads site metadata from a JSON file.
func LoadSiteMetaFromJson(filePath string) (*SiteMeta, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var metadata *SiteMeta
	if err := json.Unmarshal(buf, &metadata); err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, fmt.Errorf("site metadata file %q is empty or null", filePath)
	}

	metadata.fileInfo = fi
	if err := metadata.init(); err != nil {
		return nil, err
	}
	return metadata, nil
}
