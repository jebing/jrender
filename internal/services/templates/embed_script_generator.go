package templates

import (
	"bytes"
	_ "embed"
	"log/slog"
	"text/template"

	"revonoir.com/jrender/internal/services/renders/dtos"
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

type SharedJavascriptCoreEngineIf interface {
	GenerateSharedJavascript(data dtos.FormCoreData) (string, error)
}

type EmbedScriptGenerator struct {
	apiBaseURL                 string
	parsedTemplate             *template.Template
	cssContent                 string
	sharedJavascriptCoreEngine SharedJavascriptCoreEngineIf
}

func NewEmbedScriptGenerator(apiBaseURL string, parsedTemplate *template.Template, sharedJavascriptCoreEngine SharedJavascriptCoreEngineIf) *EmbedScriptGenerator {

	// Extract static content at initialization
	cssContent := extractCSSContent()

	return &EmbedScriptGenerator{
		apiBaseURL:                 apiBaseURL,
		parsedTemplate:             parsedTemplate,
		cssContent:                 cssContent,
		sharedJavascriptCoreEngine: sharedJavascriptCoreEngine,
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
		slog.Error("failed to generate embed script", "error", err)
		// For now, return empty script on error
		return ""
	}

	// formData := dtos.FormCoreData{}
	// sharedJavascript, err := g.sharedJavascriptCoreEngine.GenerateSharedJavascript(formData)
	// if err != nil {
	// 	slog.Error("failed to generate shared Javascript", "error", err)
	// 	return ""
	// }

	return SharedScriptTemplate + "\n" + buf.String()
}

// extractCSSContent returns only the static CSS (no template directives)
func extractCSSContent() string {
	// FormCoreCSSStatic contains only static CSS without Go template directives
	// Safe to inject directly into JavaScript without template processing
	return FormCoreCSSStatic
}
