<script context="module">
// /api/results
export async function APICreator(method, url, body, loginToken, debug = false) {
    const token = loginToken || (typeof localStorage !== 'undefined' ? localStorage.getItem('authToken') : '');
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

    const response = await fetch(url, requestOptions);
    if (debug) {
        console.log(`API Response:`, response);
    }

    if (response.status === 204) {
        return null;
    }

    const responseText = await response.text();
    if (!responseText) {
        return null;
    }

    try {
        return JSON.parse(responseText);
    } catch {
        return responseText;
    }
}

// Backward-compatible alias for existing imports.
export const APICreater = APICreator;
</script>