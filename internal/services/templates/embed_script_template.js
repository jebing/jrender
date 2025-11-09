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
