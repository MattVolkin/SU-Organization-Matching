<!-- @component Provides centralized API request handling with authentication token management and debug logging. -->
<script context="module">
/**
 * @function APICreater - Centralized API request handler that manages authentication tokens, headers, and request/response logging
 * @param {string} method - HTTP method (GET, POST, PATCH, DELETE, etc.)
 * @param {string} url - API endpoint URL
 * @param {object} body - Request body payload, automatically JSON stringified for non-GET requests
 * @param {boolean} debug - Enable console logging of requests and responses for debugging, defaults to false
 * @returns {Promise<object>} Parsed JSON response from the server
 */
export async function APICreater(method, url, body, debug = false) {
    const token = (typeof localStorage !== 'undefined' ? localStorage.getItem('authToken') : '');
    const requestHeaders = {
        'Content-Type': 'application/json',
    };

    if (token) {
        requestHeaders.Authorization = `Bearer ${token}`;
    }

    const requestOptions = {
        method,
        headers: requestHeaders,
        credentials: 'include',
    };

    if (body !== undefined && body !== null && method.toUpperCase() !== 'GET') {
        requestOptions.body = JSON.stringify(body);
    }

    if (debug) {
        console.log('API Request:', {
            method,
            url,
            headers: requestHeaders,
            body: requestOptions.body ? JSON.parse(requestOptions.body) : null,
        });
    }

    const response = await fetch(url, requestOptions);
    if (debug) {
        const responseText = await response.clone().text();
            let responseData = null;

            try {
                responseData = JSON.parse(responseText);
            } catch (error) {
                responseData = null;
            }

        console.log('API Response:', {
            status: response.status,
            ok: response.ok,
            headers: Object.fromEntries(response.headers.entries()),
            body: responseText,
                data: responseData,
        })
    }

    return response.json();
}
</script>