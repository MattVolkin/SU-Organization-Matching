<!-- @component
  Creates a header for pages with a login button and navigation links.
  The links shown depend on the user type (admin, officer, or regular user).

  Notes:
  - Defaults to normal user view if no user type is provided.

  **TODO**:
  - Links are placeholders and should be updated once real pages are created.
-->
<script>
  import { onMount, onDestroy } from 'svelte';
/**
 * @type {props} userType - defaults to user view if no user type is provided, can be 'admin', 'officer' or 'user'
 * @type {state} isMenuOpen - boolean to track whether the mobile hamburger menu is open or closed
 * @function toggleMenu - toggles the state of isMenuOpen when the hamburger menu button is clicked
 * @function closeMenu - sets isMenuOpen to false, used to close the mobile menu when a navigation link is clicked
 */
  let { userType = "user" } = $props();
  let isMenuOpen = $state(false);
  let userEmail = $state('');
  let authToken = $state('');

  async function refreshUser() {
    const tokenFromStorage = localStorage.getItem('authToken') || '';
    authToken = tokenFromStorage;
    const headers = tokenFromStorage
      ? { Authorization: `Bearer ${tokenFromStorage}` }
      : {};

    const res = await fetch('/api/user', {
      method: 'GET',
      credentials: 'include',
      headers,
    });

    if (!res.ok) {
      userEmail = '';
      authToken = '';
      localStorage.removeItem('authToken');
      return;
    }

    const data = await res.json();
    userEmail = data.email || '';
  }

  function loginWithGooglePopup() {
    const popup = window.open(
      '/login?popup=1',
      'google-login',
      'width=520,height=640,menubar=no,toolbar=no,location=yes,resizable=yes,scrollbars=yes,status=no'
    );
    if (!popup) {
      alert('Popup blocked. Please allow popups for this site and try again.');
    }
  }

  async function logout() {
    const headers = authToken
      ? { Authorization: `Bearer ${authToken}` }
      : {};
    await fetch('/logout', {
      method: 'POST',
      credentials: 'include',
      headers,
    });
    userEmail = '';
    authToken = '';
    localStorage.removeItem('authToken');
  }

  function onAuthMessage(event) {
    if (event.origin !== window.location.origin) {
      return;
    }
    if (event.data?.type !== 'google-auth-success') {
      return;
    }

    if (event.data.token) {
      authToken = event.data.token;
      localStorage.setItem('authToken', event.data.token);
    }
    userEmail = event.data.email || '';
  }

  onMount(() => {
    window.addEventListener('message', onAuthMessage);
    refreshUser();
  });

  onDestroy(() => {
    window.removeEventListener('message', onAuthMessage);
  });

  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
  }

  function closeMenu() {
    isMenuOpen = false;
  }
</script>

