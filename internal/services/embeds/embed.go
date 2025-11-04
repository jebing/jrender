package embeds

import (
	"context"
	"log/slog"
	"slices"

	"github.com/google/uuid"
	"revonoir.com/jrender/controllers/dto/jerrors"
	"revonoir.com/jrender/internal/remotes/dto"
	"revonoir.com/jrender/internal/services/embeds/dtos"
	rendersdtos "revonoir.com/jrender/internal/services/renders/dtos"
)

type EmbedScriptGeneratorIf interface {
	GenerateEmbedScript() string
}

type JFormClientIf interface {
	GetForm(ctx context.Context, formID uuid.UUID) (*dto.FormResponse, error)
}

type FormCoreEngineIf interface {
	GenerateCSSDynamic(data rendersdtos.FormCoreData) (string, error)
	GenerateHTML(data rendersdtos.FormCoreData) (string, error)
}

type EmbedService struct {
	jformClient          JFormClientIf
	embedScriptGenerator EmbedScriptGeneratorIf
	formCoreEngine       FormCoreEngineIf
}

func NewEmbedService(jformClient JFormClientIf, embedScriptGenerator EmbedScriptGeneratorIf, formCoreEngine FormCoreEngineIf) *EmbedService {
	return &EmbedService{
		jformClient:          jformClient,
		embedScriptGenerator: embedScriptGenerator,
		formCoreEngine:       formCoreEngine,
	}
}

func (s EmbedService) GenerateEmbedScript(formID string) (string, error) {
	script := s.embedScriptGenerator.GenerateEmbedScript()
	return script, nil
}

func (s EmbedService) GenerateDynamicHTML(formID uuid.UUID, lang string) (dtos.DynamicHTMLData, error) {

	form, err := s.jformClient.GetForm(context.Background(), formID)
	if err != nil {
		slog.Error("failed to get form from jform service", "error", err, "formID", formID)
		return dtos.DynamicHTMLData{}, jerrors.InternalServerError("failed to get form from jform service")
	}

	// figure out the valid language
	bestLang := s.selectBestLanguage(lang, form.FormDefinition.Languages)

	coreData := rendersdtos.FormCoreData{
		Language: bestLang,
		FormData: rendersdtos.FormData{
			FormID:         formID.String(),
			FormDefinition: form.FormDefinition,
			FormStyling:    form.FormStyling,
		},
	}

	// Generate the dynamic CSS and HTML
	css, err := s.formCoreEngine.GenerateCSSDynamic(coreData)
	if err != nil {
		slog.Error("failed to generate dynamic CSS", "error", err, "formID", formID)
		return dtos.DynamicHTMLData{}, jerrors.InternalServerError("failed to generate dynamic CSS")
	}

	html, err := s.formCoreEngine.GenerateHTML(coreData)
	if err != nil {
		slog.Error("failed to generate dynamic HTML", "error", err, "formID", formID)
		return dtos.DynamicHTMLData{}, jerrors.InternalServerError("failed to generate dynamic HTML")
	}

	return dtos.DynamicHTMLData{
		Lang:   bestLang,
		FormID: formID.String(),
		Css:    css,
		Html:   html,
	}, nil

}

func (s EmbedService) selectBestLanguage(lang string, languages rendersdtos.FormLanguageSettings) string {
	// if the lang is empty, return the default language
	if lang == "" {
		return languages.Default
	}

	// if the lang is in the supported languages, return the lang
	if slices.Contains(languages.Supported, lang) {
		return lang
	}

	// otherwise, return the default language
	return languages.Default
}
