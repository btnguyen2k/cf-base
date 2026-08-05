package cfbase

// Author is site's/document's author
type Author struct {
	// Name is the author's display name.
	Name string `json:"name" yaml:"name"`
	// Email is the author's email address.
	Email string `json:"email" yaml:"email"`
	// Avatar is the location of the author's avatar image.
	Avatar string `json:"avatar" yaml:"avatar"`
}
