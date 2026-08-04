package cfbase

const (
	DefaultSiteMode  = SiteModeDocument
	SiteModeDocument = "document"
	SiteModeBlog     = "blog"
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
