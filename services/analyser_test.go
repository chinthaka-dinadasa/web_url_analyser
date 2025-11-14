package services

import (
	"strings"
	"testing"
	"web-analyser/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestAnalyserService_PageTitleExtraction(t *testing.T) {
	analyser := NewAnalyserService(1)
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "basic title",
			html:     `<html><head><title>Test Page</title></head><body></body></html>`,
			expected: "Test Page",
		},
		{
			name:     "basic title",
			html:     `<title>Random</title>`,
			expected: "Random",
		},
		{
			name:     "basic title",
			html:     `<html>hello</html>`,
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			result := analyser.capturePageTitle(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyserService_CaptureLoginForm(t *testing.T) {
	analyser := NewAnalyserService(1)
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "simple password input in form",
			html:     `<form><input type="password" name="password"></form>`,
			expected: true,
		},
		{
			name:     "password input with different attributes",
			html:     `<form><input type="password" placeholder="Enter password"></form>`,
			expected: true,
		},
		{
			name:     "multiple password inputs",
			html:     `<form><input type="password" name="pwd1"><input type="password" name="pwd2"></form>`,
			expected: true,
		},
		{
			name:     "password input with other inputs",
			html:     `<form><input type="text" name="username"><input type="password" name="password"></form>`,
			expected: true,
		},
		{
			name:     "form without password input",
			html:     `<form><input type="text" name="username"><input type="email" name="email"></form>`,
			expected: false,
		},
		{
			name:     "empty form",
			html:     `<form></form>`,
			expected: false,
		},
		{
			name:     "no forms at all",
			html:     `<html><body><div>Hello World</div></body></html>`,
			expected: false,
		},
		{
			name:     "password input outside form",
			html:     `<input type="password"><form><input type="text"></form>`,
			expected: false,
		},
		{
			name: "multiple forms with one having password",
			html: `
            <form id="search">
                <input type="text" placeholder="Search">
            </form>
            <form id="login">
                <input type="text" name="username">
                <input type="password" name="password">
            </form>
            <form id="contact">
                <input type="text" name="name">
            </form>`,
			expected: true,
		},
		{
			name:     "nested password input",
			html:     `<form><div><section><input type="password" name="pwd"></section></div></form>`,
			expected: true,
		},
		{
			name:     "password input with mixed case",
			html:     `<form><input type="PASSWORD" name="password"></form>`,
			expected: true,
		},
		{
			name:     "malformed password input",
			html:     `<form><input type="password" ></form>`,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			result := analyser.captureLoginForm(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyserService_CaptureLinksData(t *testing.T) {
	analyser := NewAnalyserService(1)
	tests := []struct {
		name     string
		html     string
		baseUrl  string
		expected models.WebLinkDetail
	}{
		{
			name: "Internal links",
			html: `<html>
		        <body>
		            <a href="/about">About</a>
		            <a href="/contact">Contact</a>
		            <a href="#section">Anchor</a>
		        </body>
		    </html>`,
			baseUrl: "https://www.javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     3,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "Internal and external links",
			html: `<html>
		        <body>
		            <a href="/about">About</a>
		            <a href="/contact">Contact</a>
		            <a href="https://github.com/">Anchor</a>
		        </body>
		    </html>`,
			baseUrl: "https://www.javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     2,
				ExternalLinks:     1,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "Internal and external links with 404",
			html: `<html>
                <body>
                    <a href="/about">About</a>
                    <a href="/contact">Contact</a>
                    <a href="https://noname_757a971d-ac55-4651-8622-17e62b703310393.coms">Anchor</a>
                </body>
            </html>`,
			baseUrl: "https://www.javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     2,
				ExternalLinks:     1,
				UnAccessibleLinks: 1,
			},
		},
		{
			name:    "empty page",
			html:    `<html><body>No links here</body></html>`,
			baseUrl: "https://javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     0,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "links without href",
			html: `
            <html>
                <body>
                    <a name="anchor">No href</a>
                    <a>Empty anchor</a>
                    <a href="/valid">Valid link</a>
                </body>
            </html>`,
			baseUrl: "https://javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     1,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "mailto and javascript links",
			html: `
            <html>
                <body>
                    <a href="mailto:author@javatodev.com">Email</a>
                    <a href="javascript:void(0)">JS Link</a>
                    <a href="tel:+1234567890">Phone</a>
                    <a href="/normal">Normal Link</a>
                </body>
            </html>`,
			baseUrl: "https://javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     1,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			result := analyser.captureLinksData(tt.baseUrl, doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyserService_CaptureHeadingDetails(t *testing.T) {
	analyser := NewAnalyserService(1)
	tests := []struct {
		name     string
		html     string
		expected models.HeadingDetail
	}{
		{
			name: "No Headings",
			html: `<html>
		        <body>
		            <a href="/about">About</a>
		            <a href="/contact">Contact</a>
		            <a href="#section">Anchor</a>
		        </body>
		    </html>`,
			expected: models.HeadingDetail{
				H1: 0,
				H2: 0,
				H3: 0,
				H4: 0,
				H5: 0,
				H6: 0,
			},
		},
		{
			name: "Single H1 Heading",
			html: `<html>
                <body>
                    <h1>Main Title</h1>
                    <p>Some content</p>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 1,
				H2: 0,
				H3: 0,
				H4: 0,
				H5: 0,
				H6: 0,
			},
		},
		{
			name: "Multiple H1 Headings",
			html: `<html>
                <body>
                    <h1>First H1</h1>
                    <h1>Second H1</h1>
                    <h1>Third H1</h1>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 3,
				H2: 0,
				H3: 0,
				H4: 0,
				H5: 0,
				H6: 0,
			},
		},
		{
			name: "All Heading Levels",
			html: `<html>
                <body>
                    <h1>Heading 1</h1>
                    <h2>Heading 2</h2>
                    <h3>Heading 3</h3>
                    <h4>Heading 4</h4>
                    <h5>Heading 5</h5>
                    <h6>Heading 6</h6>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 1,
				H2: 1,
				H3: 1,
				H4: 1,
				H5: 1,
				H6: 1,
			},
		},
		{
			name: "Mixed Multiple Headings",
			html: `<html>
                <body>
                    <h1>Main Title</h1>
                    <h2>Section 1</h2>
                    <h3>Subsection 1.1</h3>
                    <h3>Subsection 1.2</h3>
                    <h2>Section 2</h2>
                    <h3>Subsection 2.1</h3>
                    <h4>Detail 2.1.1</h4>
                    <h4>Detail 2.1.2</h4>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 1,
				H2: 2,
				H3: 3,
				H4: 2,
				H5: 0,
				H6: 0,
			},
		},
		{
			name: "Headings with Attributes and Classes",
			html: `<html>
                <body>
                    <h1 class="main-title" id="title1">Main Title</h1>
                    <h2 data-test="section" style="color: red;">Section</h2>
                    <h3 hidden>Hidden Heading</h3>
                    <h4 data-count="1">Heading with data</h4>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 1,
				H2: 1,
				H3: 1,
				H4: 1,
				H5: 0,
				H6: 0,
			},
		},
		{
			name: "Nested Headings in Divs and Sections",
			html: `<html>
                <body>
                    <div class="header">
                        <h1>Site Title</h1>
                    </div>
                    <main>
                        <section>
                            <h2>Article Title</h2>
                            <div class="content">
                                <h3>Content Heading</h3>
                                <article>
                                    <h4>Sub Heading</h4>
                                </article>
                            </div>
                        </section>
                    </main>
                    <footer>
                        <h5>Footer Heading</h5>
                    </footer>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 1,
				H2: 1,
				H3: 1,
				H4: 1,
				H5: 1,
				H6: 0,
			},
		},
		{
			name: "Empty Heading Tags",
			html: `<html>
                <body>
                    <h1></h1>
                    <h2> </h2>
                    <h3><!-- Comment --></h3>
                    <h4><span></span></h4>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 1,
				H2: 1,
				H3: 1,
				H4: 1,
				H5: 0,
				H6: 0,
			},
		},
		{
			name: "Malformed Headings",
			html: `<html>
                <body>
                    <h1>Proper H1</h1>
                    <H2>Uppercase H2</H2>
                    <h3>Mixed <span>content</span> H3</h3>
                    <h7>Invalid H7</h7>
                </body>
            </html>`,
			expected: models.HeadingDetail{
				H1: 1,
				H2: 1,
				H3: 1,
				H4: 0,
				H5: 0,
				H6: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			result := analyser.captureHeadingDetails(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyserService_CaptureHTMLVersion(t *testing.T) {
	analyser := NewAnalyserService(1)
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5 doctype",
			html:     `<!DOCTYPE html><html><head></head><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "HTML5 doctype with attributes",
			html:     `<!DOCTYPE html SYSTEM "about:legacy-compat"><html><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "HTML5 lowercase doctype",
			html:     `<!doctype html><html><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "HTML 4.01 Strict",
			html:     `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd"><html><body></body></html>`,
			expected: "HTML4",
		},
		{
			name:     "HTML 4.01 Transitional",
			html:     `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd"><html><body></body></html>`,
			expected: "HTML4",
		},
		{
			name:     "XHTML 1.0 Strict",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"><html xmlns="http://www.w3.org/1999/xhtml"><body></body></html>`,
			expected: "XHTML",
		},
		{
			name:     "XHTML 1.1",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd"><html xmlns="http://www.w3.org/1999/xhtml"><body></body></html>`,
			expected: "XHTML",
		},
		{
			name:     "Basic HTML without doctype",
			html:     `<html><head></head><body>Hello World</body></html>`,
			expected: "HTML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "Failed to parse HTML for test: %s", tt.name)

			result := analyser.captureHTMLVersion(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}
