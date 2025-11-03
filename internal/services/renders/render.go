package renders

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"revonoir.com/jrender/controllers/dto/jerrors"
	"revonoir.com/jrender/internal/remotes/dto"
	"revonoir.com/jrender/internal/services/renders/dtos"
)

type FormCoreEngineIf interface {
	GenerateCSS(data dtos.FormCoreData) (string, error)
	GenerateHTML(data dtos.FormCoreData) (string, error)
}

type EmbeddedFormEngineIf interface {
	GenerateHTML(data *dto.FormResponse) (string, error)
}

type JFormClientIf interface {
	GetForm(ctx context.Context, formID uuid.UUID) (*dto.FormResponse, error)
}

type RenderService struct {
	jformClient        JFormClientIf
	embeddedFormEngine EmbeddedFormEngineIf
}

func NewRenderService(jformClient JFormClientIf, embeddedFormEngine EmbeddedFormEngineIf) *RenderService {
	return &RenderService{
		jformClient:        jformClient,
		embeddedFormEngine: embeddedFormEngine,
	}
}

func (s RenderService) RenderForm(formID string) (string, error) {

	formIDUUID, err := uuid.Parse(formID)
	if err != nil {
		slog.Error("invalid form ID", "error", err, "formID", formID)
		return "", jerrors.BadRequest("invalid form ID")
	}

	// call the jform service to get the form details
	form, err := s.jformClient.GetForm(context.Background(), formIDUUID)
	if err != nil {
		slog.Error("failed to get form from jform service", "error", err, "formID", formID)
		return "", err
	}

	html, err := s.embeddedFormEngine.GenerateHTML(form)
	if err != nil {
		slog.Error("failed to generate HTML", "error", err, "formID", form.ID)
		return "", err
	}

	// return the html, css, and javascript
	return html, nil
}
