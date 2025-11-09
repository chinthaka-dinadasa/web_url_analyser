package models

type WebAnalysingRequest struct {
	Url string `json:"url" binding:"required"`
}

type WebAnalysingResponse struct {
	HTMLVersion string `json:"htmlVersion"`
	PageTitle   string `json:"pageTitle"`
}