<header class="header">
  <div class="header-content">
    <h1>SU Organization Matching Tool</h1>
    {#if userEmail}
      <div class="user-info">
        <span class="user-email" title={userEmail}>{userEmail}</span>
        <button class="login-button logout-btn" type="button" onclick={logout}>Logout</button>
      </div>
    {:else}
      <button class="login-button" type="button" onclick={loginWithGooglePopup}>Login</button>
    {/if}
  </div>
  <!-- Mobile hamburger menu toggle button -->
  <button
    class="menu-toggle"
    type="button"
    aria-expanded={isMenuOpen}
    aria-controls="primary-nav"
    aria-label="Toggle navigation menu"
    onclick={toggleMenu}
  >
    <span class="menu-icon" aria-hidden="true"></span>
    Menu
  </button>
  <!-- TODO: replace hrefs with actual links to pages once they are created -->
  <nav id="primary-nav" class={`nav ${isMenuOpen ? 'open' : ''}`}>
    <a href="/" onclick={closeMenu}>Home</a>
    <a href="/about" onclick={closeMenu}>About This Project</a>
    <a href="/howto" onclick={closeMenu}>How To Use This Tool</a>
    <!-- Show admin-only and officer-only links based on user type -->
    {#if userType === 'admin'}
      <a href="/create" onclick={closeMenu}>Create New Club</a>
      <a href="/change-user-type" onclick={closeMenu}>Change User Type</a>
    {:else if userType === 'officer'}
      <a href="/manage-club" onclick={closeMenu}>Manage Club</a>
    {/if}
  </nav>
</header>

<style>


  .header {
    background-color: #2c3e50;
    color: white;
    position: relative;
    padding: 1rem 2rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    text-align: center;
  }

  .header-content {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    margin-bottom: 0;
    min-height: 2.75rem;
  }

  h1 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: 600;
    text-align: center;
    width: 100%;
  }

  .user-info {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    z-index: 2;
  }

  .user-email {
    font-size: 0.8rem;
    max-width: 30vw;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    opacity: 0.95;
    min-width: 0;
  }

  .login-button {
    width: auto;
    padding: 0.35rem 0.75rem;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 0.85rem;
    cursor: pointer;
    transition: background-color 0.3s ease;
    touch-action: manipulation;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .login-button:hover {
    background-color: #2980b9;
  }

  .logout-btn {
    position: static;
  }

  .menu-toggle {
    display: none;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.65rem 1rem;
    border: 1px solid rgba(255, 255, 255, 0.3);
    background: transparent;
    color: white;
    border-radius: 4px;
    margin-bottom: 0.75rem;
    cursor: pointer;
  }

  .menu-icon {
    width: 1rem;
    height: 2px;
    background-color: currentColor;
    position: relative;
    display: inline-block;
  }

  .menu-icon::before,
  .menu-icon::after {
    content: "";
    position: absolute;
    left: 0;
    width: 1rem;
    height: 2px;
    background-color: currentColor;
  }
/* adjust the position of the lines in the hamburger menu icon for better spacing and visual balance on mobile devices */
  .menu-icon::before {
    top: -0.35rem;
  }
/*` Adjust the position of the bottom line of the hamburger menu icon to ensure proper spacing between the lines and maintain a visually balanced appearance on mobile devices. */
  .menu-icon::after {
    top: 0.35rem;
  }

  .nav {
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    align-items: center;
    width: 100%;
    gap: 0;
  }
  /* Style the navigation links to be evenly spaced and visually distinct, with a hover effect for better user experience */
  .nav a {
    color: white;
    text-decoration: none;
    font-size: 1rem;
    transition: color 0.3s ease;
    flex: 1;
    text-align: center;
    padding: 0.5rem;
    border: none;
  }
  /* Add a subtle hover effect to nav links to let the user know they are clickable */
  .nav a:hover {
    color: #3498db;
  }

  /* Desktop scaling for login button */
  @media (min-width: 768px) {
    .header-content {
      margin-bottom: 1rem;
    }

    .user-info {
      top: clamp(0.5rem, 0.35rem + 0.35vw, 0.9rem);
      right: clamp(0.5rem, 0.25rem + 0.7vw, 1.25rem);
      gap: clamp(0.5rem, 0.3rem + 0.5vw, 1rem);
    }

    .login-button {
      font-size: clamp(0.85rem, 0.72rem + 0.35vw, 1.1rem);
      padding: clamp(0.35rem, 0.25rem + 0.2vw, 0.55rem)
        clamp(0.75rem, 0.55rem + 0.5vw, 1.3rem);
      border-radius: clamp(4px, 3px + 0.2vw, 8px);
    }

    .user-email {
      font-size: clamp(0.75rem, 0.68rem + 0.2vw, 0.92rem);
      max-width: clamp(20vw, 18rem, 45vw);
    }
  }

  /* Mobile-only hamburger nav */
  @media (max-width: 767px) {
    .header {
      position: relative;
      padding: 3rem 1rem 1rem;
    }

    .header-content {
      flex-direction: column;
      gap: 1rem;
      min-height: 0;
    }

    h1 {
      font-size: 1.25rem;
    }

    .user-info {
      top: 0.4rem;
      right: 0.4rem;
      flex-direction: column;
      align-items: flex-end;
      gap: 0.35rem;
    }

    .user-email {
      font-size: 0.65rem;
      max-width: 35vw;
    }

    .login-button {
      font-size: 0.75rem;
      padding: 0.3rem 0.6rem;
    }

    .menu-toggle {
      position: absolute;
      top: 0.5rem;
      left: 0.5rem;
      width: auto;
      display: inline-flex;
      padding: 0.4rem 0.6rem;
      font-size: 0.85rem;
      gap: 0.35rem;
      margin-bottom: 0;
      z-index: 2;
    }

    .menu-icon,
    .menu-icon::before,
    .menu-icon::after {
      width: 0.85rem;
    }

    .nav {
      display: none;
      flex-direction: column;
      margin-top: 0.5rem;
    }

    .nav.open {
      display: flex;
    }

    .nav a {
      transition: background-color 0.3s ease;
      padding: 0.75rem;
      border-top: 1px solid rgba(255, 255, 255, 0.1);
    }

    .nav a:hover {
      color: white;
      background-color: rgba(52, 152, 219, 0.2);
    }
  }
</style>