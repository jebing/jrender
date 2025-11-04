package dtos

type DynamicHTMLData struct {
	Lang   string `json:"lang,omitempty"`
	FormID string `json:"form_id,omitempty"`
	Css    string `json:"css,omitempty"`
	Html   string `json:"html,omitempty"`
}
