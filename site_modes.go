package cfbase

const (
	// DefaultSiteMode is the mode used when no valid mode is configured.
	DefaultSiteMode = SiteModeDocument

	// SiteModeDocument identifies a documentation site.
	SiteModeDocument = "document"

	// SiteModeBlog identifies a blog site.
	SiteModeBlog = "blog"
)

// IsValidSiteMode reports whether mode is supported.
func IsValidSiteMode(mode string) bool {
	switch mode {
	case SiteModeDocument, SiteModeBlog:
		return true
	default:
		return false
	}
}

// SiteModes returns all supported site modes.
func SiteModes() []string {
	return []string{SiteModeDocument, SiteModeBlog}
}
