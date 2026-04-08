<script>
    import { onMount } from 'svelte';
    import Results from '../Components/results.svelte'
    import Header from '../Components/header.svelte'
    import Footer from '../Components/footer.svelte'
    import LoginPopup from '../Components/login_popup.svelte';

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