package templates

import (
	"fmt"
	"html/template"
	"strings"

	"revonoir.com/jrender/internal/services/renders/dtos"
)

// FormSubmissionURLTemplate is the URL template for form submissions
const (
	FormSubmissionURLTemplate = "http://localhost:9000/api/public/v1/embeds/%s/submissions"

	// Input field types - fields that accept user input
	FieldTypeText     = "text"
	FieldTypeEmail    = "email"
	FieldTypePhone    = "phone"
	FieldTypeTextarea = "textarea"
	FieldTypeSelect   = "select"
	FieldTypeRadio    = "radio"
	FieldTypeCheckbox = "checkbox"
	FieldTypeCaptcha  = "captcha"

	// Metadata field types - fields that provide structure/layout/actions
	FieldTypeHeading      = "heading"
	FieldTypeParagraph    = "paragraph"
	FieldTypeSubmitButton = "submit_button"
	FieldTypeDivider      = "divider"
	FieldTypeSpacer       = "spacer"
)

// FormRenderer provides shared core form rendering functionality
type FormRenderer struct {
	funcMap        template.FuncMap
	captchaSiteKey string
}

// buildValidationAttributes generates data-error-* attributes and HTML5 validation attributes
// from field validation rules and translations
func buildValidationAttributes(field *dtos.FormField, translation dtos.FormFieldTransl) string {
	if field == nil {
		return ""
	}

	var attrs []string

	// Required validation
	if field.Required && translation.Required != "" {
		attrs = append(attrs, fmt.Sprintf(` data-error-required="%s"`, template.HTMLEscapeString(translation.Required)))
	}

	// Only add validation attributes if Validation is defined
	if field.Validation == nil {
		return strings.Join(attrs, "")
	}

	v := field.Validation

	// MinLength validation
	if v.MinLength != nil && *v.MinLength > 0 {
		attrs = append(attrs, fmt.Sprintf(` minlength="%d"`, *v.MinLength))
		if translation.MinLength != "" {
			attrs = append(attrs, fmt.Sprintf(` data-error-minlength="%s"`, template.HTMLEscapeString(translation.MinLength)))
		}
	}

	// MaxLength validation
	if v.MaxLength != nil && *v.MaxLength > 0 {
		attrs = append(attrs, fmt.Sprintf(` maxlength="%d"`, *v.MaxLength))
		if translation.MaxLength != "" {
			attrs = append(attrs, fmt.Sprintf(` data-error-maxlength="%s"`, template.HTMLEscapeString(translation.MaxLength)))
		}
	}

	// Email validation
	if v.Email != nil && *v.Email && translation.Email != "" {
		attrs = append(attrs, fmt.Sprintf(` data-error-email="%s"`, template.HTMLEscapeString(translation.Email)))
	}

	// Phone validation
	if v.Phone != nil && *v.Phone && translation.Phone != "" {
		attrs = append(attrs, fmt.Sprintf(` data-error-phone="%s"`, template.HTMLEscapeString(translation.Phone)))
	}

	return strings.Join(attrs, "")
}

// NewFormRenderer creates a new shared form renderer
func NewFormRenderer(captchaSiteKey string) *FormRenderer {
	renderer := &FormRenderer{
		captchaSiteKey: captchaSiteKey,
	}

	// Create shared function map that both preview and production can use
	renderer.funcMap = template.FuncMap{
		"getField":       getFieldByID,
		"getTranslation": getFieldTranslation,
		"renderField": func(field *dtos.FormField, translation dtos.FormFieldTransl, lang string, styling dtos.FormStyling) template.HTML {
			return renderer.RenderFieldHTML(field, translation, lang, styling)
		},
		"generateRowClasses":              generateRowClasses,
		"generateColumnClasses":           generateColumnClasses,
		"generateFieldCSS":                renderer.GenerateFieldCSS,
		"getFieldTypeStyle":               getFieldTypeStyle,
		"transformClasses":                transformTailwindClasses,
		"resolveFieldLayout":              resolveFieldLayout,
		"resolveAllResponsiveLayouts":     resolveAllResponsiveLayouts,
		"generateLayoutClasses":           generateLayoutClasses,
		"generateResponsiveLayoutClasses": generateResponsiveLayoutClasses,
		"safeCSS": func(css string) template.CSS {
			// Basic CSS sanitization - only allow safe characters
			sanitized := strings.ReplaceAll(css, "<", "")
			sanitized = strings.ReplaceAll(sanitized, ">", "")
			sanitized = strings.ReplaceAll(sanitized, "javascript:", "")
			sanitized = strings.ReplaceAll(sanitized, "expression(", "")
			return template.CSS(sanitized)
		},
		"getSubmissionURL": func(embedID string) string {
			return fmt.Sprintf(FormSubmissionURLTemplate, embedID)
		},
	}

	return renderer
}

// GetFuncMap returns the shared template function map
func (r *FormRenderer) GetFuncMap() template.FuncMap {
	return r.funcMap
}

func (r *FormRenderer) getFieldFormatClasses(fieldFormat int) string {
	var classes []string

	switch fieldFormat {
	case 1:
		classes = append(classes, "jform-font-bold")
	case 2:
		classes = append(classes, "jform-font-italic")
	case 3:
		classes = append(classes, "jform-font-bold jform-font-italic")
	case 4:
		classes = append(classes, "jform-underline")
	case 5:
		classes = append(classes, "jform-font-bold jform-underline")
	case 6:
		classes = append(classes, "jform-font-italic jform-underline")
	case 7:
		classes = append(classes, "jform-font-bold jform-font-italic jform-underline")
	}

	return strings.Join(classes, " ")
}

func (r *FormRenderer) getFieldAlignmentClasses(alignment string) string {
	var classes []string

	switch alignment {
	case "center":
		classes = append(classes, "jform-text-center")
	case "right":
		classes = append(classes, "jform-text-right")
	}

	return strings.Join(classes, " ")
}

