package remotes

import (
	"bytes"
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

func (c JformClient) SubmitForm(
	ctx context.Context,
	formID uuid.UUID,
	data map[string]interface{},
	headers http.Header,
) error {
	url := fmt.Sprintf("%s/public/api/v1/forms/%s/submissions", c.URL, formID.String())

	body, err := json.Marshal(data)
	if err != nil {
		return jerrors.InternalServerError("failed to marshal data")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return jerrors.InternalServerError("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", c.APIKey)

	// get the headers comma separated for logging
	// forward the http headers from the client to the jform service
	for key, value := range headers {
		if key == "Content-Type" {
			continue
		}
		req.Header.Set(key, strings.Join(value, ", "))
	}
	slog.Info("sending request to jform service", "url", url, "headers", headers, "body", string(body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return jerrors.InternalServerError("failed to call jform service")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return jerrors.NotFound("form not found")
	case http.StatusBadRequest:
		return jerrors.BadRequest("invalid form ID")
	case http.StatusTooManyRequests:
		// parse the response body:
		errResponse := jerrors.ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return jerrors.InternalServerError("failed to decode response")
		}

		// usage has exceeded the limit
		return jerrors.TooManyRequests(errResponse.Error.Message)
	default:
		return jerrors.InternalServerError(fmt.Sprintf("jform service returned status %d", resp.StatusCode))
	}
}

// GetForm retrieves form data from the jform service
func (c JformClient) GetForm(ctx context.Context, formID uuid.UUID) (*dto.FormResponse, error) {
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
