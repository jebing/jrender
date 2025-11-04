package templates

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed embed_script_template.js
var EmbedScriptTemplate string

// EmbedScriptTemplateData holds data for populating the embed script template
type EmbedScriptTemplateData struct {
	APIBaseURL          string
	CSSContent          string
	ValidationFunctions string
	SubmitFunctions     string
}

type EmbedScriptGenerator struct {
	apiBaseURL     string
	parsedTemplate *template.Template
	cssContent     string
}

func NewEmbedScriptGenerator(apiBaseURL string, parsedTemplate *template.Template) *EmbedScriptGenerator {

	// Extract static content at initialization
	cssContent := extractCSSContent()

	return &EmbedScriptGenerator{
		apiBaseURL:     apiBaseURL,
		parsedTemplate: parsedTemplate,
		cssContent:     cssContent,
	}
}

// GenerateEmbedScript generates the complete embed.js content
func (g EmbedScriptGenerator) GenerateEmbedScript() string {
	// Prepare template data
	data := EmbedScriptTemplateData{
		APIBaseURL: g.apiBaseURL,
		CSSContent: g.cssContent,
	}

	// Execute template
	var buf bytes.Buffer
	if err := g.parsedTemplate.Execute(&buf, data); err != nil {
		// In production, we should log this error properly
		// For now, return empty script on error
		return ""
	}

	return buf.String()
}

// extractCSSContent returns only the static CSS (no template directives)
func extractCSSContent() string {
	// FormCoreCSSStatic contains only static CSS without Go template directives
	// Safe to inject directly into JavaScript without template processing
	return FormCoreCSSStatic
}
