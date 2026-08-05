package cfbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/btnguyen2k/consu/reddo"
	"gopkg.in/yaml.v3"
)

// TopicMeta describes a topic and its presentation metadata.
type TopicMeta struct {
	fileInfo os.FileInfo `json:"-" yaml:"-"`
	index    int         `json:"-" yaml:"-"`
	id       string      `json:"-" yaml:"-"`
	dir      string      `json:"-" yaml:"-"`
	numDocs  int         `json:"-" yaml:"-"`

	// DefaultLanguage is the default language inherited from the containing site.
	DefaultLanguage string `json:"def_lang" yaml:"def_lang"`
	// Title is either a string or a map of language codes to titles.
	Title interface{} `json:"title" yaml:"title"`
	// Description is either a string or a map of language codes to descriptions.
	Description interface{} `json:"description,omitempty" yaml:"description,omitempty"`
	// Icon is the location of the topic icon.
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty"`
	// EntryImage is the location of the topic's entry image.
	EntryImage string `json:"img,omitempty" yaml:"img,omitempty"`
	// Hidden indicates whether the topic is hidden from the user interface.
	Hidden bool `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

func (tm *TopicMeta) setDirectory(dir string) bool {
	tm.dir = dir
	if !_rexpContentDir.MatchString(dir) {
		return false
	}
	matches := _rexpContentDir.FindStringSubmatch(dir)
	tm.index, _ = strconv.Atoi(matches[1])
	tm.id = matches[2]
	return true
}

// ToMap returns a map representation of the topic metadata.
func (tm *TopicMeta) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          tm.id,
		"num_docs":    tm.numDocs,
		"icon":        tm.Icon,
		"title":       tm.GetTitleMap(),
		"description": tm.GetDescriptionMap(),
		"img":         tm.EntryImage,
		"hidden":      tm.Hidden,
	}
}

// GetDescriptionMap returns the topic description keyed by language code.
func (tm *TopicMeta) GetDescriptionMap() map[string]string {
	desc := make(map[string]string)
	if tm.Description != nil {
		switch reflect.TypeOf(tm.Description).Kind() {
		case reflect.String:
			desc[tm.DefaultLanguage] = fmt.Sprintf("%s", tm.Description)
		case reflect.Map:
			temp, err := reddo.Convert(tm.Description, _typMapString)
			if err == nil && temp != nil {
				desc = temp.(map[string]string)
			}
		}
	}
	return desc
}

// GetTitleMap returns the topic title keyed by language code.
func (tm *TopicMeta) GetTitleMap() map[string]string {
	title := make(map[string]string)
	if tm.Title != nil {
		switch reflect.TypeOf(tm.Title).Kind() {
		case reflect.String:
			title[tm.DefaultLanguage] = fmt.Sprintf("%s", tm.Title)
		case reflect.Map:
			temp, err := reddo.Convert(tm.Title, _typMapString)
			if err == nil && temp != nil {
				title = temp.(map[string]string)
			}
		}
	}
	return title
}

// LoadTopicMetaAuto loads the first available meta.yaml, meta.yml, or meta.json file in dir.
func LoadTopicMetaAuto(dir string) (*TopicMeta, error) {
	yamlFiles := []string{dir + "/meta.yaml", dir + "/meta.yml"}
	for _, yamlFilePath := range yamlFiles {
		topicMeta, err := LoadTopicMetaFromYaml(yamlFilePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return topicMeta, nil
	}

	jsonFiles := []string{dir + "/meta.json"}
	for _, jsonFilePath := range jsonFiles {
		topicMeta, err := LoadTopicMetaFromJson(jsonFilePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return topicMeta, nil
	}

	return nil, fmt.Errorf("no meta file found")
}

// LoadTopicMetaFromYaml loads topic metadata from a YAML file.
func LoadTopicMetaFromYaml(filePath string) (*TopicMeta, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var metadata *TopicMeta
	if err := yaml.Unmarshal(buf, &metadata); err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, fmt.Errorf("topic metadata file %q is empty or null", filePath)
	}

	metadata.fileInfo = fi
	return metadata, nil
}

// LoadTopicMetaFromJson loads topic metadata from a JSON file.
func LoadTopicMetaFromJson(filePath string) (*TopicMeta, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var metadata *TopicMeta
	if err := json.Unmarshal(buf, &metadata); err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, fmt.Errorf("topic metadata file %q is empty or null", filePath)
	}

	metadata.fileInfo = fi
	return metadata, nil
}
