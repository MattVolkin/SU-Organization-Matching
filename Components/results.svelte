<script>
  import { onDestroy, onMount } from 'svelte';
  import AdminSwitch from './AdminSwitch.svelte'
  import Header from './header.svelte'
  import Footer from './footer.svelte'
  import LoginPopup from './login_popup.svelte'
  let adminPreviewType = $state('admin')
  let results = $state(["Computer Science Club"]) // for now this is just a placeholder, in the future this will be an array of clubs that match the user's interests and demographics, fetched from the backend server when the page loads
  let pageNum = $state(1) // for keeping track of what page of results the user is on. Not currently being used but will be helpful for future implementation 
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
      await getResults();
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

  async function getResults() {
    const response = await fetch('/results')
    const payload = await response.json().catch(() => [])
    results = Array.isArray(payload) && payload.length > 0
      ? payload
      : ["Computer Science Club"]
  }

  async function nextPage() {
    if (pageNum < results.length) {
      pageNum += 1
    }
  }

  async function prevPage() {
    if (pageNum > 1) {
      pageNum -= 1
      await getResults() // for future implementation when we have more results than we want to show on one page, this will fetch the previous page of results from the backend
    }
  }
  
  async function getClubInfo(club) {

    
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
  
<AdminSwitch enabled={true} value={adminPreviewType} onChange={(nextUserType) => (adminPreviewType = nextUserType)} />
<Header userType="admin" previewAs={adminPreviewType} />
<LoginPopup autoOpen={!isAuthChecking && !isAuthenticated} />

<main class="results-page">
  {#if isAuthChecking}
    <section class="status-card">
      <p>Checking sign-in status...</p>
    </section>
  {:else if !isAuthenticated}
    <section class="status-card">
      <p>Please complete sign-in to view organization information.</p>
    </section>
  {:else}
    <section class="result-card">
      <h1>{results[pageNum-1]}</h1>
      <h2>Hi we are {results[pageNum-1]} and we are commited to to provide a safe space to play games and hang out with other computer nerds </h2>

      <div class="pager">
        <button onclick={prevPage} disabled={pageNum === 1}>Previous</button>
        <button onclick={nextPage} disabled={pageNum >= results.length}>Next</button>
      </div>

      <h3> Activities we do include</h3>
      <ul>
        <li>Playing games</li>
        <li>Trivia nights </li>
        <li>Presentation nights</li>
        <li>Video game tournoments</li>
        <li>Scavenger hunts</li>
      </ul>
      <h3>Meeting Information </h3>
      <p> Every Thursday at 6:30 pm in FJS 310 (the CS lounge)</p>
    </section>
  {/if}
</main>

<Footer />

<style>
  .results-page {
    --page-bg: linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
    --card-bg: #ffffff;
    --card-border: #dbe7f3;
    --text-main: #10243a;
    --text-subtle: #31506e;
    --action: #1f6f8b;
    --action-hover: #195d76;

    min-height: calc(100vh - 220px);
    padding: 1rem;
    background: var(--page-bg);
    color: var(--text-main);
  }

  .status-card,
  .result-card {
    width: min(100%, 860px);
    margin: 0 auto;
    background: var(--card-bg);
    border: 1px solid var(--card-border);
    border-radius: 1rem;
    box-shadow: 0 10px 22px rgba(16, 36, 58, 0.1);
    padding: 1rem;
  }

  .status-card p {
    margin: 0;
    font-size: 1rem;
  }

  h1 {
    margin: 0;
    font-size: clamp(1.3rem, 2vw + 0.9rem, 2rem);
    line-height: 1.15;
  }

  h2 {
    margin: 0.75rem 0 1rem 0;
    font-size: clamp(1rem, 1vw + 0.85rem, 1.25rem);
    font-weight: 500;
    color: var(--text-subtle);
    line-height: 1.45;
  }

  h3 {
    margin: 1rem 0 0.5rem 0;
    font-size: 1rem;
    letter-spacing: 0.01em;
  }

  p,
  li {
    font-size: 0.98rem;
    line-height: 1.55;
  }

  ul {
    margin: 0 0 0.75rem 0;
    padding-left: 1.2rem;
  }

  .pager {
    display: flex;
    gap: 0.65rem;
    margin: 0.25rem 0 0.9rem 0;
  }

  button {
    border: none;
    border-radius: 0.55rem;
    padding: 0.55rem 0.95rem;
    font-size: 0.95rem;
    font-weight: 600;
    color: #ffffff;
    background: var(--action);
    cursor: pointer;
    transition: background-color 0.2s ease, transform 0.2s ease;
  }

  button:hover:not(:disabled) {
    background: var(--action-hover);
    transform: translateY(-1px);
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.6;
    transform: none;
  }

  @media (min-width: 900px) {
    .results-page {
      padding: 2rem;
    }

    .status-card,
    .result-card {
      padding: 1.75rem;
      border-radius: 1.1rem;
    }

    .pager {
      justify-content: flex-start;
    }
  }

  @media (max-width: 640px) {
    .results-page {
      padding: 0.85rem;
    }

    .status-card,
    .result-card {
      padding: 0.9rem;
      border-radius: 0.85rem;
    }

    .pager {
      flex-direction: column;
    }

    button {
      width: 100%;
    }
  }
</style>