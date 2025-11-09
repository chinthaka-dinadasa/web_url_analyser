package models

type WebAnalysingRequest struct {
	Url string `json:"url" binding:"required"`
}

type WebAnalysingResponse struct {
	HTMLVersion       string        `json:"htmlVersion"`
	PageTitle         string        `json:"pageTitle"`
	Heading           HeadingDetail `json:"headings"`
	InternalLinks     int16         `json:"internalLinks"`
	ExternalLinks     int16         `json:"externalLinks"`
	UnaccessibleLinks int16         `json:"unaccessibleLinks"`
	Error             string        `json:"error"`
}

type HeadingDetail struct {
	H1 int16 `json:"h1"`
	H2 int16 `json:"h2"`
	H3 int16 `json:"h3"`
	H4 int16 `json:"h4"`
	H5 int16 `json:"h5"`
	H6 int16 `json:"h6"`
}
