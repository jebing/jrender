package embeds

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"revonoir.com/jrender/controllers/dto"
	"revonoir.com/jrender/controllers/dto/jerrors"
	"revonoir.com/jrender/internal/services/embeds/dtos"
)

type EmbedServiceIf interface {
	GenerateEmbedScript(formID string) (string, error)
	GenerateDynamicHTML(formID uuid.UUID, lang string) (dtos.DynamicHTMLData, error)
}

type EmbedController struct {
	embedService EmbedServiceIf
}

func NewEmbedController(embedService EmbedServiceIf) *EmbedController {
	return &EmbedController{
		embedService: embedService,
	}
}

// HandleEmbedScript serves the universal embed.js file
func (c EmbedController) HandleEmbedScript(w http.ResponseWriter, r *http.Request) {
	// Generate embed.js content
	embedID := chi.URLParam(r, "embedId")
	script, err := c.embedService.GenerateEmbedScript(embedID)
	if err != nil {
		jerrors.WriteErrorResponse(w, err)
		return
	}

	// Set headers for long-term caching
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	// Write script
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(script))
}

// HandleGetFormData serves server-rendered form HTML with dynamic CSS
func (c EmbedController) HandleGetFormData(w http.ResponseWriter, r *http.Request) {
	embedID := chi.URLParam(r, "embedId")
	requestedLang := r.URL.Query().Get("lang")

	formID, err := uuid.Parse(embedID)
	if err != nil {
		jerrors.WriteErrorResponse(w, jerrors.BadRequest("invalid embed ID"))
		return
	}

	dynamicHTMLData, err := c.embedService.GenerateDynamicHTML(formID, requestedLang)
	if err != nil {
		jerrors.WriteErrorResponse(w, err)
		return
	}

	response := dto.Response[dtos.DynamicHTMLData]{
		Data: dynamicHTMLData,
	}

	// Set headers
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=60")
	w.Header().Set("X-Form-Language", dynamicHTMLData.Lang) // Tell client which language was actually used

	// Write HTML response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleFormSubmission handles form submission POST requests
func (c *EmbedController) HandleFormSubmission(w http.ResponseWriter, r *http.Request) {
	embedID := chi.URLParam(r, "embedId")

	// TODO: Implement actual form submission logic
	// For now, return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Form submitted successfully!",
		"embedId": embedID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
