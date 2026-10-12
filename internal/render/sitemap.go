package render

import (
	"encoding/xml"
	"strings"
)

// Sitemap returns sitemap.xml content listing paths under siteURL.
func Sitemap(siteURL string, paths []string) []byte {
	base := strings.TrimSuffix(siteURL, "/")
	var b strings.Builder
	b.WriteString(xml.Header)
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, p := range paths {
		b.WriteString("  <url><loc>")
		_ = xml.EscapeText(&b, []byte(base+p))
		b.WriteString("</loc></url>\n")
	}
	b.WriteString("</urlset>\n")
	return []byte(b.String())
}

func Robots(siteURL string) []byte {
	return []byte("User-agent: *\nAllow: /\nSitemap: " + strings.TrimSuffix(siteURL, "/") + "/sitemap.xml\n")
}