var validTags = map[string]bool{"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true}

// Helper functions for common HTML structures

// getFieldByID finds a field by its ID
func getFieldByID(fields []dtos.FormField, fieldID string) *dtos.FormField {
	for i := range fields {
		if fields[i].ID == fieldID {
			return &fields[i]
		}
	}
	return nil
}

// getFieldTranslation gets translation for a field in specified language
func getFieldTranslation(field *dtos.FormField, lang string) dtos.FormFieldTransl {
	if field == nil {
		return dtos.FormFieldTransl{}
	}
	if trans, ok := field.Translations[lang]; ok {
		return trans
	}
	// Fallback to first available translation
	for _, trans := range field.Translations {
		return trans
	}
	return dtos.FormFieldTransl{}
}

func generateRowClasses(row dtos.Row) string {
	var classes []string

	// generate the grid  for each row
	classes = append(classes, "jform-grid jform-lg-grid-cols-12")

	return strings.Join(classes, " ")
}

// generateColumnClasses generates CSS classes for grid columns
func generateColumnClasses(column dtos.Column) string {
	var classes []string

	// Priority 1: Transform responsive classes if provided (converts Tailwind to jform classes)
	if column.ResponsiveClasses != "" {
		transformedClasses := transformTailwindClasses(column.ResponsiveClasses)
		if transformedClasses != "" {
			classes = append(classes, transformedClasses)
		}
	}

	// Priority 2: Transform column classes
	if column.ColumnClasses != "" {
		transformedClasses := transformTailwindClasses(column.ColumnClasses)
		if transformedClasses != "" {
			classes = append(classes, transformedClasses)
		}
	}

	// Priority 3: Generate native grid classes from spans if no responsive classes provided
	if column.ResponsiveSpans["xl"] > 0 {
		classes = append(classes, fmt.Sprintf("jform-xl-col-%d", column.ResponsiveSpans["xl"]))
	}
	if column.ResponsiveSpans["lg"] > 0 {
		classes = append(classes, fmt.Sprintf("jform-lg-col-%d", column.ResponsiveSpans["lg"]))
	}
	if column.ResponsiveSpans["md"] > 0 {
		classes = append(classes, fmt.Sprintf("jform-md-col-%d", column.ResponsiveSpans["md"]))
	}
	if column.ResponsiveSpans["sm"] > 0 {
		classes = append(classes, fmt.Sprintf("jform-sm-col-%d", column.ResponsiveSpans["sm"]))
	}
	// Default to full width on mobile
	classes = append(classes, "jform-col-12")

	return strings.Join(classes, " ")
}

// resolveFieldLayout resolves the effective layout for a field (default + field override + responsive)
func resolveFieldLayout(styling dtos.FormStyling, fieldType string, breakpoint string) dtos.LayoutSettings {
	// Start with global default
	effectiveLayout := styling.Styling.LayoutDefault

	// Apply field-specific override if exists
	if fieldStyling, exists := styling.Styling.FieldStyling[fieldType]; exists {
		if fieldStyling.LayoutOverride != nil {
			// Override main layout
			if fieldStyling.LayoutOverride.LabelLayout != "" {
				effectiveLayout.LabelLayout = fieldStyling.LayoutOverride.LabelLayout
			}
			// Override inline settings
			if fieldStyling.LayoutOverride.InlineSettings != nil {
				effectiveLayout.InlineSettings = fieldStyling.LayoutOverride.InlineSettings
			}
			// Override responsive behaviors
			if fieldStyling.LayoutOverride.ResponsiveBehaviors != nil {
				effectiveLayout.ResponsiveBehaviors = fieldStyling.LayoutOverride.ResponsiveBehaviors
			}
		}
	}

	// Apply responsive override based on breakpoint
	if effectiveLayout.ResponsiveBehaviors != nil {
		responsiveLayout := ""
		switch breakpoint {
		case "mobile":
			responsiveLayout = effectiveLayout.ResponsiveBehaviors.Mobile
		case "tablet":
			responsiveLayout = effectiveLayout.ResponsiveBehaviors.Tablet
		case "desktop":
			responsiveLayout = effectiveLayout.ResponsiveBehaviors.Desktop
		}

		if responsiveLayout != "" {
			effectiveLayout.LabelLayout = responsiveLayout
		}
	}

	// Ensure LabelLayout always has a fallback value
	if effectiveLayout.LabelLayout == "" {
		effectiveLayout.LabelLayout = "stacked" // Default to stacked layout
	}

	return effectiveLayout
}

// resolveBaseFieldLayout resolves the base (non-responsive) layout for a field
func resolveBaseFieldLayout(styling dtos.FormStyling, fieldType string) dtos.LayoutSettings {
	// Start with global default
	effectiveLayout := styling.Styling.LayoutDefault

	// Apply field-specific override if exists (but ignore responsive behaviors for base layout)
	if fieldStyling, exists := styling.Styling.FieldStyling[fieldType]; exists {
		if fieldStyling.LayoutOverride != nil {
			// Override main layout
			if fieldStyling.LayoutOverride.LabelLayout != "" {
				effectiveLayout.LabelLayout = fieldStyling.LayoutOverride.LabelLayout
			}
			// Override inline settings
			if fieldStyling.LayoutOverride.InlineSettings != nil {
				effectiveLayout.InlineSettings = fieldStyling.LayoutOverride.InlineSettings
			}
			// Note: We don't apply responsive behaviors here since this is base layout
		}
	}

	// Ensure LabelLayout always has a fallback value
	if effectiveLayout.LabelLayout == "" {
		effectiveLayout.LabelLayout = "stacked" // Default to stacked layout
	}

	return effectiveLayout
}

// hasResponsiveBehaviors checks if any responsive behaviors are actually defined
func hasResponsiveBehaviors(styling dtos.FormStyling, fieldType string) bool {
	// Check global responsive behaviors
	if styling.Styling.LayoutDefault.ResponsiveBehaviors != nil {
		globalBehaviors := styling.Styling.LayoutDefault.ResponsiveBehaviors
		if globalBehaviors.Mobile != "" || globalBehaviors.Tablet != "" || globalBehaviors.Desktop != "" {
			return true
		}
	}

	// Check field-specific responsive behaviors
	if fieldStyling, exists := styling.Styling.FieldStyling[fieldType]; exists {
		if fieldStyling.LayoutOverride != nil && fieldStyling.LayoutOverride.ResponsiveBehaviors != nil {
			fieldBehaviors := fieldStyling.LayoutOverride.ResponsiveBehaviors
			if fieldBehaviors.Mobile != "" || fieldBehaviors.Tablet != "" || fieldBehaviors.Desktop != "" {
				return true
			}
		}
	}

	return false
}

// resolveAllResponsiveLayouts resolves layouts for all breakpoints (mobile, tablet, desktop)
// Only generates responsive layouts if responsive behaviors are actually defined
func resolveAllResponsiveLayouts(styling dtos.FormStyling, fieldType string) map[string]dtos.LayoutSettings {
	if !hasResponsiveBehaviors(styling, fieldType) {
		// No responsive behaviors defined, return empty map
		return map[string]dtos.LayoutSettings{}
	}

	return map[string]dtos.LayoutSettings{
		"mobile":  resolveFieldLayout(styling, fieldType, "mobile"),
		"tablet":  resolveFieldLayout(styling, fieldType, "tablet"),
		"desktop": resolveFieldLayout(styling, fieldType, "desktop"),
	}
}

func generateLabelLayoutClasses(layout dtos.LayoutSettings) string {
	var classes []string

	// Add inline-specific classes
	if layout.LabelLayout == "inline" && layout.InlineSettings != nil {
		// Label width class
		widthClass := fmt.Sprintf("jform-label-width-%s", strings.ReplaceAll(layout.InlineSettings.LabelWidth, "%", ""))
		classes = append(classes, widthClass)

		// Label alignment class
		alignmentClass := fmt.Sprintf("jform-label-align-%s", layout.InlineSettings.LabelAlignment)
		classes = append(classes, alignmentClass)
	}

	return strings.Join(classes, " ")
}

// generateLayoutClasses generates CSS classes for a field based on its layout configuration
func generateLayoutClasses(layout dtos.LayoutSettings, fieldType string) string {
	var classes []string

	// Add base layout class
	baseClass := fmt.Sprintf("jform-layout-%s", layout.LabelLayout)
	classes = append(classes, baseClass)

	// Add field-type specific class
	fieldLayoutClass := fmt.Sprintf("jform-%s-layout-%s", fieldType, layout.LabelLayout)
	classes = append(classes, fieldLayoutClass)

	return strings.Join(classes, " ")
}

// generateResponsiveLayoutClasses generates CSS classes for responsive layouts across all breakpoints
// Returns empty string if no responsive layouts are provided (no responsive behaviors defined)
func generateResponsiveLayoutClasses(layouts map[string]dtos.LayoutSettings, fieldType string) string {
	if len(layouts) == 0 {
		// No responsive behaviors defined, return empty string
		return ""
	}

	var classes []string

	// Generate classes for each breakpoint
	for breakpoint, layout := range layouts {
		// Base responsive layout class
		responsiveClass := fmt.Sprintf("jform-%s-layout-%s", breakpoint, layout.LabelLayout)
		classes = append(classes, responsiveClass)

		// Field-type specific responsive class
		fieldResponsiveClass := fmt.Sprintf("jform-%s-%s-layout-%s", fieldType, breakpoint, layout.LabelLayout)
		classes = append(classes, fieldResponsiveClass)

		// Responsive inline-specific classes
		if layout.LabelLayout == "inline" && layout.InlineSettings != nil {
			// Responsive label width class
			responsiveWidthClass := fmt.Sprintf("jform-%s-label-width-%s", breakpoint, strings.ReplaceAll(layout.InlineSettings.LabelWidth, "%", ""))
			classes = append(classes, responsiveWidthClass)

			// Responsive label alignment class
			responsiveAlignClass := fmt.Sprintf("jform-%s-label-align-%s", breakpoint, layout.InlineSettings.LabelAlignment)
			classes = append(classes, responsiveAlignClass)
		}
	}

	return strings.Join(classes, " ")
}

// Map of Tailwind classes to jform classes
var classMap = map[string]string{
	// Grid system
	"grid":            "jform-grid",
	"grid-cols-1":     "",                      // Default mobile behavior
	"grid-cols-12":    "jform-lg-grid-cols-12", // Apply at lg breakpoint
	"lg:grid-cols-12": "jform-lg-grid-cols-12", // Responsive grid columns

	"col-span-1":     "",
	"col-span-2":     "",
	"col-span-3":     "",
	"col-span-4":     "",
	"col-span-5":     "",
	"col-span-6":     "",
	"col-span-7":     "",
	"col-span-8":     "",
	"col-span-9":     "",
	"col-span-10":    "",
	"col-span-11":    "",
	"col-span-12":    "",
	"lg:col-span-1":  "",
	"lg:col-span-2":  "",
	"lg:col-span-3":  "",
	"lg:col-span-4":  "",
	"lg:col-span-5":  "",
	"lg:col-span-6":  "",
	"lg:col-span-7":  "",
	"lg:col-span-8":  "",
	"lg:col-span-9":  "",
	"lg:col-span-10": "",
	"lg:col-span-11": "",
	"lg:col-span-12": "",
	"md:col-span-1":  "",
	"md:col-span-2":  "",
	"md:col-span-3":  "",
	"md:col-span-4":  "",
	"md:col-span-5":  "",
	"md:col-span-6":  "",
	"md:col-span-7":  "",
	"md:col-span-8":  "",
	"md:col-span-9":  "",
	"md:col-span-10": "",
	"md:col-span-11": "",
	"md:col-span-12": "",
	"sm:col-span-1":  "",
	"sm:col-span-2":  "",
	"sm:col-span-3":  "",
	"sm:col-span-4":  "",
	"sm:col-span-5":  "",
	"sm:col-span-6":  "",
	"sm:col-span-7":  "",
	"sm:col-span-8":  "",
	"sm:col-span-9":  "",
	"sm:col-span-10": "",
	"sm:col-span-11": "",
	"sm:col-span-12": "",

	// Layout and spacing
	"gap-4":        "jform-gap-4",
	"gap-6":        "jform-gap-6",
	"space-y-4":    "jform-space-y-4",
	"space-y-8":    "jform-space-y-8",
	"hidden":       "jform-hidden",
	"block":        "jform-block",
	"lg:block":     "jform-lg-block",
	"flex":         "jform-flex",
	"items-center": "jform-items-center",

	// Container and sizing
	"max-w-7xl": "jform-max-w-7xl",
	"mx-auto":   "jform-mx-auto",
	"p-6":       "jform-p-6",
	"sm:p-8":    "jform-sm-p-8",
	"w-full":    "jform-w-full",
	"h-fit":     "jform-h-fit",

	// Background and colors
	"bg-white":           "jform-bg-white",
	"bg-gray-50":         "jform-bg-gray-50",
	"bg-green-600":       "jform-bg-green-600",
	"hover:bg-green-700": "jform-hover-bg-green-700",
	"text-white":         "jform-text-white",
	"text-gray-700":      "jform-text-gray-700",
	"text-gray-900":      "jform-text-gray-900",
	"text-red-600":       "jform-text-red-600",
	"text-blue-600":      "jform-text-blue-600",

	// Border and radius
	"rounded-lg":      "jform-rounded-lg",
	"rounded":         "jform-rounded",
	"border":          "jform-border",
	"border-gray-300": "jform-border-gray-300",
	"border-b-2":      "jform-border-b-2",
	"border-gray-200": "jform-border-gray-200",

	// Typography
	"text-2xl":      "jform-text-2xl",
	"text-lg":       "jform-text-lg",
	"text-sm":       "jform-text-sm",
	"font-bold":     "jform-font-bold",
	"font-semibold": "jform-font-semibold",
	"font-medium":   "jform-font-medium",

	// Spacing
	"mb-2": "jform-mb-2",
	"mb-4": "jform-mb-4",
	"mb-6": "jform-mb-6",
	"mt-2": "jform-mt-2",
	"mt-8": "jform-mt-8",
	"px-4": "jform-px-4",
	"px-8": "jform-px-8",
	"py-3": "jform-py-3",
	"py-4": "jform-py-4",
	"pb-2": "jform-pb-2",
	"mr-2": "jform-mr-2",
	"h-4":  "jform-h-4",
	"w-4":  "jform-w-4",

	// Position
	"sticky": "jform-sticky",
	"top-4":  "jform-top-4",

	// Effects and interactions
	"focus:ring-2":          "jform-focus-ring-2",
	"focus:ring-blue-500":   "jform-focus-ring-blue-500",
	"focus:ring-green-500":  "jform-focus-ring-green-500",
	"focus:border-blue-500": "jform-focus-border-blue-500",
	"focus:outline-none":    "jform-focus-outline-none",
	"focus:ring-offset-2":   "jform-focus-ring-offset-2",
	"transition-colors":     "jform-transition-colors",
	"duration-200":          "jform-duration-200",
}

// transformTailwindClasses converts Tailwind CSS classes to our jform-prefixed classes
func transformTailwindClasses(classes string) string {
	if classes == "" {
		return ""
	}

	// Split classes and transform each one
	inputClasses := strings.Fields(classes)
	var transformedClasses []string

	for _, class := range inputClasses {
		if mappedClass, exists := classMap[class]; exists {
			if mappedClass != "" { // Skip empty mappings
				transformedClasses = append(transformedClasses, mappedClass)
			}
		} else {
			// Keep unknown classes as-is but don't prefix them
			// This handles custom classes that aren't Tailwind
			transformedClasses = append(transformedClasses, class)
		}
	}

	return strings.Join(transformedClasses, " ")
}

// getFieldTypeStyle returns styling configuration for a specific field type
func getFieldTypeStyle(fieldStyling map[string]any, fieldType string) map[string]string {
	if styling, ok := fieldStyling[fieldType]; ok {
		if stylingMap, ok := styling.(map[string]string); ok {
			return stylingMap
		}
	}
	return make(map[string]string)
}

// wrapFieldType wraps content in field type container
func wrapFieldType(fieldType, content string) template.HTML {
	fieldTypeClass := fmt.Sprintf("field-type-%s", fieldType)
	return template.HTML(fmt.Sprintf(`<div class="%s">%s</div>`, fieldTypeClass, content))
}

// wrapFormField wraps content in form-field with layout classes
func wrapFormField(layoutClasses, content string) string {
	return fmt.Sprintf(`<div class="form-field %s">%s</div>`, layoutClasses, content)
}

// renderLabel generates label HTML with required mark
func renderLabel(fieldID, labelClasses, label string, required bool) string {
	escapedID := template.HTMLEscapeString(fieldID)
	escapedLabel := template.HTMLEscapeString(label)
	requiredMark := ""
	if required {
		requiredMark = ` <span class="jform-text-red-600">*</span>`
	}

	if labelClasses != "" {
		return fmt.Sprintf(`<label for="%s" class="%s">%s%s</label>`,
			escapedID, labelClasses, escapedLabel, requiredMark)
	}
	return fmt.Sprintf(`<label for="%s">%s%s</label>`,
		escapedID, escapedLabel, requiredMark)
}

// renderFieldSet creates fieldset structure for radio/checkbox
func renderFieldSet(legend string, labelClasses string, required bool, options string) string {
	requiredMark := ""
	if required {
		requiredMark = ` <span class="jform-text-red-600">*</span>`
	}
	return fmt.Sprintf(`<legend class="%s">%s%s</legend>%s`,
		labelClasses,
		template.HTMLEscapeString(legend), requiredMark, options)
}

// renderFieldWrapper combines field type and form field wrappers
func renderFieldWrapper(fieldType, layoutClasses, content string) template.HTML {
	return wrapFieldType(fieldType, wrapFormField(layoutClasses, content))
}

// RenderFieldHTML generates HTML for a form field - core shared logic
func (r *FormRenderer) RenderFieldHTML(field *dtos.FormField, translation dtos.FormFieldTransl, currentLang string, styling dtos.FormStyling) template.HTML {
	if field == nil {
		return ""
	}

	// No responsive behaviors - resolve base layout and generate base classes
	primaryLayout := resolveBaseFieldLayout(styling, field.Type)
	baseLayoutClasses := generateLayoutClasses(primaryLayout, field.Type)
	labelLayoutClasses := generateLabelLayoutClasses(primaryLayout)

	// Resolve responsive layouts for all breakpoints
	responsiveLayouts := resolveAllResponsiveLayouts(styling, field.Type)
	responsiveLayoutClasses := generateResponsiveLayoutClasses(responsiveLayouts, field.Type)

	if len(responsiveLayouts) > 0 {
		// Has responsive behaviors - use desktop layout for structure but no base layout classes
		if desktopLayout, exists := responsiveLayouts["desktop"]; exists {
			primaryLayout = desktopLayout
		} else {
			// Fallback if desktop not defined, use mobile or tablet
			for _, layout := range responsiveLayouts {
				primaryLayout = layout
				break
			}
		}
	}

	// Combine all layout classes for the form-field element
	allLayoutClasses := strings.TrimSpace(baseLayoutClasses + " " + responsiveLayoutClasses)

	switch field.Type {
	case "heading":
		tag := "h2" // default
		if field.Tag != "" {
			// Validate tag to prevent XSS
			if validTags[field.Tag] {
				tag = field.Tag
			}
		}

		formatClasses := r.getFieldFormatClasses(field.Format)
		alignmentClasses := r.getFieldAlignmentClasses(field.Alignment)
		label := template.HTMLEscapeString(translation.Label)

		content := fmt.Sprintf(`<%s class="%s %s">%s</%s>`,
			tag, alignmentClasses, formatClasses, label, tag)
		return wrapFieldType(field.Type, content)

	case "paragraph":
		formatClasses := r.getFieldFormatClasses(field.Format)
		alignmentClasses := r.getFieldAlignmentClasses(field.Alignment)
		label := template.HTMLEscapeString(translation.Label)

		content := fmt.Sprintf(`<p class="%s %s">%s</p>`, alignmentClasses, formatClasses, label)
		return wrapFieldType(field.Type, content)

	case FieldTypeText, FieldTypeEmail, FieldTypePhone:
		inputType := "text"
		switch field.Type {
		case FieldTypeEmail:
			inputType = "email"
		case FieldTypePhone:
			inputType = "tel"
		}

		placeholder := ""
		if translation.Placeholder != "" {
			placeholder = fmt.Sprintf(` placeholder="%s"`, template.HTMLEscapeString(translation.Placeholder))
		}

		return r.renderInputField(field, translation, allLayoutClasses, labelLayoutClasses, primaryLayout, inputType, placeholder)

	case FieldTypeTextarea:
		placeholder := ""
		if translation.Placeholder != "" {
			placeholder = fmt.Sprintf(` placeholder="%s"`, template.HTMLEscapeString(translation.Placeholder))
		}

		return r.renderTextareaField(field, translation, allLayoutClasses, labelLayoutClasses, primaryLayout, placeholder)

	case FieldTypeSelect:
		required := ""
		if field.Required {
			required = " required"
		}

		// Build validation attributes
		validationAttrs := buildValidationAttributes(field, translation)

		// Build select options
		var options strings.Builder
		if translation.Placeholder != "" {
			options.WriteString(fmt.Sprintf(`<option value="" disabled selected>%s</option>`, template.HTMLEscapeString(translation.Placeholder)))
		}

		// Use passed current language
		for _, option := range field.Options {
			optionLabel := option.Value // Fallback to value
			if label, ok := option.Translations[currentLang]; ok {
				optionLabel = label
			} else {
				// Fallback to first available translation
				for _, label := range option.Translations {
					optionLabel = label
					break
				}
			}
			options.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`,
				template.HTMLEscapeString(option.Value),
				template.HTMLEscapeString(optionLabel)))
		}

		labelHTML := renderLabel(field.ID, "inline-label "+labelLayoutClasses, translation.Label, field.Required)
		content := fmt.Sprintf(`%s
			<select id="%s" name="%s"%s%s>%s</select>`,
			labelHTML,
			template.HTMLEscapeString(field.ID),
			template.HTMLEscapeString(field.Name),
			required,
			validationAttrs,
			options.String())
		return renderFieldWrapper(field.Type, allLayoutClasses, content)

	case FieldTypeRadio:
		required := ""
		if field.Required {
			required = " required"
		}

		// Build radio options
		var options strings.Builder
		options.WriteString(`<div>`)
		for i, option := range field.Options {
			optionLabel := option.Value // Fallback to value
			if label, ok := option.Translations[currentLang]; ok {
				optionLabel = label
			} else {
				// Fallback to first available translation
				for _, label := range option.Translations {
					optionLabel = label
					break
				}
			}

			optionID := fmt.Sprintf("%s_%d", field.ID, i)
			options.WriteString(fmt.Sprintf(`
				<div class="radio-option">
					<input type="radio" id="%s" name="%s" value="%s"%s>
					<label for="%s">%s</label>
				</div>`,
				template.HTMLEscapeString(optionID),
				template.HTMLEscapeString(field.Name),
				template.HTMLEscapeString(option.Value),
				required,
				template.HTMLEscapeString(optionID),
				template.HTMLEscapeString(optionLabel)))
		}
		options.WriteString(`</div>`)

		fieldSetHTML := renderFieldSet(translation.Label, "inline-label "+labelLayoutClasses, field.Required, options.String())
		content := wrapFormField(allLayoutClasses, fieldSetHTML)
		return wrapFieldType(field.Type, content)

	case FieldTypeCheckbox:
		// Build checkbox options
		var options strings.Builder
		options.WriteString(`<div>`)
		for i, option := range field.Options {
			optionLabel := option.Value // Fallback to value
			if label, ok := option.Translations[currentLang]; ok {
				optionLabel = label
			} else {
				// Fallback to first available translation
				for _, label := range option.Translations {
					optionLabel = label
					break
				}
			}

			optionID := fmt.Sprintf("%s_%d", field.ID, i)
			options.WriteString(fmt.Sprintf(`
				<div class="checkbox-option">
					<input type="checkbox" id="%s" name="%s" value="%s">
					<label for="%s">%s</label>
				</div>`,
				template.HTMLEscapeString(optionID),
				template.HTMLEscapeString(field.Name),
				template.HTMLEscapeString(option.Value),
				template.HTMLEscapeString(optionID),
				template.HTMLEscapeString(optionLabel)))
		}
		options.WriteString(`</div>`)

		fieldSetHTML := renderFieldSet(translation.Label, "inline-label "+labelLayoutClasses, false, options.String()) // Checkboxes typically don't have required on fieldset
		content := wrapFormField(allLayoutClasses, fieldSetHTML)
		return wrapFieldType(field.Type, content)

	case FieldTypeSubmitButton:
		content := fmt.Sprintf(`<button type="submit">%s</button>`,
			template.HTMLEscapeString(translation.Label))
		return renderFieldWrapper(field.Type, "", content)

	case FieldTypeCaptcha:
		content := fmt.Sprintf(`<div class="captcha-container">
				<input type="hidden" id="g-recaptcha" name="g-recaptcha" value="%s">
			</div>
			<script src="https://www.google.com/recaptcha/api.js?render=%s" defer></script>`,
			r.captchaSiteKey, r.captchaSiteKey)
		return renderFieldWrapper(field.Type, "", content)

	default:
		content := fmt.Sprintf(`<div>Unsupported field type: %s</div>`, template.HTMLEscapeString(field.Type))
		return wrapFieldType(field.Type, content)
	}
}

// renderInputField renders an input field with the specified layout
func (r *FormRenderer) renderInputField(field *dtos.FormField, translation dtos.FormFieldTransl, layoutClasses, labelLayoutClasses string, layout dtos.LayoutSettings, inputType, placeholder string) template.HTML {
	escapedID := template.HTMLEscapeString(field.ID)
	escapedName := template.HTMLEscapeString(field.Name)
	escapedLabel := template.HTMLEscapeString(translation.Label)
	requiredAttr := ""
	if field.Required {
		requiredAttr = " required"
	}

	// Build validation attributes from field.Validation and translation
	validationAttrs := buildValidationAttributes(field, translation)

	switch layout.LabelLayout {
	case "hidden":
		// Hidden layout: only input with placeholder, no label
		content := fmt.Sprintf(`<input type="%s" id="%s" name="%s"%s%s%s>`,
			inputType, escapedID, escapedName, placeholder, requiredAttr, validationAttrs)
		return renderFieldWrapper(field.Type, layoutClasses, content)

	case "floating":
		// Floating layout: input with floating label
		floatingPlaceholder := placeholder
		if floatingPlaceholder == "" {
			floatingPlaceholder = fmt.Sprintf(` placeholder="%s"`, escapedLabel)
		}
		labelHtml := renderLabel(field.ID, "floating-label "+labelLayoutClasses, translation.Label, field.Required)
		content := fmt.Sprintf(`<div class="floating-input-container">
				<input type="%s" id="%s" name="%s"%s%s%s>
				%s
			</div>`,
			inputType, escapedID, escapedName, floatingPlaceholder, requiredAttr, validationAttrs,
			labelHtml)
		return renderFieldWrapper(field.Type, layoutClasses, content)

	case "inline":
		// Inline layout: label and input side by side
		labelHTML := renderLabel(field.ID, "inline-label "+labelLayoutClasses, translation.Label, field.Required)
		content := fmt.Sprintf(`%s
			<div class="inline-input">
				<input type="%s" id="%s" name="%s"%s%s%s>
			</div>`,
			labelHTML, inputType, escapedID, escapedName, placeholder, requiredAttr, validationAttrs)
		return renderFieldWrapper(field.Type, layoutClasses, content)

	default: // "stacked" or fallback
		// Stacked layout: label above input (traditional)
		labelHTML := renderLabel(field.ID, "", translation.Label, field.Required)
		content := fmt.Sprintf(`%s
			<input type="%s" id="%s" name="%s"%s%s%s>`,
			labelHTML, inputType, escapedID, escapedName, placeholder, requiredAttr, validationAttrs)
		return renderFieldWrapper(field.Type, layoutClasses, content)
	}
}

// renderTextareaField renders a textarea field with the specified layout
func (r *FormRenderer) renderTextareaField(field *dtos.FormField, translation dtos.FormFieldTransl, layoutClasses, labelLayoutClasses string, layout dtos.LayoutSettings, placeholder string) template.HTML {
	escapedID := template.HTMLEscapeString(field.ID)
	escapedName := template.HTMLEscapeString(field.Name)
	escapedLabel := template.HTMLEscapeString(translation.Label)
	requiredAttr := ""
	if field.Required {
		requiredAttr = " required"
	}

	// Build validation attributes from field.Validation and translation
	validationAttrs := buildValidationAttributes(field, translation)

	switch layout.LabelLayout {
	case "hidden":
		// Hidden layout: only textarea with placeholder, no label
		content := fmt.Sprintf(`<textarea id="%s" name="%s"%s%s%s></textarea>`,
			escapedID, escapedName, placeholder, requiredAttr, validationAttrs)
		return renderFieldWrapper(field.Type, layoutClasses, content)

	case "floating":
		// Floating layout: textarea with floating label
		floatingPlaceholder := placeholder
		if floatingPlaceholder == "" {
			floatingPlaceholder = fmt.Sprintf(` placeholder="%s"`, escapedLabel)
		}
		labelHtml := renderLabel(field.ID, "floating-label "+labelLayoutClasses, translation.Label, field.Required)
		content := fmt.Sprintf(`<div class="floating-input-container">
				<textarea id="%s" name="%s"%s%s%s></textarea>
				%s
			</div>`,
			escapedID, escapedName, floatingPlaceholder, requiredAttr, validationAttrs,
			labelHtml)
		return renderFieldWrapper(field.Type, layoutClasses, content)

	case "inline":
		// Inline layout: label and textarea side by side
		labelHTML := renderLabel(field.ID, "inline-label "+labelLayoutClasses, translation.Label, field.Required)
		content := fmt.Sprintf(`%s
			<div class="inline-input">
				<textarea id="%s" name="%s"%s%s%s></textarea>
			</div>`,
			labelHTML, escapedID, escapedName, placeholder, requiredAttr, validationAttrs)
		return renderFieldWrapper(field.Type, layoutClasses, content)

	default: // "stacked" or fallback
		// Stacked layout: label above textarea (traditional)
		labelHTML := renderLabel(field.ID, "", translation.Label, field.Required)
		content := fmt.Sprintf(`%s
			<textarea id="%s" name="%s"%s%s%s></textarea>`,
			labelHTML, escapedID, escapedName, placeholder, requiredAttr, validationAttrs)
		return renderFieldWrapper(field.Type, layoutClasses, content)
	}
}

// GenerateFieldCSS generates CSS rules for field types based on styling configuration
func (r *FormRenderer) GenerateFieldCSS(formStyling dtos.FormStyling) template.CSS {
	var cssRules strings.Builder

	// Generate layout-specific CSS rules first
	cssRules.WriteString(generateLayoutCSS(formStyling.Styling.LayoutDefault))

	// Generate CSS for each field type
	for fieldType, fieldStyle := range formStyling.Styling.FieldStyling {
		// Generate CSS classes for this field type (unified for both preview and production)
		if fieldStyle.Wrapper != "" {
			// Convert Tailwind classes to CSS properties or use as-is for custom CSS
			wrapperCSS := convertTailwindClassesToCSS(fieldStyle.Wrapper)
			cssRules.WriteString(fmt.Sprintf(".field-type-%s .form-field { %s }\n", fieldType, sanitizeCSS(wrapperCSS)))
		}
		if fieldStyle.Label != "" {
			labelCSS := convertTailwindClassesToCSS(fieldStyle.Label)
			cssRules.WriteString(fmt.Sprintf(".field-type-%s .form-field label { %s }\n", fieldType, sanitizeCSS(labelCSS)))
		}
		if fieldStyle.Input != "" {
			inputCSS := convertTailwindClassesToCSS(fieldStyle.Input)
			cssRules.WriteString(fmt.Sprintf(".field-type-%s .form-field input, .field-type-%s .form-field textarea, .field-type-%s .form-field select { %s }\n", fieldType, fieldType, fieldType, sanitizeCSS(inputCSS)))
		}
		if fieldStyle.Element != "" {
			// For heading, paragraph and other element types
			elementCSS := convertTailwindClassesToCSS(fieldStyle.Element)
			cssRules.WriteString(fmt.Sprintf(".field-type-%s h1, .field-type-%s h2, .field-type-%s h3, .field-type-%s h4, .field-type-%s h5, .field-type-%s h6, .field-type-%s p { %s }\n",
				fieldType, fieldType, fieldType, fieldType, fieldType, fieldType, fieldType, sanitizeCSS(elementCSS)))
		}
		if fieldStyle.Error != "" {
			errorCSS := convertTailwindClassesToCSS(fieldStyle.Error)
			cssRules.WriteString(fmt.Sprintf(".field-type-%s .form-field .error { %s }\n", fieldType, sanitizeCSS(errorCSS)))
		}
		if fieldStyle.Button != "" {
			buttonCSS := convertTailwindClassesToCSS(fieldStyle.Button)
			cssRules.WriteString(fmt.Sprintf(".field-type-%s .form-field button { %s }\n", fieldType, sanitizeCSS(buttonCSS)))
		}

		// Generate layout override CSS if exists
		if fieldStyle.LayoutOverride != nil {
			cssRules.WriteString(generateFieldLayoutCSS(fieldType, *fieldStyle.LayoutOverride))
		}
	}

	return template.CSS(cssRules.String())
}

// generateLayoutCSS generates base CSS rules for layout system
func generateLayoutCSS(defaultLayout dtos.LayoutSettings) string {
	var cssRules strings.Builder

	// Generate CSS for inline layouts with different label widths (all breakpoints)
	labelWidths := []string{"25", "30", "40", "50"}
	breakpoints := []struct {
		name  string
		query string
	}{
		{"mobile", "@media (max-width: 767px)"},
		{"tablet", "@media (min-width: 768px) and (max-width: 1023px)"},
		{"desktop", "@media (min-width: 1024px)"},
	}

	for _, bp := range breakpoints {
		for _, width := range labelWidths {
			cssRules.WriteString(fmt.Sprintf(`
%s {
	/* %s - Inline layout with %s%% label width */
	.jform-%s-label-width-%s .form-field.jform-%s-layout-inline {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
	}
	.jform-%s-label-width-%s .form-field.jform-%s-layout-inline .inline-label {
		width: %s%%;
		flex-shrink: 0;
		padding-top: 0.75rem;
	}
	.jform-%s-label-width-%s .form-field.jform-%s-layout-inline .inline-input {
		width: %d%%;
		flex-grow: 1;
	}
	
	/* %s - Layout type overrides */
	.jform-%s-layout-stacked .form-field {
		display: block !important;
	}
	.jform-%s-layout-stacked .form-field .inline-label {
		width: auto !important;
		padding-top: 0 !important;
		margin-bottom: 0.5rem !important;
		display: block !important;
	}
	.jform-%s-layout-stacked .form-field .inline-input {
		width: 100%% !important;
	}
	
	.jform-%s-layout-inline .form-field {
		display: flex !important;
		align-items: flex-start !important;
		gap: 1rem !important;
	}
	
	.jform-%s-layout-floating .form-field {
		position: relative !important;
	}
	
	.jform-%s-layout-hidden .form-field label {
		display: none !important;
	}
}
`, bp.query, bp.name, width, bp.name, width, bp.name, bp.name, width, bp.name, width, bp.name, width, bp.name, 100-parseWidth(width), bp.name, bp.name, bp.name, bp.name, bp.name, bp.name, bp.name))
		}

		// Generate alignment CSS for this breakpoint
		alignments := []string{"left", "right", "center"}
		for _, align := range alignments {
			cssRules.WriteString(fmt.Sprintf(`
%s {
	.jform-%s-label-align-%s .inline-label {
		text-align: %s !important;
	}
}
`, bp.query, bp.name, align, align))
		}
	}

	return cssRules.String()
}

// generateFieldLayoutCSS generates CSS for field-specific layout overrides
func generateFieldLayoutCSS(fieldType string, layout dtos.LayoutSettings) string {
	var cssRules strings.Builder

	// Generate field-specific layout CSS based on layout type
	switch layout.LabelLayout {
	case "floating":
		cssRules.WriteString(fmt.Sprintf(`
/* Floating label for %s fields */
.field-type-%s .jform-layout-floating .floating-input-container {
	position: relative;
}
.field-type-%s .jform-layout-floating .floating-label {
	position: absolute;
	top: 0.75rem;
	left: 0.75rem;
	background: white;
	padding: 0 0.25rem;
	transition: all 0.2s ease-in-out;
	pointer-events: none;
	color: #6b7280;
}
.field-type-%s .jform-layout-floating input:focus + .floating-label,
.field-type-%s .jform-layout-floating input:not(:placeholder-shown) + .floating-label,
.field-type-%s .jform-layout-floating textarea:focus + .floating-label,
.field-type-%s .jform-layout-floating textarea:not(:placeholder-shown) + .floating-label {
	top: -0.5rem;
	font-size: 0.75rem;
	color: #3b82f6;
}
`, fieldType, fieldType, fieldType, fieldType, fieldType, fieldType, fieldType))

	case "hidden":
		cssRules.WriteString(fmt.Sprintf(`
/* Hidden label for %s fields */
.field-type-%s .jform-layout-hidden label {
	display: none;
}
`, fieldType, fieldType))
	}

	return cssRules.String()
}

// convertTailwindClassesToCSS converts Tailwind utility classes to actual CSS properties
func convertTailwindClassesToCSS(classes string) string {
	// If it looks like CSS properties already (contains colons), return as-is
	if strings.Contains(classes, ":") && (strings.Contains(classes, ";") || !strings.Contains(classes, " ")) {
		return classes
	}

	// Map of Tailwind classes to CSS properties
	cssMap := map[string]string{
		"mb-4":                  "margin-bottom: 1rem;",
		"mb-6":                  "margin-bottom: 1.5rem;",
		"mb-2":                  "margin-bottom: 0.5rem;",
		"mt-8":                  "margin-top: 2rem;",
		"mt-2":                  "margin-top: 0.5rem;",
		"block":                 "display: block;",
		"flex":                  "display: flex;",
		"items-center":          "align-items: center;",
		"text-sm":               "font-size: 0.875rem; line-height: 1.25rem;",
		"text-lg":               "font-size: 1.125rem; line-height: 1.75rem;",
		"text-2xl":              "font-size: 1.5rem; line-height: 2rem;",
		"font-medium":           "font-weight: 500;",
		"font-semibold":         "font-weight: 600;",
		"font-bold":             "font-weight: 700;",
		"text-gray-700":         "color: #374151;",
		"text-gray-900":         "color: #111827;",
		"text-red-600":          "color: #dc2626;",
		"text-blue-600":         "color: #2563eb;",
		"text-white":            "color: #ffffff;",
		"bg-green-600":          "background-color: #059669;",
		"hover:bg-green-700":    "background-color: #047857;", // Simplified, doesn't handle :hover state
		"bg-gray-50":            "background-color: #f9fafb;",
		"bg-white":              "background-color: #ffffff;",
		"w-full":                "width: 100%;",
		"px-4":                  "padding-left: 1rem; padding-right: 1rem;",
		"py-3":                  "padding-top: 0.75rem; padding-bottom: 0.75rem;",
		"py-4":                  "padding-top: 1rem; padding-bottom: 1rem;",
		"px-8":                  "padding-left: 2rem; padding-right: 2rem;",
		"p-6":                   "padding: 1.5rem;",
		"border":                "border-width: 1px;",
		"border-gray-300":       "border-color: #d1d5db;",
		"border-b-2":            "border-bottom-width: 2px;",
		"border-gray-200":       "border-color: #e5e7eb;",
		"rounded-lg":            "border-radius: 0.5rem;",
		"rounded":               "border-radius: 0.25rem;",
		"focus:ring-2":          "outline: none;", // Simplified focus state
		"focus:ring-blue-500":   "box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5);",
		"focus:border-blue-500": "border-color: #3b82f6;",
		"transition-colors":     "transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out;",
		"duration-200":          "transition-duration: 0.2s;",
		"mr-2":                  "margin-right: 0.5rem;",
		"h-4":                   "height: 1rem;",
		"w-4":                   "width: 1rem;",
		"pb-2":                  "padding-bottom: 0.5rem;",
	}

	// Split classes and convert each one
	inputClasses := strings.Fields(classes)
	var cssProperties []string

	for _, class := range inputClasses {
		if cssProperty, exists := cssMap[class]; exists {
			cssProperties = append(cssProperties, cssProperty)
		}
		// If not found, skip unknown Tailwind classes
	}

	return strings.Join(cssProperties, " ")
}

// sanitizeCSS performs basic CSS sanitization
func sanitizeCSS(css string) string {
	// Remove potentially dangerous CSS
	sanitized := strings.ReplaceAll(css, "javascript:", "")
	sanitized = strings.ReplaceAll(sanitized, "expression(", "")
	sanitized = strings.ReplaceAll(sanitized, "@import", "")
	sanitized = strings.ReplaceAll(sanitized, "<", "")
	sanitized = strings.ReplaceAll(sanitized, ">", "")
	return sanitized
}

// parseWidth converts width string to integer
func parseWidth(width string) int {
	switch width {
	case "25":
		return 25
	case "30":
		return 30
	case "40":
		return 40
	case "50":
		return 50
	default:
		return 30
	}
}
