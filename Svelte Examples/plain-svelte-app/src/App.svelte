<script>
  import { onDestroy, onMount } from 'svelte';
  import AdminSwitch from '../../../Components/AdminHome.svelte'
  import AdminHome from '../../../Components/AdminHome.svelte';
  import SwipingApp from '../../../Components/SwipingApp.svelte';
  import SettingsPage from '../../../Components/SettingPage.svelte';
  import LoginPopup from '../../../Components/login_popup.svelte';
  import ResultsPage from '../../../Components/results.svelte';
  import Multi_Choice from '../../../Components/Mult_choice_demo.svelte';
    import MultChoiceDemo from '../../../Components/Mult_choice_demo.svelte';
  let isAuthChecking = $state(true)
  let isAuthenticated = $state(false)
   async function promptLoginIfNeeded() {
    const tokenFromStorage = localStorage.getItem('authToken') || '';
    const headers = tokenFromStorage
      ? { Authorization: `Bearer ${tokenFromStorage}` }
      : {};

    const response = await fetch('/api/user', {
      method: 'GET',
      credentials: 'include',
      headers,
    });

    if (response.ok) {
      isAuthenticated = true;
      isAuthChecking = false;
      return;
    }

    isAuthenticated = false;
    isAuthChecking = false;
  }

  function handleAuthLogout() {
    isAuthenticated = false;
    isAuthChecking = false;
  }

  function handleAuthLogin() {
    promptLoginIfNeeded();
  }
    onMount(() => {
    window.addEventListener('auth-login', handleAuthLogin);
    window.addEventListener('auth-logout', handleAuthLogout);
    promptLoginIfNeeded();
  });

  onDestroy(() => {
    window.removeEventListener('auth-login', handleAuthLogin);
    window.removeEventListener('auth-logout', handleAuthLogout);
  });
</script>
 <!-- <AdminHome />  -->
  <LoginPopup/>
 <SettingsPage/>
 <div class="parent">

  
</div>