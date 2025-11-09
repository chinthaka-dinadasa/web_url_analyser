package models

type WebAnalysingRequest struct {
	Url string `json:"url" binding:"required"`
}

type WebAnalysingResponse struct {
	HTMLVersion       string        `json:"htmlVersion"`
	PageTitle         string        `json:"pageTitle"`
	Heading           HeadingDetail `json:"headings"`
	InternalLinks     int           `json:"internalLinks"`
	ExternalLinks     int           `json:"externalLinks"`
	UnaccessibleLinks int           `json:"unaccessibleLinks"`
	Error             string        `json:"error"`
}

type HeadingDetail struct {
	H1 int `json:"h1"`
	H2 int `json:"h2"`
	H3 int `json:"h3"`
	H4 int `json:"h4"`
	H5 int `json:"h5"`
	H6 int `json:"h6"`
}
