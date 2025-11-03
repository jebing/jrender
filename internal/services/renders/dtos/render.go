package dtos

// FormPreviewData represents the data structure for form preview
type FormData struct {
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	FormDefinition FormDefinition `json:"form_definition"`
	FormStyling    FormStyling    `json:"form_styling"`
	FormID         string         `json:"form_id,omitempty"`
}

// FormCoreData represents the data structure for core form rendering
type FormCoreData struct {
	FormData
	DefaultLanguage string // Default language for the form
	GridColumns     int    // Number of grid columns (default: 12)
}

type FormLanguageSettings struct {
	Default   string   `json:"default" validate:"required"`
	Supported []string `json:"supported" validate:"required,min=1"`
}

type FormDefinition struct {
	Languages FormLanguageSettings `json:"languages" validate:"required"`
	Fields    []FormField          `json:"fields" validate:"required"`
}

// FormField represents a single form field with multi-language support
type FormField struct {
	ID                      string                     `json:"id" validate:"required"`
	Name                    string                     `json:"name" validate:"required"`
	Type                    string                     `json:"type" validate:"required,oneof=text email textarea select checkbox radio phone captcha heading paragraph submit_button divider spacer"`
	Required                bool                       `json:"required"`
	AllowMultipleSelections *bool                      `json:"allowMultipleSelections,omitempty"` // Only for select fields
	Options                 []FormFieldSelectOption    `json:"options,omitempty"`                 // Only for select/radio/checkbox fields
	Validation              *FormFieldValidation       `json:"validation,omitempty"`
	Translations            map[string]FormFieldTransl `json:"translations" validate:"required,min=1"`
	// Heading field properties
	Tag       string `json:"tag,omitempty"`       // HTML tag for heading fields (h1, h2, h3, etc.)
	Format    int    `json:"format,omitempty"`    // Format of the heading / paragraph (1: Bold, 2: Italic, 3: Bold Italic, 4: Underline, 5: Bold Underline, 6: Italic Underline, 7: Bold Italic, Underline)
	Alignment string `json:"alignment,omitempty"` // Alignment of the heading / paragraph (left, center, right)
}

type FormFieldValidation struct {
	Email            *bool    `json:"email,omitempty"` // Simple email validation flag
	Phone            *bool    `json:"phone,omitempty"` // Simple phone validation flag
	MinLength        *int     `json:"minLength,omitempty"`
	MaxLength        *int     `json:"maxLength,omitempty"`
	Step             *float64 `json:"step,omitempty"`
	MaxSize          *int64   `json:"maxSize,omitempty"`
	MaxFiles         *int     `json:"maxFiles,omitempty"`
	Accept           []string `json:"accept,omitempty"`
	MimeTypes        []string `json:"mimeTypes,omitempty"`
	MinSelected      *int     `json:"minSelected,omitempty"`
	MaxSelected      *int     `json:"maxSelected,omitempty"`
	MinDate          *string  `json:"minDate,omitempty"`
	MaxDate          *string  `json:"maxDate,omitempty"`
	AlphanumericOnly *bool    `json:"alphanumericOnly,omitempty"` // Only letters and numbers allowed
	NoSpecialChars   *bool    `json:"noSpecialChars,omitempty"`   // No special characters allowed
}

type FormFieldTransl struct {
	Label         string            `json:"label" validate:"required"`
	Placeholder   string            `json:"placeholder,omitempty"`
	Required      string            `json:"required,omitempty"`
	MinLength     string            `json:"minLength,omitempty"`
	MaxLength     string            `json:"maxLength,omitempty"`
	Email         string            `json:"email,omitempty"`
	Phone         string            `json:"phone,omitempty"`
	HelpText      string            `json:"helpText,omitempty"`
	ErrorMessages map[string]string `json:"errorMessages,omitempty"`
	Options       []FormFieldOption `json:"options,omitempty"`
}

type FormFieldOption struct {
	Value string `json:"value" validate:"required"`
	Label string `json:"label" validate:"required"`
}

// FormFieldSelectOption represents options for select/radio/checkbox fields with translations
type FormFieldSelectOption struct {
	Value        string            `json:"value" validate:"required"`
	Translations map[string]string `json:"translations" validate:"required,min=1"`
}

// FormStyling represents the complete canvas layout and styling configuration
type FormStyling struct {
	CanvasLayout CanvasLayout `json:"canvas_layout"`
	Styling      Styling      `json:"styling"`
}

// CanvasLayout represents the grid system configuration
type CanvasLayout struct {
	GridSystem            string            `json:"grid_system"`
	ResponsiveBreakpoints map[string]string `json:"responsive_breakpoints"`
	ContainerClasses      string            `json:"container_classes"`
	Rows                  []Row             `json:"rows"`
}

// Row represents a row in the canvas layout
type Row struct {
	ID      string   `json:"id"`
	Columns []Column `json:"columns"`
}

// Column represents a column within a row
type Column struct {
	ID                string           `json:"id"`
	ResponsiveSpans   map[string]int   `json:"responsive_spans"`
	ResponsiveClasses string           `json:"responsive_classes"`
	ColumnClasses     string           `json:"column_classes"`
	Gap               string           `json:"gap"`
	Fields            []FieldReference `json:"fields"`
}

// FieldReference represents a reference to a field by ID
type FieldReference struct {
	FieldID string `json:"field_id"`
}

// Styling represents the CSS styling configuration
type Styling struct {
	FormContainer FormContainerStyle        `json:"form_container"`
	LayoutDefault LayoutSettings            `json:"layout_default"` // Mandatory default layout
	FieldStyling  map[string]FieldTypeStyle `json:"field_styling"`
}

// FormContainerStyle represents styling for the form container
type FormContainerStyle struct {
	Classes string `json:"classes"`
}

// FieldTypeStyle represents styling for a specific field type
type FieldTypeStyle struct {
	Wrapper        string          `json:"wrapper,omitempty"`
	Label          string          `json:"label,omitempty"`
	Input          string          `json:"input,omitempty"`
	Element        string          `json:"element,omitempty"` // For heading fields
	Error          string          `json:"error,omitempty"`
	Button         string          `json:"button,omitempty"`          // For submit_button fields
	LayoutOverride *LayoutSettings `json:"layout_override,omitempty"` // Optional field-specific layout
}

// LayoutSettings represents layout configuration for form fields
type LayoutSettings struct {
	LabelLayout         string               `json:"label_layout,omitempty"` // "inline" or "stacked"
	InlineSettings      *InlineSettings      `json:"inline_settings,omitempty"`
	ResponsiveBehaviors *ResponsiveBehaviors `json:"responsive_behaviors,omitempty"`
}

// InlineSettings represents configuration for inline layout
type InlineSettings struct {
	LabelWidth     string `json:"label_width"`     // e.g., "30%"
	LabelAlignment string `json:"label_alignment"` // "left", "center", "right"
}

// ResponsiveBehaviors represents responsive layout behavior
type ResponsiveBehaviors struct {
	Mobile  string `json:"mobile"`  // "stacked" or "inline"
	Tablet  string `json:"tablet"`  // "stacked" or "inline"
	Desktop string `json:"desktop"` // "stacked" or "inline"
}
