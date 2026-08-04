package cfbase

// Author is site's/document's author
type Author struct {
	Name   string `json:"name" yaml:"name"`
	Email  string `json:"email" yaml:"email"`
	Avatar string `json:"avatar" yaml:"avatar"`
}
