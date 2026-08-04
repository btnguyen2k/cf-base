package cfbase

const (
	DefaultSiteMode  = SiteModeDocument
	SiteModeDocument = "document"
	SiteModeBlog     = "blog"
)

var (
	AllSiteModes = []string{SiteModeDocument, SiteModeBlog}
)
