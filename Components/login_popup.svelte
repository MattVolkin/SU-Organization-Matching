<script>
  import { onMount, onDestroy } from 'svelte';

  let {
    autoOpen = false,
    onSuccess = () => {},
    onBlocked = () => {},
  } = $props();

  let hasAutoOpened = $state(false);

  function openLoginPopup() {
    const popup = window.open(
      '/login?popup=1',
      'google-login',
      'width=520,height=640,menubar=no,toolbar=no,location=yes,resizable=yes,scrollbars=yes,status=no'
    );

    if (!popup) {
      onBlocked();
    }
  }

  function onAuthMessage(event) {
    if (event.origin !== window.location.origin) {
      return;
    }
    if (event.data?.type !== 'google-auth-success') {
      return;
    }

    if (event.data.token) {
      localStorage.setItem('authToken', event.data.token);
    }

    onSuccess({
      email: event.data.email || '',
      token: event.data.token || '',
    });
  }

  $effect(() => {
    if (autoOpen && !hasAutoOpened) {
      hasAutoOpened = true;
      openLoginPopup();
    }

    if (!autoOpen) {
      hasAutoOpened = false;
    }
  });

  onMount(() => {
    window.addEventListener('message', onAuthMessage);
  });

  onDestroy(() => {
    window.removeEventListener('message', onAuthMessage);
  });
</script>
