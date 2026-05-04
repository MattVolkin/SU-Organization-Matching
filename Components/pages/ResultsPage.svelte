<!-- @component Results page wrapper that loads user authentication and displays the organization match results interface. -->
<script>
/**
 * @type {state} userType - Current user's role loaded from authentication token
 * @function loadUserType - Fetches and sets user's role from backend
 * @lifecycle onMount - Loads user type on component initialization
 */
    // Auth + role bootstrap for role-aware header/results views.
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

<!-- Auth UI entry point -->
<LoginPopup/>

<!-- Top navigation -->
<Header userType={userType} />

<!-- Main results content -->
<Results userType={userType} />

<!-- Global footer -->
<Footer />