// Froconnect Form Embed Script
// This script enables embeddable forms with server-side rendering and client-side validation
(function() {
    'use strict';

    // Base URL for API calls
    var API_BASE_URL = '{{.APIBaseURL}}';

    // ======================
    // 1. CSS INJECTION
    // ======================
    function injectStyles() {
        if (document.getElementById('froconnect-styles')) return;

        var style = document.createElement('style');
        style.id = 'froconnect-styles';
        style.textContent = `{{.CSSContent}}`;
        document.head.appendChild(style);
    }

    // ======================
    // 2. LANGUAGE DETECTION
    // ======================
    function detectLanguage(container) {
        // Priority 1: Attribute (user explicitly set)
        var attrLang = container.getAttribute('data-froconnect-lang');
        if (attrLang) {
            return attrLang;
        }

        // Priority 2: Browser language (de-DE → de)
        var browserLang = (navigator.language || navigator.userLanguage || '').split('-')[0];
        if (browserLang) {
            return browserLang;
        }

        // Priority 3: Default (let server decide)
        return '';
    }

    // ======================
    // 3. FORM HTML FETCHING
    // ======================
    function fetchFormHTML(embedId, lang) {
        var url = API_BASE_URL + '/api/public/v1/embeds/' + embedId + '/data';
        if (lang) {
            url += '?lang=' + encodeURIComponent(lang);
        }

        return fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'text/json; charset=utf-8'
            }
        })
        .then(function(response) {
            if (!response.ok) {
                throw new Error('Failed to load form (status: ' + response.status + ')');
            }
            // Get the language actually used by server from header
            var actualLang = response.headers.get('X-Form-Language');
            return response.json().then(function(resp) {
                var html = resp.data.html;
                var css = resp.data.css;
                return { html: html, css: css, language: actualLang };
            });
        });
    }

    // ======================
    // 4. FORM VALUE MANAGEMENT
    // ======================
    function getFormValues(container) {
        var form = container.querySelector('form');
        if (!form) return {};

        var values = {};
        var fields = form.querySelectorAll('input, textarea, select');
        fields.forEach(function(field) {
            var name = field.name || field.id;
            if (name) {
                if (field.type === 'checkbox') {
                    values[name] = field.checked;
                } else if (field.type === 'radio') {
                    if (field.checked) {
                        values[name] = field.value;
                    }
                } else {
                    values[name] = field.value;
                }
            }
        });
        return values;
    }

    function restoreFormValues(container, values) {
        var form = container.querySelector('form');
        if (!form || !values) return;

        Object.keys(values).forEach(function(name) {
            var field = form.querySelector('[name="' + name + '"], [id="' + name + '"]');
            if (field) {
                if (field.type === 'checkbox') {
                    field.checked = values[name];
                } else if (field.type === 'radio') {
                    if (field.value === values[name]) {
                        field.checked = true;
                    }
                } else {
                    field.value = values[name] || '';
                }
            }
        });
    }

    // ======================
    // 5. FORM RENDERING
    // ======================
    function renderForm(html, css, container, embedId, submissionUrl) {
        // Inject HTML (server-rendered form with dynamic CSS and HTML)
        container.innerHTML = html;

        // Add CSS
        var style = document.createElement('style');
        style.id = 'froconnect-dynamic-css';
        style.textContent = css;
        document.head.appendChild(style);

        // Initialize form handlers (validation + submission)
        var form = container.querySelector('form');
        var messageContainer = container.querySelector('[id^="jform-message-"]');

        if (form && messageContainer) {
            initializeFormHandlers(embedId, container, submissionUrl);
        }
    }

    // ======================
    // 6. VALIDATION FUNCTIONS
    // ======================
    function setupValidationListeners(form) {
        const fields = form.querySelectorAll('input, textarea, select');
        fields.forEach(function(field) {
            // Validate on blur
            field.addEventListener('blur', function() {
                var result = validateField(field);
                if (!result.valid) {
                    showFieldError(field, result.message);
                }
                updateSubmitButtonState(form);
            });

            // Clear errors and revalidate on input
            field.addEventListener('input', function() {
                clearFieldError(field);
                updateSubmitButtonState(form);
            });
        });
    }

    function validateField(field) {
        var value = field.value.trim();

        // Check required
        if (field.hasAttribute('required') && value === '') {
            var msg = field.getAttribute('data-error-required') || 'This field is required';
            return { valid: false, message: msg };
        }

        // If field is empty and not required, skip other validations
        if (value === '') {
            return { valid: true };
        }

        // Check minlength
        if (field.hasAttribute('minlength')) {
            var minLen = parseInt(field.getAttribute('minlength'));
            if (value.length < minLen) {
                var msg = field.getAttribute('data-error-minlength') || 'Too short';
                return { valid: false, message: msg };
            }
        }

        // Check maxlength
        if (field.hasAttribute('maxlength')) {
            var maxLen = parseInt(field.getAttribute('maxlength'));
            if (value.length > maxLen) {
                var msg = field.getAttribute('data-error-maxlength') || 'Too long';
                return { valid: false, message: msg };
            }
        }

        // Check email format
        if (field.type === 'email') {
            var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(value)) {
                var msg = field.getAttribute('data-error-email') || 'Invalid email address';
                return { valid: false, message: msg };
            }
        }

        // Check phone format
        if (field.type === 'tel' && field.hasAttribute('data-error-phone')) {
            var phoneRegex = /^[0-9+\-\s()]+$/;
            if (!phoneRegex.test(value)) {
                var msg = field.getAttribute('data-error-phone') || 'Invalid phone number';
                return { valid: false, message: msg };
            }
        }

        // Check pattern
        if (field.hasAttribute('pattern')) {
            var pattern = new RegExp(field.getAttribute('pattern'));
            if (!pattern.test(value)) {
                var msg = field.getAttribute('data-error-pattern') || 'Invalid format';
                return { valid: false, message: msg };
            }
        }

        // Check min for number fields
        if (field.hasAttribute('min') && field.type === 'number') {
            var minVal = parseFloat(field.getAttribute('min'));
            if (parseFloat(value) < minVal) {
                var msg = field.getAttribute('data-error-min') || 'Value too small';
                return { valid: false, message: msg };
            }
        }

        // Check max for number fields
        if (field.hasAttribute('max') && field.type === 'number') {
            var maxVal = parseFloat(field.getAttribute('max'));
            if (parseFloat(value) > maxVal) {
                var msg = field.getAttribute('data-error-max') || 'Value too large';
                return { valid: false, message: msg };
            }
        }

        return { valid: true };
    }

    function showFieldError(field, message) {
        // Remove any existing error
        clearFieldError(field);

        // Add invalid styling to field
        field.classList.add('jform-field-invalid', 'jform-shake');

        // Remove shake animation after it completes
        setTimeout(function() {
            field.classList.remove('jform-shake');
        }, 400);

        // Create and insert error message
        var errorDiv = document.createElement('div');
        errorDiv.className = 'jform-field-error-message';
        errorDiv.setAttribute('data-error-for', field.name || field.id);
        errorDiv.innerHTML =
            '<svg class="jform-field-error-icon" width="16" height="16" fill="currentColor" viewBox="0 0 20 20">' +
            '<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/></svg>' +
            '<span>' + escapeHtml(message) + '</span>';

        // Insert after the field
        field.parentNode.insertBefore(errorDiv, field.nextSibling);
    }

    function clearFieldError(field) {
        // Remove invalid styling
        field.classList.remove('jform-field-invalid', 'jform-shake');

        // Find and remove error message
        var errorMsg = field.parentNode.querySelector('[data-error-for="' + (field.name || field.id) + '"]');
        if (errorMsg) {
            errorMsg.remove();
        }
    }

    function updateSubmitButtonState(form) {
        var submitBtn = form.querySelector('button[type="submit"]');
        if (!submitBtn) return;

        var allValid = true;
        var fields = form.querySelectorAll('input, textarea, select');

        fields.forEach(function(field) {
            var result = validateField(field);
            if (!result.valid) {
                allValid = false;
            }
        });

        submitBtn.disabled = !allValid;
        if (!allValid) {
            submitBtn.classList.add('jform-btn-disabled');
        } else {
            submitBtn.classList.remove('jform-btn-disabled');
        }
    }

    function validateAllFields(form) {
        var isValid = true;
        var fields = form.querySelectorAll('input, textarea, select');
        var firstInvalidField = null;

        fields.forEach(function(field) {
            var result = validateField(field);
            if (!result.valid) {
                showFieldError(field, result.message);
                isValid = false;
                if (!firstInvalidField) {
                    firstInvalidField = field;
                }
            }
        });

        // Focus first invalid field
        if (firstInvalidField) {
            firstInvalidField.focus();
        }

        return isValid;
    }

    function escapeHtml(text) {
        var div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // ======================
    // 7. FORM SUBMISSION
    // ======================
    function initializeFormHandlers(embedId, container, submissionUrl) {
        const form = container.querySelector('form');
        const messageContainer = container.querySelector('[id^="jform-message-"]');

        if (!form) return;

        // Setup validation event listeners
        setupValidationListeners(form);

        // Initial button state check
        updateSubmitButtonState(form);

        form.addEventListener('submit', function(e) {
            e.preventDefault();
            handleFormSubmission(form, messageContainer, submissionUrl);
        });
    }

    function handleFormSubmission(form, messageContainer, submissionUrl) {
        // Get submit button
        const submitBtn = form.querySelector('button[type="submit"]');
        if (!submitBtn) return;

        // Clear any previous messages
        clearMessage(messageContainer);

        // Validate all fields before submission
        var isValid = validateAllFields(form);

        // If validation fails, stop here
        if (!isValid) {
            return;
        }

        // Add loading state
        form.classList.add('jform-submitting');
        submitBtn.classList.add('jform-btn-loading');
        submitBtn.classList.remove('jform-btn-disabled');
        submitBtn.disabled = true;

        // Store original button text
        const originalText = submitBtn.innerHTML;

        // Add spinner
        submitBtn.innerHTML = originalText + '<span class="jform-spinner"></span>';

        // Prepare form data
        const formData = new FormData(form);

        // Submit form
        fetch(submissionUrl, {
            method: 'POST',
            body: formData,
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        })
        .then(function(response) {
            if (!response.ok) {
                return response.json().then(function(err) {
                    throw new Error(err.message || 'Submission failed');
                });
            }
            return response.json();
        })
        .then(function(data) {
            // Success
            showMessage(messageContainer, 'success', data.message || 'Form submitted successfully!');
            form.reset();

            // Clear all field errors
            var fields = form.querySelectorAll('input, textarea, select');
            fields.forEach(function(field) {
                clearFieldError(field);
            });

            // Remove loading state after brief delay
            setTimeout(function() {
                removeLoadingState(form, submitBtn, originalText);
                // Update button state after reset
                updateSubmitButtonState(form);
            }, 500);
        })
        .catch(function(error) {
            // Error
            showMessage(messageContainer, 'error', error.message || 'An error occurred. Please try again.');
            removeLoadingState(form, submitBtn, originalText);
        });
    }

    function removeLoadingState(form, submitBtn, originalText) {
        form.classList.remove('jform-submitting');
        submitBtn.classList.remove('jform-btn-loading');
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
    }

    function showMessage(container, type, message) {
        if (!container) return;

        const icon = type === 'success'
            ? '<svg class="jform-message-icon" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>'
            : '<svg class="jform-message-icon" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/></svg>';

        container.innerHTML =
            '<div class="jform-message jform-message-' + type + '">' +
                icon +
                '<div class="jform-message-content">' + escapeHtml(message) + '</div>' +
            '</div>';

        container.classList.remove('jform-hidden');

        // Auto-hide success messages after 5 seconds
        if (type === 'success') {
            setTimeout(function() {
                if (container.firstElementChild) {
                    container.firstElementChild.classList.add('jform-fade-out');
                    setTimeout(function() {
                        clearMessage(container);
                    }, 300);
                }
            }, 5000);
        }
    }

    function clearMessage(container) {
        if (!container) return;
        container.innerHTML = '';
        container.classList.add('jform-hidden');
    }

    // ======================
    // 8. MUTATION OBSERVER (Auto-reload)
    // ======================
    function setupMutationObserver(embedId, container) {
        var reloadOnChange = container.getAttribute('data-froconnect-reload-on-change');
        if (reloadOnChange !== 'true') return;

        var observer = new MutationObserver(function(mutations) {
            mutations.forEach(function(mutation) {
                if (mutation.attributeName === 'data-froconnect-lang') {
                    var newLang = container.getAttribute('data-froconnect-lang');
                    console.log('[Froconnect] Language changed to:', newLang);
                    reloadForm(embedId, container);
                }
            });
        });

        observer.observe(container, {
            attributes: true,
            attributeFilter: ['data-froconnect-lang']
        });

        console.log('[Froconnect] Auto-reload enabled for form:', embedId);
    }

    // ======================
    // 9. FORM INITIALIZATION
    // ======================
    function initForm(embedId, container) {
        // Show loading state
        container.innerHTML = '<div class="froconnect-loading" style="padding: 2rem; text-align: center; color: #6b7280;">Loading form...</div>';

        // Detect language from attribute or browser
        var lang = detectLanguage(container);

        // Fetch form HTML (server renders with detected/requested language)
        fetchFormHTML(embedId, lang)
            .then(function(result) {
                // Update language attribute to reflect actual language used by server
                if (result.language) {
                    container.setAttribute('data-froconnect-lang', result.language);
                }

                // Render form (inject HTML and initialize handlers)
                var submissionUrl = API_BASE_URL + '/api/public/v1/embeds/' + embedId + '/submissions';
                renderForm(result.html, result.css,container, embedId, submissionUrl);

                // Setup mutation observer for auto-reload
                setupMutationObserver(embedId, container);

                console.log('[Froconnect] Form initialized:', embedId, 'Language:', result.language || lang);
            })
            .catch(function(error) {
                console.error('[Froconnect] Error loading form:', error);
                container.innerHTML = '<div class="froconnect-error" style="padding: 2rem; text-align: center; color: #dc2626; background: #fee2e2; border-radius: 0.5rem;">' +
                    '<strong>Unable to load form.</strong><br>Please try refreshing the page.' +
                    '</div>';
            });
    }

    function reloadForm(embedId, container) {
        var lang = container.getAttribute('data-froconnect-lang');

        // Save current values
        var values = getFormValues(container);

        // Show updating indicator
        var existingForm = container.querySelector('form');
        if (existingForm) {
            var loadingOverlay = document.createElement('div');
            loadingOverlay.style.cssText = 'position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: rgba(255,255,255,0.8); display: flex; align-items: center; justify-content: center; z-index: 1000;';
            loadingOverlay.innerHTML = '<div style="padding: 1rem; background: white; border-radius: 0.5rem; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">Updating...</div>';
            container.style.position = 'relative';
            container.appendChild(loadingOverlay);
        }

        // Fetch and re-render
        fetchFormHTML(embedId, lang)
            .then(function(result) {
                var submissionUrl = API_BASE_URL + '/api/public/v1/embeds/' + embedId + '/submissions';
                renderForm(result.html, result.css, container, embedId, submissionUrl);
                restoreFormValues(container, values);
                console.log('[Froconnect] Form reloaded with language:', lang);
            })
            .catch(function(error) {
                console.error('[Froconnect] Error reloading form:', error);
                // Remove loading overlay on error
                var overlay = container.querySelector('[style*="position: absolute"]');
                if (overlay) overlay.remove();
            });
    }

    // ======================
    // 10. PUBLIC API
    // ======================
    window.Froconnect = window.Froconnect || {};
    window.Froconnect.reload = function(embedId) {
        if (embedId) {
            // Reload specific form
            var container = document.querySelector('[data-froconnect-form="' + embedId + '"]');
            if (container) {
                reloadForm(embedId, container);
            } else {
                console.warn('[Froconnect] Form not found:', embedId);
            }
        } else {
            // Reload all forms
            var containers = document.querySelectorAll('[data-froconnect-form]');
            containers.forEach(function(container) {
                var id = container.getAttribute('data-froconnect-form');
                reloadForm(id, container);
            });
        }
    };

    // ======================
    // 11. AUTO-INITIALIZATION
    // ======================
    function initializeForms() {
        // Inject CSS once
        injectStyles();

        // Find and initialize all forms
        var containers = document.querySelectorAll('[data-froconnect-form]');
        console.log('[Froconnect] Found', containers.length, 'form(s) to initialize');

        containers.forEach(function(container) {
            var embedId = container.getAttribute('data-froconnect-form');
            if (embedId) {
                initForm(embedId, container);
            } else {
                console.warn('[Froconnect] Form container missing data-froconnect-form attribute');
            }
        });
    }

    // Wait for DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initializeForms);
    } else {
        initializeForms();
    }

    console.log('[Froconnect] Embed script loaded');
})();
