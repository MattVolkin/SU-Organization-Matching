<!-- @component Creates a login popup that appears when the user needs to sign in. -->
<script>
/**
 * @type {props} autoOpen - boolean to control whether the login popup should automatically open when the component is rendered, defaults to false
 * @type {props} onSuccess - callback function that is called when the user successfully logs in, receives an object with the user's email and token as an argument
 * @type {props} onBlocked - callback function that is called when the login popup is blocked by the browser, defaults to an empty function
 * @type {state} popupBlocked - boolean to track whether the login popup was blocked by the browser, used to conditionally render a warning message in the UI
 * @function openLoginPopup - function that opens the login popup window and sets popupBlocked to true if the popup was blocked by the browser
 * @function onAuthMessage - event handler function that listens for messages from the login popup window, checks if the message is a successful authentication message, and if so, stores the token in localStorage, calls the onSuccess callback with the user's email and token, and dispatches a custom 'auth-login' event on the window with the user's email and token as the detail
 * @effect - an effect that resets popupBlocked to false whenever autoOpen changes
 * @lifecycle onMount - adds the onAuthMessage event listener to the window when the component is mounted
 * @lifecycle onDestroy - removes the onAuthMessage event listener from the window when the component is destroyed
 */
  import { onMount, onDestroy } from 'svelte';

  let {
    autoOpen = false,
    onSuccess = () => {},
    onBlocked = () => {},
  } = $props();

  let popupBlocked = $state(false);

  function setBodyLock(locked) {
    if (typeof document === 'undefined') {
      return;
    }

    document.body.classList.toggle('login-required-open', locked);
  }

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

    window.dispatchEvent(new CustomEvent('auth-login', {
      detail: {
        email: event.data.email || '',
        token: event.data.token || '',
      },
    }));

    popupBlocked = false;
  }

  $effect(() => {
    if (!autoOpen) {
      popupBlocked = false;
    }

    setBodyLock(autoOpen);
  });

  onMount(() => {
    window.addEventListener('message', onAuthMessage);
  });

  onDestroy(() => {
    window.removeEventListener('message', onAuthMessage);
    setBodyLock(false);
  });
</script>

{#if autoOpen}
  <div class="login-overlay" role="presentation" aria-hidden="false">
    <div class="login-popup" role="dialog" aria-modal="true" aria-label="Sign in">
      <h2>Sign in required</h2>
      <p>Use Google sign-in to continue viewing your matches.</p>
      <button type="button" onclick={openLoginPopup}>Continue with Google</button>
      {#if popupBlocked}
        <p class="warn">Popup was blocked. Allow popups for this site, then click the button again.</p>
      {/if}
    </div>
  </div>
{/if}

<style>
  :global(body.login-required-open) {
    overflow: hidden;
  }

  .login-overlay {
    position: fixed;
    inset: 0;
    display: grid;
    place-items: center;
    padding: 1rem;
    background: rgba(7, 16, 27, 0.6);
    backdrop-filter: blur(2px);
    z-index: 2000;
  }

  .login-popup {
    width: min(24rem, calc(100vw - 2rem));
    padding: 1rem;
    background: #ffffff;
    border: 1px solid #d1d5db;
    border-radius: 0.75rem;
    box-shadow: 0 16px 44px rgba(0, 0, 0, 0.28);
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
      width: auto;
    }
  }
</style>
