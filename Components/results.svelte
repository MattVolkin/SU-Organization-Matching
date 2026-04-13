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

  const defaultResultImage = "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 640 360'><defs><linearGradient id='sky' x1='0' y1='0' x2='0' y2='1'><stop offset='0%' stop-color='%23dbeafe'/><stop offset='100%' stop-color='%23f8fafc'/></linearGradient></defs><rect width='640' height='360' fill='url(%23sky)'/><rect y='235' width='640' height='125' fill='%23d1d5db'/><path d='M0 250L120 190L210 235L330 165L440 235L520 205L640 250V360H0Z' fill='%239ca3af'/><circle cx='110' cy='84' r='34' fill='%23fde68a'/><text x='320' y='330' text-anchor='middle' font-family='Arial, sans-serif' font-size='22' fill='%234b5563'>Organization Image</text></svg>";
  const orgPhotoLibrary = new Map([
    ['computerscienceclub', new URL('./OrgPhotos/CSClub.jpg', import.meta.url).href],
    ['piratesforpride', new URL('./OrgPhotos/P4P.jpg', import.meta.url).href],
    ['reproductivejusticealliance', new URL('./OrgPhotos/ReproductiveJusticeAlliance.png', import.meta.url).href],
    ['sutertulias', new URL('./OrgPhotos/SUTertulias.jpeg', import.meta.url).href],
    ['thegame', new URL('./OrgPhotos/THEGAME.jpeg', import.meta.url).href],
    ['unitedmethodiststudentfellowship', new URL('./OrgPhotos/UnitedMethodistStudentFellowship.jpeg', import.meta.url).href],
  ])

  let results = $state([]) // stores matched clubs returned by the results API
  let pageNum = $state(1) // for keeping track of what page of results the user is on.
  let isAuthChecking = $state(true)
  let isAuthenticated = $state(false)
  let resultsPageElement = $state(null)
  let ClubsPerPage = 5
  let threshold = 30

  function scrollToResultsTop() {
    if (resultsPageElement) {
      resultsPageElement.scrollIntoView({ behavior: 'smooth', block: 'start' })
      return
    }

    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  function normalizeClubScore(value) {
    const numeric = typeof value === 'number' ? value : Number(value || 0)
    if (!Number.isFinite(numeric)) {
      return 0
    }
    if (numeric > 0 && numeric <= 1) {
      return numeric * 100
    }
    return numeric
  }

  function getClubName(club) {
    return typeof club === 'string' ? club : (club?.clubName || club?.ClubName || 'Unknown Club')
  }

  function getClubID(club) {
    if (!club || typeof club !== 'object') {
      return 0
    }

    const rawId = club?.id ?? club?.ID
    const numericId = typeof rawId === 'number' ? rawId : Number(rawId)
    return Number.isFinite(numericId) ? numericId : 0
  }

  function getTrimmedText(value) {
    return typeof value === 'string' ? value.trim() : ''
  }

  function getClubDescription(club, fallbackClub = null) {
    if (typeof club === 'string') {
      return typeof fallbackClub === 'string' ? '' : getClubDescription(fallbackClub)
    }

    const primaryDescription = getTrimmedText(
      club?.description ?? club?.Description ?? club?.clubDescription ?? club?.ClubDescription
    )
    if (primaryDescription) {
      return primaryDescription
    }

    if (!fallbackClub || typeof fallbackClub === 'string') {
      return ''
    }

    return getTrimmedText(
      fallbackClub?.description ?? fallbackClub?.Description ?? fallbackClub?.clubDescription ?? fallbackClub?.ClubDescription
    )
  }

  function getClubScore(club) {
    return normalizeClubScore(club?.matchPercentage ?? club?.score ?? 0)
  }

  function normalizeClubKey(value) {
    return String(value || '').toLowerCase().replace(/[^a-z0-9]/g, '')
  }

  function isOfficerPlaceholder(value) {
    const text = typeof value === 'string' ? value.trim().toLowerCase() : ''
    if (!text) {
      return false
    }

    return text.includes('updated by officer')
  }

  function getClubImagePath(club) {
    const configuredPath = typeof club?.imagePath === 'string' ? club.imagePath.trim() : ''
    if (configuredPath) {
      return configuredPath
    }

    const normalizedName = normalizeClubKey(getClubName(club))
    return orgPhotoLibrary.get(normalizedName) || ''
  }

  function getVisibleContactInfo(value) {
    const text = typeof value === 'string' ? value.trim() : ''
    if (!text) {
      return ''
    }

    if (isOfficerPlaceholder(text)) {
      return ''
    }

    return text
  }

  function getVisibleExternalLink(value) {
    const text = typeof value === 'string' ? value.trim() : ''
    if (!text || isOfficerPlaceholder(text)) {
      return ''
    }

    return text
  }

  function formatUpdatedAt(updatedAt) {
    if (!updatedAt) {
      return ''
    }

    const parsedDate = new Date(updatedAt)
    if (Number.isNaN(parsedDate.getTime())) {
      return ''
    }

    return parsedDate.toLocaleString()
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
    const [detailPayload, scorePayload] = await Promise.all([
      APICreater('GET', '/api/results', null),
      APICreater('GET', '/api/results?includeScores=true', null),
    ])

    const detailList = Array.isArray(detailPayload) ? detailPayload : []
    const scoreList = Array.isArray(scorePayload) ? scorePayload : []

    const detailById = new Map()
    const detailByName = new Map()
    for (const item of detailList) {
      const clubName = getClubName(item)
      const clubId = getClubID(item)

      if (clubId > 0) {
        detailById.set(clubId, item)
      }
      detailByName.set(normalizeClubKey(clubName), item)
    }

    const sourceList = scoreList.length > 0 ? scoreList : detailList

    results = sourceList.map((item) => {
      const clubId = getClubID(item)
      const clubName = getClubName(item)
      const matchedDetail = clubId > 0
        ? (detailById.get(clubId) || detailByName.get(normalizeClubKey(clubName)))
        : detailByName.get(normalizeClubKey(clubName))

      const mergedClub = {
        ...(item && typeof item === 'object' ? item : {}),
        ...(matchedDetail && typeof matchedDetail === 'object' ? matchedDetail : {}),
      }
      const resolvedId = getClubID(mergedClub) || clubId

      return {
        id: resolvedId,
        clubName,
        matchPercentage: getClubScore(item),
        description: getClubDescription(mergedClub, item),
        meetingTime: getVisibleContactInfo(mergedClub?.meetingTime ?? mergedClub?.MeetingTime),
        imagePath: getClubImagePath(mergedClub),
        externalLink: getVisibleExternalLink(mergedClub?.externalLink ?? mergedClub?.ExternalLink),
        contactInfo: getVisibleContactInfo(mergedClub?.contactInfo ?? mergedClub?.ContactInfo),
        includeOfficerEmails: Boolean(mergedClub?.includeOfficerEmails),
        updatedAt: typeof (mergedClub?.updatedAt ?? mergedClub?.UpdatedAt) === 'string'
          ? (mergedClub?.updatedAt ?? mergedClub?.UpdatedAt).trim()
          : '',
      }
    })

    pageNum = 1
  }

  async function nextPage() {
    if (pageNum < getTotalPages()) {
      pageNum += 1
      scrollToResultsTop()
    }
  }

  async function prevPage() {
    if (pageNum > 1) {
      pageNum -= 1
      scrollToResultsTop()
    }
  }

  onMount(() => {
    window.addEventListener('auth-login', handleAuthLogin)
    window.addEventListener('auth-logout', handleAuthLogout)
    promptLoginIfNeeded();
  });

  onDestroy(() => {
    window.removeEventListener('auth-login', handleAuthLogin)
    window.removeEventListener('auth-logout', handleAuthLogout)
  });
    
</script>

<LoginPopup autoOpen={!isAuthChecking && !isAuthenticated} />

<main class="results-page" bind:this={resultsPageElement}>
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

            <img
              class="club-hero-image"
              src={club.imagePath || defaultResultImage}
              alt={`Image for ${club.clubName}`}
            />

            <p class="match-score">Match score: {club.matchPercentage.toFixed(0)}%</p>
            <h2>{club.description || 'No description available yet.'}</h2>

            {#if club.meetingTime}
              <h3>Meeting Information</h3>
              <p>{club.meetingTime}</p>
            {/if}

            <div class="club-meta">
              {#if club.externalLink}
                <a href={club.externalLink} target="_blank" rel="noopener noreferrer">Visit website</a>
              {/if}

              {#if club.contactInfo}
                <p><strong>Contact:</strong> {club.contactInfo}</p>
              {/if}

              
            </div>
          </article>
        {/each}

        <div class="pager">
          <button onclick={prevPage} disabled={pageNum === 1}>Previous Organizations</button>
          <button onclick={nextPage} disabled={pageNum >= getTotalPages()}>Next Organizations</button>
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
    height: clamp(180px, 28vw, 320px);
    object-fit: contain;
    object-position: center;
    background: rgba(8, 22, 44, 0.25);
    border-radius: 0.8rem;
    margin: 0.9rem 0;
    border: none;
  }

  h3 {
    margin: 1rem 0 0.5rem 0;
    font-size: 1rem;
    letter-spacing: 0.01em;
  }

  p {
    font-size: 0.98rem;
    line-height: 1.55;
  }

  .club-meta {
    display: grid;
    gap: 0.3rem;
    margin-top: 0.35rem;
  }

  .club-meta a {
    width: fit-content;
    color: var(--action);
    font-weight: 600;
    text-decoration: none;
  }

  .club-meta a:hover {
    text-decoration: underline;
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