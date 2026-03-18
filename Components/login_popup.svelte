
<script>
  import { onMount, onDestroy } from 'svelte';

  let {
    autoOpen = false,
    onSuccess = () => {},
    onBlocked = () => {},
  } = $props();

  let hasAutoOpened = $state(false);
  let popupBlocked = $state(false);

  function openLoginPopup() {
    popupBlocked = false;
    const popup = window.open(
      '/login?popup=1',
      'google-login',
      'width=520,height=640,menubar=no,toolbar=no,location=yes,resizable=yes,scrollbars=yes,status=no'
    );

    if (!popup) {
      popupBlocked = true;
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

    popupBlocked = false;
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

{#if autoOpen}
  <dialog open class="login-popup" aria-label="Sign in">
    <h2>Sign in required</h2>
    <p>Use Google sign-in to continue viewing your matches.</p>
    <button type="button" onclick={openLoginPopup}>Continue with Google</button>
    {#if popupBlocked}
      <p class="warn">Popup was blocked. Allow popups for this site, then click the button again.</p>
    {/if}
  </dialog>
{/if}

<style>
  .login-popup {
    position: fixed;
    top: 6.5rem;
    right: 1rem;
    width: min(24rem, calc(100vw - 2rem));
    padding: 1rem;
    background: #ffffff;
    border: 1px solid #d1d5db;
    border-radius: 0.75rem;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    z-index: 20;
  }

  h2 {
    margin: 0 0 0.5rem 0;
    font-size: 1.05rem;
  }

  p {
    margin: 0 0 0.75rem 0;
    line-height: 1.4;
  }

  button {
    border: none;
    border-radius: 0.5rem;
    padding: 0.55rem 0.85rem;
    font-weight: 600;
    cursor: pointer;
    background: #0f766e;
    color: #ffffff;
  }

  button:hover {
    background: #0d665f;
  }

  .warn {
    margin-top: 0.75rem;
    color: #9f1239;
    font-size: 0.9rem;
  }

  @media (max-width: 640px) {
    .login-popup {
      left: 1rem;
      right: 1rem;
      top: 5.75rem;
      width: auto;
    }
  }
</style>
