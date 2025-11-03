package dto

import (
	"github.com/google/uuid"
	"revonoir.com/jrender/internal/services/renders/dtos"
)

// FormApiResponse represents the response wrapper from jform service
type FormApiResponse struct {
	Data FormResponse `json:"data"`
}

// FormResponse represents complete form data from jform service
type FormResponse struct {
	ID             uuid.UUID           `json:"id"`
	Name           string              `json:"name"`
	Description    *string             `json:"description"`
	FormDefinition dtos.FormDefinition `json:"form_definition"`
	FormStyling    dtos.FormStyling    `json:"form_styling,omitempty"`
}
