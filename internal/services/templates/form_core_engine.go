package templates

import (
	"io"
	"strings"

	"revonoir.com/jrender/internal/services/renders/dtos"
)

// FormCoreCSSStatic contains only the static CSS (no template directives)
// Used for embed.js to inject universal styles
const FormCoreCSSStatic = `
        /* 12-Column Responsive Grid System */
        .jform-grid {
            display: grid;
            grid-template-columns: repeat(1, minmax(0, 1fr));
            gap: 1.5rem;
            width: 100%;
        }

        /* Default 12-column grid - applies when jform-lg-grid-cols-12 class is present */
        .jform-grid.jform-lg-grid-cols-12 {
            grid-template-columns: repeat(1, minmax(0, 1fr)); /* Mobile: single column */
        }

        /* Column spans for mobile-first approach */
        .jform-col-1 { grid-column: span 1 / span 1; }
        .jform-col-2 { grid-column: span 1 / span 1; } /* Mobile: stack */
        .jform-col-3 { grid-column: span 1 / span 1; }
        .jform-col-4 { grid-column: span 1 / span 1; }
        .jform-col-5 { grid-column: span 1 / span 1; }
        .jform-col-6 { grid-column: span 1 / span 1; }
        .jform-col-7 { grid-column: span 1 / span 1; }
        .jform-col-8 { grid-column: span 1 / span 1; }
        .jform-col-9 { grid-column: span 1 / span 1; }
        .jform-col-10 { grid-column: span 1 / span 1; }
        .jform-col-11 { grid-column: span 1 / span 1; }
        .jform-col-12 { grid-column: span 1 / span 1; }

        /* Small screens (640px and up) */
        @media (min-width: 640px) {
            .jform-grid { gap: 1.5rem; }
            .jform-sm-col-1 { grid-column: span 1 / span 1; }
            .jform-sm-col-2 { grid-column: span 2 / span 2; }
            .jform-sm-col-3 { grid-column: span 3 / span 3; }
            .jform-sm-col-4 { grid-column: span 4 / span 4; }
            .jform-sm-col-5 { grid-column: span 5 / span 5; }
            .jform-sm-col-6 { grid-column: span 6 / span 6; }
            .jform-sm-col-7 { grid-column: span 7 / span 7; }
            .jform-sm-col-8 { grid-column: span 8 / span 8; }
            .jform-sm-col-9 { grid-column: span 9 / span 9; }
            .jform-sm-col-10 { grid-column: span 10 / span 10; }
            .jform-sm-col-11 { grid-column: span 11 / span 11; }
            .jform-sm-col-12 { grid-column: span 12 / span 12; }
            .jform-sm-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)); }
        }

        /* Medium screens (768px and up) */
        @media (min-width: 768px) {
            .jform-grid { gap: 1.5rem; }
            .jform-md-col-1 { grid-column: span 1 / span 1; }
            .jform-md-col-2 { grid-column: span 2 / span 2; }
            .jform-md-col-3 { grid-column: span 3 / span 3; }
            .jform-md-col-4 { grid-column: span 4 / span 4; }
            .jform-md-col-5 { grid-column: span 5 / span 5; }
            .jform-md-col-6 { grid-column: span 6 / span 6; }
            .jform-md-col-7 { grid-column: span 7 / span 7; }
            .jform-md-col-8 { grid-column: span 8 / span 8; }
            .jform-md-col-9 { grid-column: span 9 / span 9; }
            .jform-md-col-10 { grid-column: span 10 / span 10; }
            .jform-md-col-11 { grid-column: span 11 / span 11; }
            .jform-md-col-12 { grid-column: span 12 / span 12; }
            .jform-md-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)); }
        }

        /* Large screens (1024px and up) */
        @media (min-width: 1024px) {
            .jform-grid { gap: 1.5rem; }
            .jform-lg-col-1 { grid-column: span 1 / span 1; }
            .jform-lg-col-2 { grid-column: span 2 / span 2; }
            .jform-lg-col-3 { grid-column: span 3 / span 3; }
            .jform-lg-col-4 { grid-column: span 4 / span 4; }
            .jform-lg-col-5 { grid-column: span 5 / span 5; }
            .jform-lg-col-6 { grid-column: span 6 / span 6; }
            .jform-lg-col-7 { grid-column: span 7 / span 7; }
            .jform-lg-col-8 { grid-column: span 8 / span 8; }
            .jform-lg-col-9 { grid-column: span 9 / span 9; }
            .jform-lg-col-10 { grid-column: span 10 / span 10; }
            .jform-lg-col-11 { grid-column: span 11 / span 11; }
            .jform-lg-col-12 { grid-column: span 12 / span 12; }
            .jform-lg-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)) !important; }

            /* Apply 12-column grid when class is present */
            .jform-grid.jform-lg-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)); }
        }

        /* Extra large screens (1280px and up) */
        @media (min-width: 1280px) {
            .jform-xl-col-1 { grid-column: span 1 / span 1; }
            .jform-xl-col-2 { grid-column: span 2 / span 2; }
            .jform-xl-col-3 { grid-column: span 3 / span 3; }
            .jform-xl-col-4 { grid-column: span 4 / span 4; }
            .jform-xl-col-5 { grid-column: span 5 / span 5; }
            .jform-xl-col-6 { grid-column: span 6 / span 6; }
            .jform-xl-col-7 { grid-column: span 7 / span 7; }
            .jform-xl-col-8 { grid-column: span 8 / span 8; }
            .jform-xl-col-9 { grid-column: span 9 / span 9; }
            .jform-xl-col-10 { grid-column: span 10 / span 10; }
            .jform-xl-col-11 { grid-column: span 11 / span 11; }
            .jform-xl-col-12 { grid-column: span 12 / span 12; }
        }

        /* Utility Classes for Common Tailwind Patterns */
        .jform-hidden { display: none; }
        .jform-block { display: block; }
        .jform-flex { display: flex; }
        .jform-items-center { align-items: center; }
        .jform-gap-4 { gap: 1rem; }
        .jform-gap-6 { gap: 1.5rem; }
        .jform-space-y-8 > * + * { margin-top: 2rem; }
        .jform-space-y-4 > * + * { margin-top: 1rem; }

        /* Responsive visibility utilities */
        @media (min-width: 1024px) {
            .jform-lg-block { display: block; }
            .jform-hidden.jform-lg-block { display: block; }
        }

        /* Container and layout utilities */
        .jform-max-w-7xl { max-width: 80rem; }
        .jform-mx-auto { margin-left: auto; margin-right: auto; }
        .jform-p-6 { padding: 1.5rem; }
        .jform-sm-p-8 { padding: 1.5rem; }
        @media (min-width: 640px) {
            .jform-sm-p-8 { padding: 2rem; }
        }
        .jform-bg-white { background-color: #ffffff; }
        .jform-bg-gray-50 { background-color: #f9fafb; }
        .jform-rounded-lg { border-radius: 0.5rem; }
        .jform-h-fit { height: fit-content; }
        .jform-sticky { position: sticky; }
        .jform-top-4 { top: 1rem; }

        /* Text and typography utilities */
        .jform-text-2xl { font-size: 1.5rem; line-height: 2rem; }
        .jform-font-bold { font-weight: 700; }
        .jform-text-gray-900 { color: #111827; }
        .jform-border-b-2 { border-bottom-width: 2px; }
        .jform-border-gray-200 { border-color: #e5e7eb; }
        .jform-pb-2 { padding-bottom: 0.5rem; }
        .jform-text-lg { font-size: 1.125rem; line-height: 1.75rem; }
        .jform-font-semibold { font-weight: 600; }
        .jform-font-italic { font-style: italic; }
        .jform-underline { text-decoration: underline; }

        /* Text alignment utilities */
        .jform-text-center { text-align: center; }
        .jform-text-right { text-align: right; }

        /* Form element utilities */
        .jform-mb-4 { margin-bottom: 1rem; }
        .jform-mb-6 { margin-bottom: 1.5rem; }
        .jform-block { display: block; }
        .jform-text-sm { font-size: 0.875rem; line-height: 1.25rem; }
        .jform-font-medium { font-weight: 500; }
        .jform-text-gray-700 { color: #374151; }
        .jform-mb-2 { margin-bottom: 0.5rem; }
        .jform-w-full { width: 100%; }
        .jform-px-4 { padding-left: 1rem; padding-right: 1rem; }
        .jform-py-3 { padding-top: 0.75rem; padding-bottom: 0.75rem; }
        .jform-border { border-width: 1px; }
        .jform-border-gray-300 { border-color: #d1d5db; }
        .jform-rounded-lg { border-radius: 0.5rem; }
        .jform-focus-ring-2:focus { outline: none; box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-focus-ring-blue-500:focus { outline: none; box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-focus-border-blue-500:focus { border-color: #3b82f6; }
        .jform-transition-colors { transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out; }
        .jform-duration-200 { transition-duration: 0.2s; }
        .jform-mt-2 { margin-top: 0.5rem; }
        .jform-text-red-600 { color: #dc2626; }
        .jform-mt-8 { margin-top: 2rem; }
        .jform-bg-green-600 { background-color: #059669; }
        .jform-hover-bg-green-700:hover { background-color: #047857; }
        .jform-text-white { color: #ffffff; }
        .jform-font-bold { font-weight: 700; }
        .jform-py-4 { padding-top: 1rem; padding-bottom: 1rem; }
        .jform-px-8 { padding-left: 2rem; padding-right: 2rem; }
        .jform-focus-outline-none:focus { outline: none; }
        .jform-focus-ring-2:focus { box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-focus-ring-green-500:focus { box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.5); }
        .jform-focus-ring-offset-2:focus { box-shadow: 0 0 0 2px transparent, 0 0 0 4px rgba(16, 185, 129, 0.5); }
        .jform-text-lg { font-size: 1.125rem; line-height: 1.75rem; }
        .jform-mr-2 { margin-right: 0.5rem; }
        .jform-h-4 { height: 1rem; }
        .jform-w-4 { width: 1rem; }
        .jform-text-blue-600 { color: #2563eb; }
        .jform-focus-ring-blue-500:focus { box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-border-gray-300 { border-color: #d1d5db; }
        .jform-rounded { border-radius: 0.25rem; }

        /* Default form styling (fallbacks) */
        .form-container {
            background: transparent;
            padding: 1rem;
            border-radius: 8px;
        }

        .form-field {
            margin-bottom: 1.5rem;
        }

        .form-field label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 500;
            color: #374151;
        }

        .form-field input,
        .form-field textarea {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #d1d5db;
            border-radius: 4px;
            font-size: 1rem;
            transition: border-color 0.2s;
        }

        .form-field input:focus,
        .form-field textarea:focus {
            outline: none;
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }

        .form-field textarea {
            min-height: 100px;
            resize: vertical;
        }

        h1, h2, h3, h4, h5, h6 {
            margin: 0 0 1rem 0;
            font-weight: 600;
            line-height: 1.2;
        }

        h1 { font-size: 2.5rem; }
        h2 { font-size: 2rem; }
        h3 { font-size: 1.75rem; }
        h4 { font-size: 1.5rem; }
        h5 { font-size: 1.25rem; }
        h6 { font-size: 1.125rem; }

        p {
            margin: 0 0 1rem 0;
            line-height: 1.6;
            color: #6b7280;
        }

        /* Form field element styles */
        select {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #d1d5db;
            border-radius: 4px;
            font-size: 1rem;
            transition: border-color 0.2s;
            background-color: white;
        }

        select:focus {
            outline: none;
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }

        fieldset {
            border: none;
            padding: 0;
            margin: 0;
        }

        legend {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 500;
            color: #374151;
            padding: 0;
        }

        .radio-option, .checkbox-option {
            display: flex;
            align-items: center;
            margin-bottom: 0.5rem;
        }

        .radio-option input, .checkbox-option input {
            width: auto;
            margin-right: 0.5rem;
        }

        .radio-option label, .checkbox-option label {
            margin-bottom: 0;
            cursor: pointer;
        }

        button[type="submit"] {
            width: 100%;
            background-color: #3b82f6;
            color: white;
            padding: 0.75rem 1.5rem;
            border: none;
            border-radius: 4px;
            font-size: 1rem;
            cursor: pointer;
            transition: background-color 0.2s;
        }

        button[type="submit"]:hover {
            background-color: #2563eb;
        }

        .captcha-placeholder {
            border: 2px dashed #d1d5db;
            border-radius: 4px;
            padding: 1rem;
            text-align: center;
            background-color: #f9fafb;
        }

        .captcha-box {
            color: #6b7280;
        }

        .captcha-box span {
            display: block;
            margin-bottom: 0.5rem;
            font-size: 1rem;
        }

        .captcha-box small {
            font-size: 0.875rem;
            color: #9ca3af;
        }

        /* Form Submission Animation Styles */

        /* Button loading state */
        .jform-btn-loading {
            opacity: 0.7;
            cursor: not-allowed;
            position: relative;
            padding-right: 3rem;
        }

        /* Button disabled state */
        .jform-btn-disabled {
            opacity: 0.5;
            cursor: not-allowed !important;
            background-color: #9ca3af !important;
        }

        .jform-btn-disabled:hover {
            background-color: #9ca3af !important;
        }

        /* Spinner animation */
        .jform-spinner {
            display: inline-block;
            width: 1rem;
            height: 1rem;
            border: 2px solid rgba(255, 255, 255, 0.3);
            border-top-color: white;
            border-radius: 50%;
            animation: jform-spin 0.6s linear infinite;
            position: absolute;
            right: 1.5rem;
            top: 50%;
            margin-top: -0.5rem;
        }

        @keyframes jform-spin {
            to { transform: rotate(360deg); }
        }

        /* Form submitting state - dims and disables the form */
        .jform-submitting {
            opacity: 0.6;
            pointer-events: none;
            transition: opacity 0.3s ease;
        }

        .jform-submitting input,
        .jform-submitting textarea,
        .jform-submitting select,
        .jform-submitting button {
            cursor: not-allowed;
        }

        /* Message container styles */
        .jform-message {
            padding: 1rem;
            border-radius: 0.5rem;
            margin-bottom: 1rem;
            animation: jform-slide-down 0.3s ease-out;
            display: flex;
            align-items: flex-start;
            gap: 0.75rem;
        }

        .jform-message-success {
            background-color: #d1fae5;
            color: #065f46;
            border: 1px solid #10b981;
        }

        .jform-message-error {
            background-color: #fee2e2;
            color: #991b1b;
            border: 1px solid #ef4444;
        }

        .jform-message-info {
            background-color: #dbeafe;
            color: #1e40af;
            border: 1px solid #3b82f6;
        }

        .jform-message-icon {
            flex-shrink: 0;
            width: 1.25rem;
            height: 1.25rem;
            margin-top: 0.125rem;
        }

        .jform-message-content {
            flex-grow: 1;
        }

        /* Slide down animation for messages */
        @keyframes jform-slide-down {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        /* Fade out animation */
        .jform-fade-out {
            animation: jform-fade-out 0.3s ease-out forwards;
        }

        @keyframes jform-fade-out {
            to {
                opacity: 0;
                transform: translateY(-10px);
            }
        }

        /* Additional utility classes for animations */
        .jform-transition-all {
            transition: all 0.3s ease;
        }

        .jform-opacity-0 { opacity: 0; }
        .jform-opacity-50 { opacity: 0.5; }
        .jform-opacity-60 { opacity: 0.6; }
        .jform-opacity-100 { opacity: 1; }

        .jform-scale-95 { transform: scale(0.95); }
        .jform-scale-100 { transform: scale(1); }

        .jform-pointer-events-none { pointer-events: none; }
        .jform-pointer-events-auto { pointer-events: auto; }

        /* Hidden utility */
        .jform-hidden { display: none; }

        /* Field Validation Styles */

        /* Invalid field styling - red border and shadow */
        .jform-field-invalid {
            border-color: #ef4444 !important;
            box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1) !important;
        }

        /* Error message below field */
        .jform-field-error-message {
            color: #dc2626;
            font-size: 0.875rem;
            margin-top: 0.5rem;
            display: flex;
            align-items: flex-start;
            gap: 0.25rem;
            animation: jform-slide-down 0.2s ease-out;
        }

        .jform-field-error-icon {
            flex-shrink: 0;
            margin-top: 0.125rem;
        }

        /* Shake animation for invalid fields */
        @keyframes jform-shake {
            0%, 100% { transform: translateX(0); }
            10%, 30%, 50%, 70%, 90% { transform: translateX(-4px); }
            20%, 40%, 60%, 80% { transform: translateX(4px); }
        }

        .jform-shake {
            animation: jform-shake 0.4s ease-in-out;
        }

        /* Layout System Base Styles */

        /* Stacked Layout (default) */
        .form-field.jform-layout-stacked {
            display: block;
        }
        .form-field.jform-layout-stacked label {
            display: block;
            margin-bottom: 0.5rem;
        }

        /* Inline Layout */
        .form-field.jform-layout-inline {
            display: flex;
            align-items: flex-start;
            gap: 1rem;
        }
        .form-field.jform-layout-inline .inline-label {
            flex-shrink: 0;
            padding-top: 0.75rem;
            margin-bottom: 0;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-25 {
            width: 25%;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-30 {
            width: 30%;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-40 {
            width: 40%;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-50 {
            width: 50%;
        }
        .form-field.jform-layout-inline .inline-input {
            flex-grow: 1;
        }
        .form-field.jform-layout-inline .inline-input input,
        .form-field.jform-layout-inline .inline-input textarea,
        .form-field.jform-layout-inline .inline-input select {
            width: 100%;
        }

        /* Floating Layout */
        .form-field.jform-layout-floating {
            position: relative;
        }
        .form-field.jform-layout-floating .floating-input-container {
            position: relative;
        }
        .form-field.jform-layout-floating .floating-label {
            position: absolute;
            top: 0.75rem;
            left: 0.75rem;
            background: white;
            padding: 0 0.25rem;
            transition: all 0.2s ease-in-out;
            pointer-events: none;
            color: #6b7280;
            font-size: 1rem;
            z-index: 1;
        }
        .form-field.jform-layout-floating input:focus + .floating-label,
        .form-field.jform-layout-floating input:not(:placeholder-shown) + .floating-label,
        .form-field.jform-layout-floating textarea:focus + .floating-label,
        .form-field.jform-layout-floating textarea:not(:placeholder-shown) + .floating-label {
            top: -0.5rem;
            font-size: 0.75rem;
            color: #3b82f6;
        }
        .form-field.jform-layout-floating input,
        .form-field.jform-layout-floating textarea {
            padding-top: 1.5rem;
            padding-bottom: 0.5rem;
        }

        /* Hidden Layout */
        .form-field.jform-layout-hidden label {
            display: none !important;
        }

        /* Responsive Layout Behavior - Direct Selectors (Option 1) */
        @media (max-width: 767px) {
            /* Mobile responsive layout overrides */
            .form-field.jform-mobile-layout-stacked {
                display: block;
            }
            .form-field.jform-mobile-layout-stacked label {
                display: block;
                margin-bottom: 0.5rem;
            }

            .form-field.jform-mobile-layout-inline {
                display: flex;
                align-items: flex-start;
                gap: 1rem;
            }
            .form-field.jform-mobile-layout-inline .inline-label {
                flex-shrink: 0;
                padding-top: 0.75rem;
                margin-bottom: 0;
            }
            .form-field.jform-mobile-layout-inline .inline-input {
                flex-grow: 1;
            }

            .form-field.jform-mobile-layout-floating {
                position: relative;
            }
            .form-field.jform-mobile-layout-floating .floating-input-container {
                position: relative;
            }
            .form-field.jform-mobile-layout-floating .floating-label {
                position: absolute;
                top: 0.75rem;
                left: 0.75rem;
                background: white;
                padding: 0 0.25rem;
                transition: all 0.2s ease-in-out;
                pointer-events: none;
                color: #6b7280;
                font-size: 1rem;
                z-index: 1;
            }

            .form-field.jform-mobile-layout-hidden label {
                display: none !important;
            }
        }

        @media (min-width: 768px) and (max-width: 1023px) {
            /* Tablet responsive layout overrides */
            .form-field.jform-tablet-layout-stacked {
                display: block;
            }
            .form-field.jform-tablet-layout-stacked label {
                display: block;
                margin-bottom: 0.5rem;
            }

            .form-field.jform-tablet-layout-inline {
                display: flex;
                align-items: flex-start;
                gap: 1rem;
            }
            .form-field.jform-tablet-layout-inline .inline-label {
                flex-shrink: 0;
                padding-top: 0.75rem;
                margin-bottom: 0;
            }
            .form-field.jform-tablet-layout-inline .inline-input {
                flex-grow: 1;
            }

            .form-field.jform-tablet-layout-floating {
                position: relative;
            }
            .form-field.jform-tablet-layout-floating .floating-input-container {
                position: relative;
            }
            .form-field.jform-tablet-layout-floating .floating-label {
                position: absolute;
                top: 0.75rem;
                left: 0.75rem;
                background: white;
                padding: 0 0.25rem;
                transition: all 0.2s ease-in-out;
                pointer-events: none;
                color: #6b7280;
                font-size: 1rem;
                z-index: 1;
            }

            .form-field.jform-tablet-layout-hidden label {
                display: none !important;
            }
        }

        @media (min-width: 1024px) {
            /* Desktop responsive layout overrides */
            .form-field.jform-desktop-layout-stacked {
                display: block;
            }
            .form-field.jform-desktop-layout-stacked label {
                display: block;
                margin-bottom: 0.5rem;
            }

            .form-field.jform-desktop-layout-inline {
                display: flex;
                align-items: flex-start;
                gap: 1rem;
            }
            .form-field.jform-desktop-layout-inline .inline-label {
                flex-shrink: 0;
                padding-top: 0.75rem;
                margin-bottom: 0;
            }
            .form-field.jform-desktop-layout-inline .inline-input {
                flex-grow: 1;
            }

            .form-field.jform-desktop-layout-floating {
                position: relative;
            }
            .form-field.jform-desktop-layout-floating .floating-input-container {
                position: relative;
            }
            .form-field.jform-desktop-layout-floating .floating-label {
                position: absolute;
                top: 0.75rem;
                left: 0.75rem;
                background: white;
                padding: 0 0.25rem;
                transition: all 0.2s ease-in-out;
                pointer-events: none;
                color: #6b7280;
                font-size: 1rem;
                z-index: 1;
            }

            .form-field.jform-desktop-layout-hidden label {
                display: none !important;
            }
        }

        /* Override theme styles for reliability */
        .form-container {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif !important;
            line-height: 1.6 !important;
            color: #111827 !important;
        }

        input, select, textarea {
            font-family: inherit !important;
            font-size: 1rem !important;
            line-height: 1.5 !important;
        }`

