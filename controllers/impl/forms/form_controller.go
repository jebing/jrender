package forms

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"revonoir.com/jrender/controllers/dto/jerrors"
)

type RenderServiceIf interface {
	RenderForm(formID string) (string, error)
}

type FormController struct {
	renderService RenderServiceIf
}

func NewFormController(renderService RenderServiceIf) *FormController {
	return &FormController{
		renderService: renderService,
	}
}

func (c *FormController) DisplayDirectForm(w http.ResponseWriter, r *http.Request) {
	formID := chi.URLParam(r, "formID")
	html, err := c.renderService.RenderForm(formID)
	if err != nil {
		jerrors.WriteErrorResponse(w, err)
		return
	}
	w.Write([]byte(html))
}
