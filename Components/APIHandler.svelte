<script context="module">
// /api/results
export async function APICreater(method, url, body, loginToken, debug = false) {
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
    };

    if (body !== undefined && body !== null && method.toUpperCase() !== 'GET') {
        requestOptions.body = JSON.stringify(body);
    }

    const response = await fetch(url, requestOptions);
    if (debug) {
        console.log(`API Response:`, response);
    }

    return response.json();
}
</script>