<script>
    import { onMount } from 'svelte';
    import Results from '../results.svelte'
    import Header from '../header.svelte'
    import Footer from '../footer.svelte'
    import LoginPopup from '../login_popup.svelte';

    let userType = 'user';

    async function loadUserType() {
        const tokenFromStorage = localStorage.getItem('authToken') || '';
        const headers = tokenFromStorage
            ? { Authorization: `Bearer ${tokenFromStorage}` }
            : {};

        const response = await fetch('/api/user', {
            method: 'GET',
            credentials: 'include',
            headers,
        });

        if (!response.ok) {
            userType = 'user';
            return;
        }

        const data = await response.json().catch(() => ({}));
        const role = String(data?.role || '').toLowerCase();
        userType = role === 'admin' || role === 'officer' ? role : 'user';
    }

    onMount(() => {
        loadUserType();
    });
</script>
<LoginPopup/>
<Header userType={userType} />
<Results userType={userType} />
<Footer />