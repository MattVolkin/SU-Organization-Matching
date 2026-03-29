<script>
import AdminSwitch from './AdminSwitch.svelte'
import Header from './header.svelte'
import Footer from './footer.svelte'
async function APIPostCreater(url, body, loginToken) {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${loginToken}`
        },
        body: JSON.stringify(body)
    })

    return response.json()
}
    

</script>
