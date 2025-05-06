/**
 * CSRF protection utilities for AJAX requests
 */

// Get CSRF token from meta tag
export function getCSRFToken() {
  const token = document.querySelector('meta[name="csrf-token"]');
  return token ? token.getAttribute('content') : '';
}

// Add CSRF token to fetch requests
export function fetchWithCSRF(url, options = {}) {
  // Default to GET if method not specified
  const method = options.method || 'GET';
  
  // Only add token for state-changing methods
  if (method !== 'GET' && method !== 'HEAD') {
    // Initialize headers if not present
    if (!options.headers) {
      options.headers = {};
    }

    // Add CSRF token to headers
    options.headers['X-CSRF-Token'] = getCSRFToken();
  }

  // Return the fetch promise
  return fetch(url, options);
}

// For forms, ensure CSRF token is added automatically to any dynamically created forms
export function addCSRFToForm(form) {
  // Check if form already has a CSRF token
  if (!form.querySelector('input[name="gorilla.csrf.Token"]')) {
    const token = getCSRFToken();
    if (token) {
      const input = document.createElement('input');
      input.type = 'hidden';
      input.name = 'gorilla.csrf.Token';
      input.value = token;
      form.appendChild(input);
    }
  }
}

// Initialize - Add CSRF token to any dynamically created forms
document.addEventListener('DOMContentLoaded', () => {
  // Add a meta tag with the CSRF token
  const csrfTokenMeta = document.createElement('meta');
  csrfTokenMeta.name = 'csrf-token';
  
  // Get the token from the page (using the first CSRF input field as source)
  const firstCSRFInput = document.querySelector('input[name="gorilla.csrf.Token"]');
  if (firstCSRFInput) {
    csrfTokenMeta.content = firstCSRFInput.value;
    document.head.appendChild(csrfTokenMeta);
  }

  // Monitor for dynamically added forms
  const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.addedNodes) {
        mutation.addedNodes.forEach(node => {
          // Check if the added node is a form or contains forms
          if (node.nodeName === 'FORM') {
            addCSRFToForm(node);
          } else if (node.querySelectorAll) {
            const forms = node.querySelectorAll('form');
            forms.forEach(form => addCSRFToForm(form));
          }
        });
      }
    }
  });

  // Start observing the document with the configured parameters
  observer.observe(document.body, {
    childList: true,
    subtree: true
  });
});