const FormCoreCSSDynamic = `
    {{generateFieldCSS .FormStyling}}
`

// formCoreHTMLTemplate contains the actual form structure
const FormCoreHTMLTemplate = `
<div class="form-container {{transformClasses .FormStyling.Styling.FormContainer.Classes}}">
    <!-- Message container for submission feedback -->
    <div id="jform-message-{{.FormID}}" class="jform-hidden"></div>

    <!-- Fallback notice for users without JavaScript -->
    <noscript>
        <div class="jform-message jform-message-info">
            <svg class="jform-message-icon" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
            </svg>
            <div class="jform-message-content">
                <strong>JavaScript is disabled.</strong> The form will still work, but the page will reload when you submit.
            </div>
        </div>
    </noscript>

    <form action="{{getSubmissionURL .FormID}}" method="POST" class="{{transformClasses .FormStyling.CanvasLayout.ContainerClasses}}" data-jform-id="{{.FormID}}">
        {{range .FormStyling.CanvasLayout.Rows}}
        <div id="{{.ID}}" class="{{generateRowClasses .}}">
            {{range .Columns}}
            <div id="{{.ID}}" class="{{generateColumnClasses .}}">
                {{range .Fields}}
                    {{$field := getField $.FormDefinition.Fields .FieldID}}
                    {{if $field}}
                        {{$translation := getTranslation $field $.DefaultLanguage}}
                        {{renderField $field $translation $.DefaultLanguage $.FormStyling}}
                    {{end}}
                {{end}}
            </div>
            {{end}}
        </div>
        {{end}}
    </form>    
</div>`

