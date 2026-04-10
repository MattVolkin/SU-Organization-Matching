<!-- @component this creates the result page template 
  **TODO**:
  - Replace placeholder data with actual API calls 
  -->
<script>
  /**
   * @type {state} results - an array of club names that match the user's interests and demographics, fetched from the backend server when the page loads, currently initialized with a placeholder value
   * @type {state} pageNum - keeps track of what page of results the user is on, initialized to 1, not currently being used but will be helpful for future implementation when we have more results than we want to show on one page
   * @type {state} isAuthChecking - boolean to track whether we are currently checking the user's authentication status, initialized to true, used to conditionally render the UI while we check if the user is signed in
   * @type {state} isAuthenticated - boolean to track whether the user is authenticated or not, initialized to false, used to conditionally render the UI based on whether the user is signed in or not
   * @function promptLoginIfNeeded - an async function that checks if the user is authenticated by making a request to the backend server with the token from localStorage, updates isAuthenticated and isAuthChecking based on the response, and calls getResults if the user is authenticated
   * @function handleAuthLogout - an event handler function that sets isAuthenticated to false and isAuthChecking to false when the user logs out, used to update the UI accordingly
   * @function handleAuthLogin - an event handler function that calls promptLoginIfNeeded when the user logs in, used to check the user's authentication status and update the UI accordingly
   * @function getResults - an async function that fetches the club results from the backend server and updates the results state, currently uses a placeholder implementation that sets results to a static array if the fetch fails or returns an invalid response
   * @function nextPage - an async function that increments pageNum to show the next page of results, currently not fully implemented as it does not fetch new results from the backend but will be helpful for future implementation when we have more results than we want to show on one page
   * @function prevPage - an async function that decrements pageNum to show the previous page of results, currently not fully implemented as it does not fetch new results from the backend but will be helpful for future implementation when we have more results than we want to show on one page
   * @function getClubInfo - a placeholder function that can be used to fetch and display more detailed information about a specific club when the user clicks on a club in the results, currently does not have an implementation but can be expanded in the future to show a modal or navigate to a club details page with more information about the club
   * @lifecycle onMount - adds event listeners for 'auth-login' and 'auth-logout' events to handle changes in authentication status, and calls promptLoginIfNeeded to check the user's authentication status when the component is mounted
   * @lifecycle onDestroy - removes the event listeners for 'auth-login' and 'auth-logout' events when the component is destroyed to prevent memory leaks
    */    
  import { onDestroy, onMount } from 'svelte';
  import Header from './header.svelte'
  import Footer from './footer.svelte'
  import LoginPopup from './login_popup.svelte'
  import { APICreater } from './APIHandler.svelte';

  const selectedImageStorageKey = 'club-selected-images';
  const legacySelectedImageStorageKey = 'club-selected-images-legacy';
  const defaultResultImage = '';

  let results = $state([]) // stores matched clubs returned by the results API
  let pageNum = $state(1) // for keeping track of what page of results the user is on.
  let isAuthChecking = $state(true)
  let isAuthenticated = $state(false)
  let selectedImageByClub = $state({})
  let ClubsPerPage = 5
  let threshold = 50

  function normalizeClubKey(value) {
    return String(value || '').trim().toLowerCase()
  }

  function normalizeClubScore(value) {
    return typeof value === 'number' ? value : Number(value || 0)
  }

  function getClubName(club) {
    return typeof club === 'string' ? club : (club?.clubName || club?.ClubName || 'Unknown Club')
  }

  function getClubScore(club) {
    return normalizeClubScore(club?.matchPercentage ?? club?.score ?? 0)
  }

  function getFilteredResults() {
    return results.filter((club) => getClubScore(club) >= threshold)
  }

  function getTotalPages() {
    const filteredResults = getFilteredResults()
    return filteredResults.length > 0 ? Math.ceil(filteredResults.length / ClubsPerPage) : 0
  }

  function getVisibleResults() {
    const filteredResults = getFilteredResults()
    const startIndex = (pageNum - 1) * ClubsPerPage
    return filteredResults.slice(startIndex, startIndex + ClubsPerPage)
  }

  function getSelectedImageForClub(club) {
    const currentClubName = getClubName(club)
    if (!currentClubName) {
      return defaultResultImage
    }

    if (selectedImageByClub[currentClubName]) {
      return selectedImageByClub[currentClubName]
    }

    const normalizedCurrent = normalizeClubKey(currentClubName)
    const matchedEntry = Object.entries(selectedImageByClub).find(([clubName]) => normalizeClubKey(clubName) === normalizedCurrent)
    if (matchedEntry) {
      return matchedEntry[1]
    }

    return defaultResultImage
  }

  function loadSelectedClubImages() {
    try {
      const savedSelection = localStorage.getItem(selectedImageStorageKey)
      const legacySelection = localStorage.getItem(legacySelectedImageStorageKey)
      selectedImageByClub = savedSelection
        ? JSON.parse(savedSelection)
        : (legacySelection ? JSON.parse(legacySelection) : {})
    } catch (error) {
      console.error('Unable to load selected club image', error)
      selectedImageByClub = {}
    }
  }

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

    isAuthenticated = false
    isAuthChecking = false
  }

  function handleAuthLogout() {
    isAuthenticated = false
    isAuthChecking = false
  }

  function handleAuthLogin() {
    promptLoginIfNeeded()
  }

  async function getResults() {
    const payload = await APICreater('GET', '/api/results?includeScores=1', null)
    results = Array.isArray(payload) && payload.length > 0
      ? payload.map((item) => ({
          id: typeof item?.id === 'number' ? item.id : 0,
          clubName: getClubName(item),
          matchPercentage: getClubScore(item),
        }))
      : []
    pageNum = 1
  }

  async function nextPage() {
    if (pageNum < getTotalPages()) {
      pageNum += 1
    }
  }

  async function prevPage() {
    if (pageNum > 1) {
      pageNum -= 1
    }
  }
  
  async function getClubInfo(club) {

    
  }

  onMount(() => {
    window.addEventListener('auth-login', handleAuthLogin)
    window.addEventListener('auth-logout', handleAuthLogout)
    loadSelectedClubImages();
    promptLoginIfNeeded();
  });

  onDestroy(() => {
    window.removeEventListener('auth-login', handleAuthLogin)
    window.removeEventListener('auth-logout', handleAuthLogout)
  });
    
