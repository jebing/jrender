package templates

import (
	"bytes"
	"html/template"
	"log/slog"

	"revonoir.com/jrender/controllers/dto/jerrors"
	"revonoir.com/jrender/internal/remotes/dto"
	"revonoir.com/jrender/internal/services/renders/dtos"
)

const CompleteHTMLTemplate = `<!DOCTYPE html>
<html lang="{{.DefaultLanguage}}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Name}}</title>
    <style>
        /* Production form styles - Clean and embed-friendly */
        
        /* CSS Reset for form elements */
        * {
            box-sizing: border-box;
        }
        
        body {
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #111827;
            background: transparent;
        }
        {{.CoreCSSStatic}}
        {{.CoreCSSDynamic}}
    </style>
</head>
<body>
    {{.CoreHTML}}
    {{.CoreJavascript}}
</body>
</html>`

type RequestWrapper struct {
	Name            string
	DefaultLanguage string
	CoreCSSStatic   template.CSS
	CoreCSSDynamic  template.CSS
	CoreHTML        template.HTML
	CoreJavascript  template.HTML
}

type FormCoreEngineIf interface {
	GenerateCSSStatic(data dtos.FormCoreData) (string, error)
	GenerateCSSDynamic(data dtos.FormCoreData) (string, error)
	GenerateHTML(data dtos.FormCoreData) (string, error)
	GenerateJavascript(data dtos.FormCoreData) (string, error)
}

type EmbeddedFormEngine struct {
	coreEngine      FormCoreEngineIf
	wrapperTemplate TemplateIf
}

func NewEmbeddedFormEngine(coreEngine FormCoreEngineIf, wrapperTemplate TemplateIf) *EmbeddedFormEngine {
	return &EmbeddedFormEngine{
		coreEngine:      coreEngine,
		wrapperTemplate: wrapperTemplate,
	}
}

func (fe EmbeddedFormEngine) GenerateHTML(data *dto.FormResponse) (string, error) {

	formData := dtos.FormData{
		Name:           data.Name,
		Description:    *data.Description,
		FormDefinition: data.FormDefinition,
		FormStyling:    data.FormStyling,
		FormID:         data.ID.String(),
	}
	// Create core data
	coreData := dtos.FormCoreData{
		FormData:        formData,
		DefaultLanguage: data.FormDefinition.Languages.Default,
		GridColumns:     12, // Default to 12-column grid
	}

	// Generate core CSS and HTML
	coreCSSDynamic, err := fe.coreEngine.GenerateCSSDynamic(coreData)
	if err != nil {
		return "", jerrors.InternalServerError("failed to generate core CSS dynamic")
	}

	coreCSSStatic, err := fe.coreEngine.GenerateCSSStatic(coreData)
	if err != nil {
		return "", jerrors.InternalServerError("failed to generate core CSS static")
	}

	coreHTML, err := fe.coreEngine.GenerateHTML(coreData)
	if err != nil {
		return "", jerrors.InternalServerError("failed to generate core HTML")
	}

	coreJavascript, err := fe.coreEngine.GenerateJavascript(coreData)
	if err != nil {
		return "", jerrors.InternalServerError("failed to generate core Javascript")
	}

	wrapper := RequestWrapper{
		Name:            data.Name,
		DefaultLanguage: coreData.DefaultLanguage,
		CoreCSSStatic:   template.CSS(coreCSSStatic),
		CoreCSSDynamic:  template.CSS(coreCSSDynamic),
		CoreHTML:        template.HTML(coreHTML),
		CoreJavascript:  template.HTML(coreJavascript),
	}

	var buf bytes.Buffer
	err = fe.wrapperTemplate.Execute(&buf, wrapper)
	if err != nil {
		slog.Error("failed to generate the complete HTML", "error", err)
		return "", jerrors.InternalServerError("failed to generate the complete HTML")
	}

	return buf.String(), nil
}

/**

    <script>
(function() {
    const form = document.querySelector('form');
    if (!form) return;

    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        // Check if reCAPTCHA is needed and generate token
        const tokenField = document.getElementById('g-recaptcha');
        if (tokenField && tokenField.value && typeof grecaptcha !== 'undefined') {
            const siteKey = tokenField.value;
            if (!siteKey) {
                console.error('reCAPTCHA site key is not set');
                return;
            }

            try {
                await new Promise((resolve, reject) => {
                    // Set timeout to prevent indefinite waiting
                    const timeout = setTimeout(() => {
                        reject(new Error('reCAPTCHA timeout: verification took too long'));
                    }, 5000);

                    try {
                        grecaptcha.ready(function() {
                            try {
                                grecaptcha.execute(siteKey, {action: 'submit'}).then(function(token) {
                                    clearTimeout(timeout);
                                    tokenField.value = token;
                                    resolve();
                                }).catch(function(error) {
                                    clearTimeout(timeout);
                                    reject(error);
                                });
                            } catch (error) {
                                clearTimeout(timeout);
                                reject(error);
                            }
                        });
                    } catch (error) {
                        clearTimeout(timeout);
                        reject(error);
                    }
                });
            } catch (error) {
                console.error('reCAPTCHA error:', error);

                // Determine error type for better user message
                let errorMessage = '<strong>Security verification failed.</strong> ';
                if (error.message && error.message.includes('Invalid site key')) {
                    errorMessage += 'Configuration error detected. Please contact support.';
                } else if (error.message && error.message.includes('timeout')) {
                    errorMessage += 'Verification timed out. Please try again.';
                } else {
                    errorMessage += 'Please refresh the page and try again.';
                }

                // Show error message to user
                const errorDiv = document.createElement('div');
                errorDiv.style.cssText = 'color: #dc2626; padding: 1rem; margin-bottom: 1rem; border: 1px solid #fecaca; background-color: #fef2f2; border-radius: 0.375rem;';
                errorDiv.innerHTML = errorMessage;
                form.insertBefore(errorDiv, form.firstChild);

                // Block form submission
                return;
            }
        }

        // Collect form data as JSON
        const formData = new FormData(form);
        const data = {};
        formData.forEach((value, key) => {
            if (data[key]) {
                data[key] = data[key] + ';' + value;
            } else {
                data[key] = value;
            }
        });

        // Submit via AJAX
        try {
            const response = await fetch(form.action, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify(data),
                credentials: 'include'
            });

            const result = await response.json();

            if (response.ok) {
                // Show success message
                form.innerHTML = '<div style="padding: 2rem; text-align: center; color: #059669;"><h3>Thank you!</h3><p>Your form has been submitted successfully.</p></div>';

                // Notify parent (optional - no error if no listener)
                if (window.parent !== window) {
                    window.parent.postMessage({
                        type: 'form-submitted',
                        success: true,
                        data: result
                    }, '*');
                }
            } else {
                // Show error
                const errorDiv = document.createElement('div');
                errorDiv.style.cssText = 'color: #dc2626; padding: 1rem; margin-bottom: 1rem; border: 1px solid #fecaca; background-color: #fef2f2; border-radius: 0.375rem;';
                errorDiv.textContent = 'Error: ' + (result.message || 'Failed to submit form');
                form.insertBefore(errorDiv, form.firstChild);
            }
        } catch (error) {
            // Network error - fallback to regular submission
            console.error('Form submission error:', error);
            form.submit();
        }

    });
})();
    </script>
*/