const FormCoreJsTemplate = `<script type="text/javascript">
    (function() {
        'use strict';

        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', initFormSubmission);
        } else {
            initFormSubmission();
        }

        function initFormSubmission() {
            const formId = '{{.FormID}}';
            const form = document.querySelector('[data-jform-id="' + formId + '"]');
            const messageContainer = document.getElementById('jform-message-' + formId);

            if (!form) return;

            // Setup validation event listeners
            setupValidationListeners(form);

            // Initial button state check
            updateSubmitButtonState(form);

            form.addEventListener('submit', function(e) {
                e.preventDefault();
                handleFormSubmission(form, messageContainer);
            });
        }

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

            // Check phone format (basic check for digits)
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

        function handleFormSubmission(form, messageContainer) {
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
            const actionURL = form.getAttribute('action');

            // Submit form
            fetch(actionURL, {
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

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
    })();
</script>`

type TemplateIf interface {
	Execute(wr io.Writer, data any) error
}

// FormCoreEngine provides the shared core form rendering functionality
type FormCoreEngine struct {
	cssTemplate  TemplateIf
	htmlTemplate TemplateIf
	jsTemplate   TemplateIf
}

func NewFormCoreEngine(cssTemplate TemplateIf, htmlTemplate TemplateIf, jsTemplate TemplateIf) *FormCoreEngine {
	return &FormCoreEngine{
		cssTemplate:  cssTemplate,
		htmlTemplate: htmlTemplate,
		jsTemplate:   jsTemplate,
	}
}

