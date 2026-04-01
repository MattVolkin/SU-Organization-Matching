<script context="module">
// /api/results
export async function APICreater(method,url, body, loginToken, debug = false) {
    const response = await fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${loginToken}`
        },
        body: JSON.stringify(body)
    })
     if (debug) {
        console.log(`API Response:`, response);
    }


    return response.json()
}
</script>