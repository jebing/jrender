package remotes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"revonoir.com/jrender/conns/configs"
	"revonoir.com/jrender/controllers/dto/jerrors"
	"revonoir.com/jrender/internal/remotes/dto"
)

type JformClient struct {
	URL        string
	httpClient *http.Client
	APIKey     string
}

func NewJformClient(config configs.Configuration) *JformClient {
	return &JformClient{
		URL: config.Remote.JForm,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey: config.Remote.ApiKey,
	}
}

// GetForm retrieves form data from the jform service
func (c *JformClient) GetForm(ctx context.Context, formID uuid.UUID) (*dto.FormResponse, error) {
	url := fmt.Sprintf("%s/public/api/v1/forms/%s", c.URL, formID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, jerrors.InternalServerError("failed to create request")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "jrender/1.0")
	req.Header.Set("X-Api-Key", c.APIKey)

	// get the headers comma separated for logging
	headers := ""
	for key, value := range req.Header {
		headers += fmt.Sprintf("%s: %s, ", key, strings.Join(value, ", "))
	}
	slog.Info("sending request to jform service", "url", url, "headers", headers)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, jerrors.InternalServerError("failed to call jform service")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var apiResponse dto.FormApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			return nil, jerrors.InternalServerError("failed to decode response")
		}
		return &apiResponse.Data, nil
	case http.StatusNotFound:
		return nil, jerrors.NotFound("form not found")
	case http.StatusBadRequest:
		return nil, jerrors.BadRequest("invalid form ID")
	default:
		return nil, jerrors.InternalServerError(fmt.Sprintf("jform service returned status %d", resp.StatusCode))
	}
}