func (fce FormCoreEngine) GenerateCSSStatic(data dtos.FormCoreData) (string, error) {
	return FormCoreCSSStatic, nil
}

// GenerateCSS generates the core CSS styles
func (fce FormCoreEngine) GenerateCSSDynamic(data dtos.FormCoreData) (string, error) {
	var buf strings.Builder
	err := fce.cssTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateHTML generates the core form HTML
func (fce FormCoreEngine) GenerateHTML(data dtos.FormCoreData) (string, error) {
	var buf strings.Builder
	err := fce.htmlTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (fce FormCoreEngine) GenerateJavascript(data dtos.FormCoreData) (string, error) {
	var buf strings.Builder
	err := fce.jsTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FormCoreTemplate contains the full CSS template including dynamic parts
// Used for direct form rendering (backward compatibility)
const FormCoreTemplate = `
        /* 12-Column Responsive Grid System */
        .jform-grid {
            display: grid;
            grid-template-columns: repeat(1, minmax(0, 1fr));
            gap: 1.5rem;
            width: 100%;
        }
        
        /* Default 12-column grid - applies when jform-lg-grid-cols-12 class is present */
        .jform-grid.jform-lg-grid-cols-12 {
            grid-template-columns: repeat(1, minmax(0, 1fr)); /* Mobile: single column */
        }
        
        /* Column spans for mobile-first approach */
        .jform-col-1 { grid-column: span 1 / span 1; }
        .jform-col-2 { grid-column: span 1 / span 1; } /* Mobile: stack */
        .jform-col-3 { grid-column: span 1 / span 1; }
        .jform-col-4 { grid-column: span 1 / span 1; }
        .jform-col-5 { grid-column: span 1 / span 1; }
        .jform-col-6 { grid-column: span 1 / span 1; }
        .jform-col-7 { grid-column: span 1 / span 1; }
        .jform-col-8 { grid-column: span 1 / span 1; }
        .jform-col-9 { grid-column: span 1 / span 1; }
        .jform-col-10 { grid-column: span 1 / span 1; }
        .jform-col-11 { grid-column: span 1 / span 1; }
        .jform-col-12 { grid-column: span 1 / span 1; }
        
        /* Small screens (640px and up) */
        @media (min-width: 640px) {
            .jform-grid { gap: 1.5rem; }
            .jform-sm-col-1 { grid-column: span 1 / span 1; }
            .jform-sm-col-2 { grid-column: span 2 / span 2; }
            .jform-sm-col-3 { grid-column: span 3 / span 3; }
            .jform-sm-col-4 { grid-column: span 4 / span 4; }
            .jform-sm-col-5 { grid-column: span 5 / span 5; }
            .jform-sm-col-6 { grid-column: span 6 / span 6; }
            .jform-sm-col-7 { grid-column: span 7 / span 7; }
            .jform-sm-col-8 { grid-column: span 8 / span 8; }
            .jform-sm-col-9 { grid-column: span 9 / span 9; }
            .jform-sm-col-10 { grid-column: span 10 / span 10; }
            .jform-sm-col-11 { grid-column: span 11 / span 11; }
            .jform-sm-col-12 { grid-column: span 12 / span 12; }
            .jform-sm-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)); }
        }
        
        /* Medium screens (768px and up) */
        @media (min-width: 768px) {
            .jform-grid { gap: 1.5rem; }
            .jform-md-col-1 { grid-column: span 1 / span 1; }
            .jform-md-col-2 { grid-column: span 2 / span 2; }
            .jform-md-col-3 { grid-column: span 3 / span 3; }
            .jform-md-col-4 { grid-column: span 4 / span 4; }
            .jform-md-col-5 { grid-column: span 5 / span 5; }
            .jform-md-col-6 { grid-column: span 6 / span 6; }
            .jform-md-col-7 { grid-column: span 7 / span 7; }
            .jform-md-col-8 { grid-column: span 8 / span 8; }
            .jform-md-col-9 { grid-column: span 9 / span 9; }
            .jform-md-col-10 { grid-column: span 10 / span 10; }
            .jform-md-col-11 { grid-column: span 11 / span 11; }
            .jform-md-col-12 { grid-column: span 12 / span 12; }
            .jform-md-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)); }
        }
        
        /* Large screens (1024px and up) */
        @media (min-width: 1024px) {
            .jform-grid { gap: 1.5rem; }
            .jform-lg-col-1 { grid-column: span 1 / span 1; }
            .jform-lg-col-2 { grid-column: span 2 / span 2; }
            .jform-lg-col-3 { grid-column: span 3 / span 3; }
            .jform-lg-col-4 { grid-column: span 4 / span 4; }
            .jform-lg-col-5 { grid-column: span 5 / span 5; }
            .jform-lg-col-6 { grid-column: span 6 / span 6; }
            .jform-lg-col-7 { grid-column: span 7 / span 7; }
            .jform-lg-col-8 { grid-column: span 8 / span 8; }
            .jform-lg-col-9 { grid-column: span 9 / span 9; }
            .jform-lg-col-10 { grid-column: span 10 / span 10; }
            .jform-lg-col-11 { grid-column: span 11 / span 11; }
            .jform-lg-col-12 { grid-column: span 12 / span 12; }
            .jform-lg-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)) !important; }
            
            /* Apply 12-column grid when class is present */
            .jform-grid.jform-lg-grid-cols-12 { grid-template-columns: repeat(12, minmax(0, 1fr)); }
        }
        
        /* Extra large screens (1280px and up) */
        @media (min-width: 1280px) {
            .jform-xl-col-1 { grid-column: span 1 / span 1; }
            .jform-xl-col-2 { grid-column: span 2 / span 2; }
            .jform-xl-col-3 { grid-column: span 3 / span 3; }
            .jform-xl-col-4 { grid-column: span 4 / span 4; }
            .jform-xl-col-5 { grid-column: span 5 / span 5; }
            .jform-xl-col-6 { grid-column: span 6 / span 6; }
            .jform-xl-col-7 { grid-column: span 7 / span 7; }
            .jform-xl-col-8 { grid-column: span 8 / span 8; }
            .jform-xl-col-9 { grid-column: span 9 / span 9; }
            .jform-xl-col-10 { grid-column: span 10 / span 10; }
            .jform-xl-col-11 { grid-column: span 11 / span 11; }
            .jform-xl-col-12 { grid-column: span 12 / span 12; }
        }
        
        /* Utility Classes for Common Tailwind Patterns */
        .jform-hidden { display: none; }
        .jform-block { display: block; }
        .jform-flex { display: flex; }
        .jform-items-center { align-items: center; }
        .jform-gap-4 { gap: 1rem; }
        .jform-gap-6 { gap: 1.5rem; }
        .jform-space-y-8 > * + * { margin-top: 2rem; }
        .jform-space-y-4 > * + * { margin-top: 1rem; }
        
        /* Responsive visibility utilities */
        @media (min-width: 1024px) {
            .jform-lg-block { display: block; }
            .jform-hidden.jform-lg-block { display: block; }
        }
        
        /* Container and layout utilities */
        .jform-max-w-7xl { max-width: 80rem; }
        .jform-mx-auto { margin-left: auto; margin-right: auto; }
        .jform-p-6 { padding: 1.5rem; }
        .jform-sm-p-8 { padding: 1.5rem; }
        @media (min-width: 640px) {
            .jform-sm-p-8 { padding: 2rem; }
        }
        .jform-bg-white { background-color: #ffffff; }
        .jform-bg-gray-50 { background-color: #f9fafb; }
        .jform-rounded-lg { border-radius: 0.5rem; }
        .jform-h-fit { height: fit-content; }
        .jform-sticky { position: sticky; }
        .jform-top-4 { top: 1rem; }
        
        /* Text and typography utilities */
        .jform-text-2xl { font-size: 1.5rem; line-height: 2rem; }
        .jform-font-bold { font-weight: 700; }
        .jform-text-gray-900 { color: #111827; }
        .jform-border-b-2 { border-bottom-width: 2px; }
        .jform-border-gray-200 { border-color: #e5e7eb; }
        .jform-pb-2 { padding-bottom: 0.5rem; }
        .jform-text-lg { font-size: 1.125rem; line-height: 1.75rem; }
        .jform-font-semibold { font-weight: 600; }
        .jform-font-italic { font-style: italic; }
        .jform-underline { text-decoration: underline; }
        
        /* Text alignment utilities */
        .jform-text-center { text-align: center; }
        .jform-text-right { text-align: right; }
        
        /* Form element utilities */
        .jform-mb-4 { margin-bottom: 1rem; }
        .jform-mb-6 { margin-bottom: 1.5rem; }
        .jform-block { display: block; }
        .jform-text-sm { font-size: 0.875rem; line-height: 1.25rem; }
        .jform-font-medium { font-weight: 500; }
        .jform-text-gray-700 { color: #374151; }
        .jform-mb-2 { margin-bottom: 0.5rem; }
        .jform-w-full { width: 100%; }
        .jform-px-4 { padding-left: 1rem; padding-right: 1rem; }
        .jform-py-3 { padding-top: 0.75rem; padding-bottom: 0.75rem; }
        .jform-border { border-width: 1px; }
        .jform-border-gray-300 { border-color: #d1d5db; }
        .jform-rounded-lg { border-radius: 0.5rem; }
        .jform-focus-ring-2:focus { outline: none; box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-focus-ring-blue-500:focus { outline: none; box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-focus-border-blue-500:focus { border-color: #3b82f6; }
        .jform-transition-colors { transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out; }
        .jform-duration-200 { transition-duration: 0.2s; }
        .jform-mt-2 { margin-top: 0.5rem; }
        .jform-text-red-600 { color: #dc2626; }
        .jform-mt-8 { margin-top: 2rem; }
        .jform-bg-green-600 { background-color: #059669; }
        .jform-hover-bg-green-700:hover { background-color: #047857; }
        .jform-text-white { color: #ffffff; }
        .jform-font-bold { font-weight: 700; }
        .jform-py-4 { padding-top: 1rem; padding-bottom: 1rem; }
        .jform-px-8 { padding-left: 2rem; padding-right: 2rem; }
        .jform-focus-outline-none:focus { outline: none; }
        .jform-focus-ring-2:focus { box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-focus-ring-green-500:focus { box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.5); }
        .jform-focus-ring-offset-2:focus { box-shadow: 0 0 0 2px transparent, 0 0 0 4px rgba(16, 185, 129, 0.5); }
        .jform-text-lg { font-size: 1.125rem; line-height: 1.75rem; }
        .jform-mr-2 { margin-right: 0.5rem; }
        .jform-h-4 { height: 1rem; }
        .jform-w-4 { width: 1rem; }
        .jform-text-blue-600 { color: #2563eb; }
        .jform-focus-ring-blue-500:focus { box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); }
        .jform-border-gray-300 { border-color: #d1d5db; }
        .jform-rounded { border-radius: 0.25rem; }

        /* Default form styling (fallbacks) */
        .form-container {
            background: transparent;
            padding: 1rem;
            border-radius: 8px;
        }
        
        .form-field {
            margin-bottom: 1.5rem;
        }
        
        .form-field label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 500;
            color: #374151;
        }
        
        .form-field input,
        .form-field textarea {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #d1d5db;
            border-radius: 4px;
            font-size: 1rem;
            transition: border-color 0.2s;
        }
        
        .form-field input:focus,
        .form-field textarea:focus {
            outline: none;
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }
        
        .form-field textarea {
            min-height: 100px;
            resize: vertical;
        }
        
        h1, h2, h3, h4, h5, h6 {
            margin: 0 0 1rem 0;
            font-weight: 600;
            line-height: 1.2;
        }
        
        h1 { font-size: 2.5rem; }
        h2 { font-size: 2rem; }
        h3 { font-size: 1.75rem; }
        h4 { font-size: 1.5rem; }
        h5 { font-size: 1.25rem; }
        h6 { font-size: 1.125rem; }
        
        p {
            margin: 0 0 1rem 0;
            line-height: 1.6;
            color: #6b7280;
        }
        
        /* Form field element styles */
        select {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #d1d5db;
            border-radius: 4px;
            font-size: 1rem;
            transition: border-color 0.2s;
            background-color: white;
        }
        
        select:focus {
            outline: none;
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }
        
        fieldset {
            border: none;
            padding: 0;
            margin: 0;
        }
        
        legend {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 500;
            color: #374151;
            padding: 0;
        }
        
        .radio-option, .checkbox-option {
            display: flex;
            align-items: center;
            margin-bottom: 0.5rem;
        }
        
        .radio-option input, .checkbox-option input {
            width: auto;
            margin-right: 0.5rem;
        }
        
        .radio-option label, .checkbox-option label {
            margin-bottom: 0;
            cursor: pointer;
        }
        
        button[type="submit"] {
            width: 100%;
            background-color: #3b82f6;
            color: white;
            padding: 0.75rem 1.5rem;
            border: none;
            border-radius: 4px;
            font-size: 1rem;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        
        button[type="submit"]:hover {
            background-color: #2563eb;
        }
        
        .captcha-placeholder {
            border: 2px dashed #d1d5db;
            border-radius: 4px;
            padding: 1rem;
            text-align: center;
            background-color: #f9fafb;
        }
        
        .captcha-box {
            color: #6b7280;
        }
        
        .captcha-box span {
            display: block;
            margin-bottom: 0.5rem;
            font-size: 1rem;
        }
        
        .captcha-box small {
            font-size: 0.875rem;
            color: #9ca3af;
        }

        /* Form Submission Animation Styles */

        /* Button loading state */
        .jform-btn-loading {
            opacity: 0.7;
            cursor: not-allowed;
            position: relative;
            padding-right: 3rem;
        }

        /* Button disabled state */
        .jform-btn-disabled {
            opacity: 0.5;
            cursor: not-allowed !important;
            background-color: #9ca3af !important;
        }

        .jform-btn-disabled:hover {
            background-color: #9ca3af !important;
        }

        /* Spinner animation */
        .jform-spinner {
            display: inline-block;
            width: 1rem;
            height: 1rem;
            border: 2px solid rgba(255, 255, 255, 0.3);
            border-top-color: white;
            border-radius: 50%;
            animation: jform-spin 0.6s linear infinite;
            position: absolute;
            right: 1.5rem;
            top: 50%;
            margin-top: -0.5rem;
        }

        @keyframes jform-spin {
            to { transform: rotate(360deg); }
        }

        /* Form submitting state - dims and disables the form */
        .jform-submitting {
            opacity: 0.6;
            pointer-events: none;
            transition: opacity 0.3s ease;
        }

        .jform-submitting input,
        .jform-submitting textarea,
        .jform-submitting select,
        .jform-submitting button {
            cursor: not-allowed;
        }

        /* Message container styles */
        .jform-message {
            padding: 1rem;
            border-radius: 0.5rem;
            margin-bottom: 1rem;
            animation: jform-slide-down 0.3s ease-out;
            display: flex;
            align-items: flex-start;
            gap: 0.75rem;
        }

        .jform-message-success {
            background-color: #d1fae5;
            color: #065f46;
            border: 1px solid #10b981;
        }

        .jform-message-error {
            background-color: #fee2e2;
            color: #991b1b;
            border: 1px solid #ef4444;
        }

        .jform-message-info {
            background-color: #dbeafe;
            color: #1e40af;
            border: 1px solid #3b82f6;
        }

        .jform-message-icon {
            flex-shrink: 0;
            width: 1.25rem;
            height: 1.25rem;
            margin-top: 0.125rem;
        }

        .jform-message-content {
            flex-grow: 1;
        }

        /* Slide down animation for messages */
        @keyframes jform-slide-down {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        /* Fade out animation */
        .jform-fade-out {
            animation: jform-fade-out 0.3s ease-out forwards;
        }

        @keyframes jform-fade-out {
            to {
                opacity: 0;
                transform: translateY(-10px);
            }
        }

        /* Additional utility classes for animations */
        .jform-transition-all {
            transition: all 0.3s ease;
        }

        .jform-opacity-0 { opacity: 0; }
        .jform-opacity-50 { opacity: 0.5; }
        .jform-opacity-60 { opacity: 0.6; }
        .jform-opacity-100 { opacity: 1; }

        .jform-scale-95 { transform: scale(0.95); }
        .jform-scale-100 { transform: scale(1); }

        .jform-pointer-events-none { pointer-events: none; }
        .jform-pointer-events-auto { pointer-events: auto; }

        /* Hidden utility */
        .jform-hidden { display: none; }

        /* Field Validation Styles */

        /* Invalid field styling - red border and shadow */
        .jform-field-invalid {
            border-color: #ef4444 !important;
            box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1) !important;
        }

        /* Error message below field */
        .jform-field-error-message {
            color: #dc2626;
            font-size: 0.875rem;
            margin-top: 0.5rem;
            display: flex;
            align-items: flex-start;
            gap: 0.25rem;
            animation: jform-slide-down 0.2s ease-out;
        }

        .jform-field-error-icon {
            flex-shrink: 0;
            margin-top: 0.125rem;
        }

        /* Shake animation for invalid fields */
        @keyframes jform-shake {
            0%, 100% { transform: translateX(0); }
            10%, 30%, 50%, 70%, 90% { transform: translateX(-4px); }
            20%, 40%, 60%, 80% { transform: translateX(4px); }
        }

        .jform-shake {
            animation: jform-shake 0.4s ease-in-out;
        }

        /* Layout System Base Styles */
        
        /* Stacked Layout (default) */
        .form-field.jform-layout-stacked {
            display: block;
        }
        .form-field.jform-layout-stacked label {
            display: block;
            margin-bottom: 0.5rem;
        }
        
        /* Inline Layout */
        .form-field.jform-layout-inline {
            display: flex;
            align-items: flex-start;
            gap: 1rem;
        }
        .form-field.jform-layout-inline .inline-label {
            flex-shrink: 0;
            padding-top: 0.75rem;
            margin-bottom: 0;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-25 {
            width: 25%;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-30 {
            width: 30%;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-40 {
            width: 40%;
        }
        .form-field.jform-layout-inline .inline-label.jform-label-width-50 {
            width: 50%;
        }
        .form-field.jform-layout-inline .inline-input {
            flex-grow: 1;
        }
        .form-field.jform-layout-inline .inline-input input,
        .form-field.jform-layout-inline .inline-input textarea,
        .form-field.jform-layout-inline .inline-input select {
            width: 100%;
        }
        
        /* Floating Layout */
        .form-field.jform-layout-floating {
            position: relative;
        }
        .form-field.jform-layout-floating .floating-input-container {
            position: relative;
        }
        .form-field.jform-layout-floating .floating-label {
            position: absolute;
            top: 0.75rem;
            left: 0.75rem;
            background: white;
            padding: 0 0.25rem;
            transition: all 0.2s ease-in-out;
            pointer-events: none;
            color: #6b7280;
            font-size: 1rem;
            z-index: 1;
        }
        .form-field.jform-layout-floating input:focus + .floating-label,
        .form-field.jform-layout-floating input:not(:placeholder-shown) + .floating-label,
        .form-field.jform-layout-floating textarea:focus + .floating-label,
        .form-field.jform-layout-floating textarea:not(:placeholder-shown) + .floating-label {
            top: -0.5rem;
            font-size: 0.75rem;
            color: #3b82f6;
        }
        .form-field.jform-layout-floating input,
        .form-field.jform-layout-floating textarea {
            padding-top: 1.5rem;
            padding-bottom: 0.5rem;
        }
        
        /* Hidden Layout */
        .form-field.jform-layout-hidden label {
            display: none !important;
        }
        
        /* Responsive Layout Behavior - Direct Selectors (Option 1) */
        @media (max-width: 767px) {
            /* Mobile responsive layout overrides */
            .form-field.jform-mobile-layout-stacked {
                display: block;
            }
            .form-field.jform-mobile-layout-stacked label {
                display: block;
                margin-bottom: 0.5rem;
            }
            
            .form-field.jform-mobile-layout-inline {
                display: flex;
                align-items: flex-start;
                gap: 1rem;
            }
            .form-field.jform-mobile-layout-inline .inline-label {
                flex-shrink: 0;
                padding-top: 0.75rem;
                margin-bottom: 0;
            }
            .form-field.jform-mobile-layout-inline .inline-input {
                flex-grow: 1;
            }
            
            .form-field.jform-mobile-layout-floating {
                position: relative;
            }
            .form-field.jform-mobile-layout-floating .floating-input-container {
                position: relative;
            }
            .form-field.jform-mobile-layout-floating .floating-label {
                position: absolute;
                top: 0.75rem;
                left: 0.75rem;
                background: white;
                padding: 0 0.25rem;
                transition: all 0.2s ease-in-out;
                pointer-events: none;
                color: #6b7280;
                font-size: 1rem;
                z-index: 1;
            }
            
            .form-field.jform-mobile-layout-hidden label {
                display: none !important;
            }
        }
        
        @media (min-width: 768px) and (max-width: 1023px) {
            /* Tablet responsive layout overrides */
            .form-field.jform-tablet-layout-stacked {
                display: block;
            }
            .form-field.jform-tablet-layout-stacked label {
                display: block;
                margin-bottom: 0.5rem;
            }
            
            .form-field.jform-tablet-layout-inline {
                display: flex;
                align-items: flex-start;
                gap: 1rem;
            }
            .form-field.jform-tablet-layout-inline .inline-label {
                flex-shrink: 0;
                padding-top: 0.75rem;
                margin-bottom: 0;
            }
            .form-field.jform-tablet-layout-inline .inline-input {
                flex-grow: 1;
            }
            
            .form-field.jform-tablet-layout-floating {
                position: relative;
            }
            .form-field.jform-tablet-layout-floating .floating-input-container {
                position: relative;
            }
            .form-field.jform-tablet-layout-floating .floating-label {
                position: absolute;
                top: 0.75rem;
                left: 0.75rem;
                background: white;
                padding: 0 0.25rem;
                transition: all 0.2s ease-in-out;
                pointer-events: none;
                color: #6b7280;
                font-size: 1rem;
                z-index: 1;
            }
            
            .form-field.jform-tablet-layout-hidden label {
                display: none !important;
            }
        }
        
        @media (min-width: 1024px) {
            /* Desktop responsive layout overrides */
            .form-field.jform-desktop-layout-stacked {
                display: block;
            }
            .form-field.jform-desktop-layout-stacked label {
                display: block;
                margin-bottom: 0.5rem;
            }
            
            .form-field.jform-desktop-layout-inline {
                display: flex;
                align-items: flex-start;
                gap: 1rem;
            }
            .form-field.jform-desktop-layout-inline .inline-label {
                flex-shrink: 0;
                padding-top: 0.75rem;
                margin-bottom: 0;
            }
            .form-field.jform-desktop-layout-inline .inline-input {
                flex-grow: 1;
            }
            
            .form-field.jform-desktop-layout-floating {
                position: relative;
            }
            .form-field.jform-desktop-layout-floating .floating-input-container {
                position: relative;
            }
            .form-field.jform-desktop-layout-floating .floating-label {
                position: absolute;
                top: 0.75rem;
                left: 0.75rem;
                background: white;
                padding: 0 0.25rem;
                transition: all 0.2s ease-in-out;
                pointer-events: none;
                color: #6b7280;
                font-size: 1rem;
                z-index: 1;
            }
            
            .form-field.jform-desktop-layout-hidden label {
                display: none !important;
            }
        }
        
        /* Dynamic field type styling - converts Tailwind utility classes and layout system */
        {{generateFieldCSS .FormStyling}}
        
        /* Override theme styles for reliability */
        .form-container {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif !important;
            line-height: 1.6 !important;
            color: #111827 !important;
        }
        
        input, select, textarea {
            font-family: inherit !important;
            font-size: 1rem !important;
            line-height: 1.5 !important;
        }`
