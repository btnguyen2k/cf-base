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

// TopicMeta capture metadata of a topic.
type TopicMeta struct {
	fileInfo    os.FileInfo `json:"-" yaml:"-"`                                         // internal use only!
	index       int         `json:"-" yaml:"-"`                                         // topic index, for ordering
	id          string      `json:"-" yaml:"-"`                                         // topic id
	dir         string      `json:"-" yaml:"-"`                                         // name of directory where topic's data locates
	numDocs     int         `json:"-" yaml:"-"`                                         // number of documents in this topic
	Title       interface{} `json:"title" yaml:"title"`                                 // topic's title, can be a single string, or a map[language-code:string]string
	Description interface{} `json:"description,omitempty" yaml:"description,omitempty"` // short description, can be a single string, or a map[language-code:string]string
	Icon        string      `json:"icon,omitempty" yaml:"icon,omitempty"`               // topic's icon
	EntryImage  string      `json:"img,omitempty" yaml:"img,omitempty"`                 // topic's entry image
	Hidden      bool        `json:"hidden,omitempty" yaml:"hidden,omitempty"`           // if 'true', this topic is "hidden" from GUI
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

// ToMap returns the topic metadata as a map.
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
func (tm *TopicMeta) GetDescriptionMap(defaultLang string) map[string]string {
	desc := make(map[string]string)
	if tm.Description != nil {
		switch reflect.TypeOf(tm.Description).Kind() {
		case reflect.String:
			desc[defaultLang] = fmt.Sprintf("%s", tm.Description)
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
func (tm *TopicMeta) GetTitleMap(defaultLang string) map[string]string {
	title := make(map[string]string)
	if tm.Title != nil {
		switch reflect.TypeOf(tm.Title).Kind() {
		case reflect.String:
			title[defaultLang] = fmt.Sprintf("%s", tm.Title)
		case reflect.Map:
			temp, err := reddo.Convert(tm.Title, _typMapString)
			if err == nil && temp != nil {
				title = temp.(map[string]string)
			}
		}
	}
	return title
}

// LoadTopicMetaAuto loads metadata from meta.yaml, meta.yml, or meta.json in dir.
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
	err = yaml.Unmarshal(buf, &metadata)
	if err == nil {
		metadata.fileInfo = fi
	}
	return metadata, err
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
	err = json.Unmarshal(buf, &metadata)
	if err == nil {
		metadata.fileInfo = fi
	}
	return metadata, err
}
