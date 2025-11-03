package forms

import (
	"fmt"
	"html/template"

	"github.com/go-chi/chi/v5"
	"revonoir.com/jrender/conns/configs"
	"revonoir.com/jrender/internal/remotes"
	"revonoir.com/jrender/internal/services/renders"
	"revonoir.com/jrender/internal/services/templates"
)

func Route(r chi.Router, config configs.Configuration) {
	jformClient := remotes.NewJformClient(config)

	renderer := templates.NewFormRenderer(config.Captcha.Provider.ReCaptcha.SiteKey)

	// Create CSS template with shared functions
	cssTemplate := template.New("form_core_css").Funcs(renderer.GetFuncMap())
	cssTemplate, err := cssTemplate.Parse(templates.FormCoreTemplate)
	if err != nil {
		panic("Failed to parse form core CSS template: " + err.Error())
	}

	// Create HTML template with shared functions
	htmlTemplate := template.New("form_core_html").Funcs(renderer.GetFuncMap())
	htmlTemplate, err = htmlTemplate.Parse(templates.FormCoreHTMLTemplate)
	if err != nil {
		panic("Failed to parse form core HTML template: " + err.Error())
	}
	formCoreEngine := templates.NewFormCoreEngine(cssTemplate, htmlTemplate)

	// Create wrapper template
	wrapperTemplate := template.New("form_production_wrapper")
	wrapperTemplate, err = wrapperTemplate.Parse(templates.CompleteHTMLTemplate)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse form production wrapper template: %v", err))
	}

	embeddedFormEngine := templates.NewEmbeddedFormEngine(formCoreEngine, wrapperTemplate)
	renderService := renders.NewRenderService(jformClient, embeddedFormEngine)
	formController := NewFormController(renderService)
	r.Route("/f", func(r chi.Router) {
		r.Get("/{formID}", formController.DisplayDirectForm)
	})
}
