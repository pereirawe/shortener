package api

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// SEOData holds the metadata extracted from a URL
type SEOData struct {
	Available   bool
	Title       string
	Description string
	Image       string
}

var seoHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return http.ErrUseLastResponse
		}
		return nil
	},
}

// fetchSEO requests the given URL and parses SEO metadata from its HTML.
// If the URL is unreachable or returns a non-2xx status, Available is false.
func fetchSEO(rawURL string) SEOData {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return SEOData{Available: false}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Shortener-Bot/1.0)")

	resp, err := seoHTTPClient.Do(req)
	if err != nil {
		return SEOData{Available: false}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return SEOData{Available: false}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		// Non-HTML resource (image, PDF, etc.) — URL is reachable but no SEO to parse
		return SEOData{Available: true}
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return SEOData{Available: true}
	}

	data := SEOData{Available: true}
	parseNode(doc, &data)
	return data
}

// parseNode walks the HTML node tree to extract SEO meta tags and the title.
func parseNode(n *html.Node, data *SEOData) {
	if n.Type == html.ElementNode {
		switch strings.ToLower(n.Data) {
		case "title":
			if n.FirstChild != nil && data.Title == "" {
				data.Title = strings.TrimSpace(n.FirstChild.Data)
			}
		case "meta":
			name := attrVal(n, "name")
			prop := attrVal(n, "property")
			content := attrVal(n, "content")

			switch {
			case strings.EqualFold(name, "description") && data.Description == "":
				data.Description = content
			case strings.EqualFold(prop, "og:title") && data.Title == "":
				data.Title = content
			case strings.EqualFold(prop, "og:description") && data.Description == "":
				data.Description = content
			case strings.EqualFold(prop, "og:image") && data.Image == "":
				data.Image = content
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c, data)
	}
}

// attrVal returns the value of the named attribute from an HTML node.
func attrVal(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, name) {
			return a.Val
		}
	}
	return ""
}
