// Froconnect Shared Embedded Script between the Standalone / iframe embeds and the Embed JS


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
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ======================
// 7. FORM SUBMISSION
// ======================
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
                if (err?.error?.message) {
                    throw new Error(err.error.message);
                } else {
                    throw new Error('Submission failed');
                }
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
            container.firstElementChild.classList.add('jform-fade-out');
            setTimeout(function() {
                clearMessage(container);
            }, 300);
        }, 5000);
    }
}

function clearMessage(container) {
    if (!container) return;
    container.innerHTML = '';
    container.classList.add('jform-hidden');
}