</script>

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
      {#if getVisibleResults().length === 0}
        <p>No clubs meet the current threshold of {threshold}%.</p>
      {:else}
        {#each getVisibleResults() as club, index (club.id || `${club.clubName}-${index}`)}
          <article class="club-card">
            <h1>{club.clubName}</h1>

            {#if getSelectedImageForClub(club)}
              <img
                class="club-hero-image"
                src={getSelectedImageForClub(club)}
                alt={`Selected club image for ${club.clubName}`}
              />
            {/if}

            <p class="match-score">Match score: {club.matchPercentage.toFixed(0)}%</p>
            <h2>Hi, we are {club.clubName}! We are commited to to provide a safe space to play games and hang out with other computer nerds. </h2>

            <h3> Activities we do include:</h3>
            <ul>
              <li>Playing games</li>
              <li>Trivia nights </li>
              <li>Presentation nights</li>
              <li>Video game tournaments</li>
              <li>Scavenger hunts</li>
            </ul>
            <h3>Meeting Information: </h3>
            <p> Every Thursday at 6:30 pm in FJS 310</p>
          </article>
        {/each}

        <div class="pager">
          <button onclick={prevPage} disabled={pageNum === 1}>Previous</button>
          <button onclick={nextPage} disabled={pageNum >= getTotalPages()}>Next</button>
        </div>
      {/if}
    </section>
  {/if}
</main>


<style>
  :global(html),
  :global(body) {
    margin: 0;
    padding: 0;
    min-height: 100%;
    background: #eef6ff;
  }

  .results-page {
    --page-bg: linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
    --card-bg: #ffffff;
    --card-border: #dbe7f3;
    --text-main: #10243a;
    --text-subtle: #31506e;
    --action: #1f6f8b;
    --action-hover: #195d76;
    --focus-ring: #60a5fa;

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
    border: none;
    border-radius: 1rem;
    box-shadow: 0 10px 22px rgba(16, 36, 58, 0.1);
    padding: 1rem;
  }

  .status-card p {
    margin: 0;
    font-size: 1rem;
  }

  .result-card {
    display: grid;
    gap: 1rem;
  }

  .club-card {
    display: grid;
    gap: 0.5rem;
    padding: 0 0 1rem 0;
    border-bottom: 1px solid var(--card-border);
  }

  .club-card:last-of-type {
    border-bottom: none;
    padding-bottom: 0;
  }

  .match-score {
    margin: 0;
    font-weight: 700;
    color: var(--action);
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

  .club-hero-image {
    width: 100%;
    max-height: 280px;
    object-fit: cover;
    border-radius: 0.8rem;
    margin: 0.9rem 0;
    border: none;
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

  li::marker {
    color: var(--action);
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

  button:focus-visible {
    outline: 2px solid var(--focus-ring);
    outline-offset: 2px;
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

  @media (prefers-color-scheme: dark) {
    .results-page {
      --page-bg: linear-gradient(180deg, #0b1220 0%, #111827 100%);
      --card-bg: #0f172a;
      --card-border: #273449;
      --text-main: #e5edf8;
      --text-subtle: #b6c7df;
      --action: #2b8fb5;
      --action-hover: #3aa3cb;
      --focus-ring: #93c5fd;
    }

    .status-card,
    .result-card {
      box-shadow: 0 14px 28px rgba(0, 0, 0, 0.45);
    }

    button:disabled {
      opacity: 0.5;
      background: #355268;
      color: #cbd5e1;
    }
  }
</style>