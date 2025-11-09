package embeds

import (
	"html/template"
	texttemplate "text/template"

	"github.com/go-chi/chi/v5"
	"revonoir.com/jrender/conns/configs"
	"revonoir.com/jrender/internal/remotes"
	"revonoir.com/jrender/internal/services/embeds"
	"revonoir.com/jrender/internal/services/templates"
)

func Route(r chi.Router, config configs.Configuration) {

	// Create FormRenderer to get template function map
	renderer := templates.NewFormRenderer(config.Base.URL, config.Captcha.Provider.ReCaptcha.SiteKey)

	// Parse CSS template with shared functions
	cssTemplate := template.New("form_core_css").Funcs(renderer.GetFuncMap())
	cssTemplate, err := cssTemplate.Parse(templates.FormCoreTemplate)
	if err != nil {
		panic("Failed to parse form core CSS template: " + err.Error())
	}

	// Parse HTML template with shared functions
	htmlTemplate := template.New("form_core_html").Funcs(renderer.GetFuncMap())
	htmlTemplate, err = htmlTemplate.Parse(templates.FormCoreHTMLTemplate)
	if err != nil {
		panic("Failed to parse form core HTML template: " + err.Error())
	}

	// Create Javascript template with shared functions
	jsTemplate := template.New("form_core_js").Funcs(renderer.GetFuncMap())
	jsTemplate, err = jsTemplate.Parse(templates.FormCoreJsTemplate)
	if err != nil {
		panic("Failed to parse form core Javascript template: " + err.Error())
	}

	sharedJsTemplate := texttemplate.New("form_core_shared_js").Funcs(renderer.GetFuncMap())
	sharedJsTemplate, err = sharedJsTemplate.Parse(templates.FormSharedJavascriptTemplate)
	if err != nil {
		panic("Failed to parse form core shared Javascript template: " + err.Error())
	}

	// Create FormCoreEngine for server-side rendering
	formCoreEngine := templates.NewFormCoreEngine(cssTemplate, htmlTemplate, jsTemplate, sharedJsTemplate)

	// Parse the embed script template (text/template, not html/template)
	embedScriptTmpl, err := texttemplate.New("embed_script").Parse(templates.EmbedScriptTemplate)
	if err != nil {
		panic("Failed to parse embed script template: " + err.Error())
	}

	embedScriptGenerator := templates.NewEmbedScriptGenerator(config.Base.URL, embedScriptTmpl, formCoreEngine)

	jformClient := remotes.NewJformClient(config)
	embedService := embeds.NewEmbedService(jformClient, embedScriptGenerator, formCoreEngine)

	// Create EmbedController with both dependencies
	embedController := NewEmbedController(embedService)

	// Universal embed script (cached, shared across all forms)
	r.Get("/embedv1.js", embedController.HandleEmbedScript)

	// Public API routes for embed forms
	r.Route("/api/public/v1/embeds", func(r chi.Router) {
		// Get form data JSON (for dynamic loading)
		r.Get("/{embedId}/data", embedController.HandleGetFormData)

		// Form submission endpoint
		r.Post("/{embedId}/submissions", embedController.HandleFormSubmission)
	})
}
