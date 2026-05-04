<!-- @component Admin home page wrapper that loads user authentication and displays the AdminHome management interface. -->
<script>
/**
 * @type {state} userType - Current user's role loaded from authentication token
 * @function loadUserType - Fetches and sets user's role from backend, redirects non-admins
 * @lifecycle onMount - Loads user type on component initialization
 */
import { onMount } from 'svelte';

import AdminHome from '../AdminHome.svelte';
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
<AdminHome />

