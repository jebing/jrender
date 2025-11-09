package embeds

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"revonoir.com/jrender/controllers/dto"
	"revonoir.com/jrender/controllers/dto/jerrors"
	"revonoir.com/jrender/internal/services/embeds/dtos"
)

type EmbedServiceIf interface {
	GenerateEmbedScript(formID string) (string, error)
	GenerateDynamicHTML(formID uuid.UUID, lang string) (dtos.DynamicHTMLData, error)
	SubmitForm(formID uuid.UUID, data map[string]interface{}, headers http.Header) error
}

type EmbedController struct {
	embedService EmbedServiceIf
}

func NewEmbedController(embedService EmbedServiceIf) *EmbedController {
	return &EmbedController{
		embedService: embedService,
	}
}

func (c EmbedController) HandleSubmitForm(w http.ResponseWriter, r *http.Request) {
	embedID := chi.URLParam(r, "embedId")
	formID, err := uuid.Parse(embedID)
	if err != nil {
		jerrors.WriteErrorResponse(w, jerrors.BadRequest("invalid embed ID"))
		return
	}

	data := map[string]interface{}{}
	// limit the body to avoid ddos
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*10) // 10MB
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		slog.Error("failed to decode form data", "error", err)
		jerrors.WriteErrorResponse(w, jerrors.BadRequest("invalid form data"))
		return
	}
	headers := r.Header

	err = c.embedService.SubmitForm(formID, data, headers)
	if err != nil {
		jerrors.WriteErrorResponse(w, err)
		return
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

	formID, err := uuid.Parse(embedID)
	if err != nil {
		jerrors.WriteErrorResponse(w, jerrors.BadRequest("invalid embed ID"))
		return
	}

	// Parse submission data
	if err := r.ParseForm(); err != nil {
		slog.Warn("invalid submission data", "error", err)
		jerrors.WriteErrorResponse(w, jerrors.BadRequest("invalid submission data"))
		return
	}

	submissionData := make(map[string]interface{})

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		// Parse JSON data
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*10) // 10MB
		if err := json.NewDecoder(r.Body).Decode(&submissionData); err != nil {
			jerrors.WriteErrorResponse(w, jerrors.BadRequest("invalid form data"))
			return
		}
	} else {
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			jerrors.WriteErrorResponse(w, jerrors.BadRequest("invalid form data"))
			return
		}

		keyValueForm := r.PostForm
		for key, value := range keyValueForm {
			// if there are more than one value, convert to comma separated string
			if len(value) > 1 {
				submissionData[key] = strings.Join(value, ",")
			} else {
				submissionData[key] = value[0]
			}
		}
	}

	// TODO: Implement actual form submission logic
	err = c.embedService.SubmitForm(formID, submissionData, r.Header)
	if err != nil {
		jerrors.WriteErrorResponse(w, err)
		return
	}

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
