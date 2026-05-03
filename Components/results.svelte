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
  import { APICreater } from './APIHandler.svelte';

  const defaultResultImage = "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 640 360'><defs><linearGradient id='sky' x1='0' y1='0' x2='0' y2='1'><stop offset='0%' stop-color='%23dbeafe'/><stop offset='100%' stop-color='%23f8fafc'/></linearGradient></defs><rect width='640' height='360' fill='url(%23sky)'/><rect y='235' width='640' height='125' fill='%23d1d5db'/><path d='M0 250L120 190L210 235L330 165L440 235L520 205L640 250V360H0Z' fill='%239ca3af'/><circle cx='110' cy='84' r='34' fill='%23fde68a'/><text x='320' y='330' text-anchor='middle' font-family='Arial, sans-serif' font-size='22' fill='%234b5563'>Organization Image</text></svg>";
 

  let results = $state([]) // stores matched clubs returned by the results API
  let pageNum = $state(1) // for keeping track of what page of results the user is on.
  let isAuthChecking = $state(true)
  let isAuthenticated = $state(false)
  let activeRole = $state('user')
  let resultsErrorMessage = $state('')
  let resultsPageElement = $state(null)
  let ClubsPerPage = 5
  let threshold =55
  let maxClub=20
  let minClub=5

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

  function getClubDisplayKey(club) {
    const clubId = getClubID(club)
    if (clubId > 0) {
      return `id:${clubId}`
    }

    return `name:${normalizeClubKey(getClubName(club))}`
  }

  function extractListPayload(payload) {
    if (Array.isArray(payload)) {
      return payload
    }

    if (!payload || typeof payload !== 'object') {
      return []
    }

    const candidateKeys = ['results', 'data', 'organizations', 'orgs', 'items']
    for (const key of candidateKeys) {
      if (Array.isArray(payload[key])) {
        return payload[key]
      }
    }

    return []
  }

  function extractErrorMessage(payload) {
    if (!payload || typeof payload !== 'object') {
      return ''
    }

    if (typeof payload.error === 'string' && payload.error.trim()) {
      return payload.error.trim()
    }

    if (typeof payload.message === 'string' && payload.message.trim()) {
      const lower = payload.message.trim().toLowerCase()
      if (lower.includes('error') || lower.includes('unauthorized') || lower.includes('failed')) {
        return payload.message.trim()
      }
    }

    return ''
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

    return ''
  }

  function getVisibleContactInfo(value) {
    const text = typeof value === 'string' ? value.trim() : ''
    if (!text) {
      return ''
    }

    if (isOfficerPlaceholder(text)) {
      return ''
    }

    const officerEmailsMatch = text.match(/officer emails?\s*:\s*(.+)$/i)
    const emailSource = officerEmailsMatch ? officerEmailsMatch[1] : text
    const emailMatches = emailSource.match(/[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}/gi)

    if (emailMatches && emailMatches.length > 0) {
      return emailMatches.join(', ')
    }

    return ''
  }

  function getVisibleExternalLink(value) {
    const text = typeof value === 'string' ? value.trim() : ''
    if (!text || isOfficerPlaceholder(text)) {
      return ''
    }

    return text
  }

  function getVisibleActivitiesDescription(club) {
    const description = getTrimmedText(
      club?.activitiesDescreption ??
      club?.activitiesDescription ??
      club?.ActivitiesDescreption ??
      club?.ActivitiesDescription
    )

    if (description) {
      return description
    }

    const activities = club?.activities ?? club?.Activities
    if (Array.isArray(activities)) {
      return activities
        .map((entry) => String(entry || '').trim())
        .filter((entry) => entry.length > 0)
        .join(', ')
    }

    if (typeof activities === 'string') {
      return activities
        .split(',')
        .map((entry) => entry.trim())
        .filter((entry) => entry.length > 0)
        .join(', ')
    }

    return ''
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
    const resultsInThreshold = results.filter((club) => getClubScore(club) >= threshold)
    const selectedKeys = new Set(resultsInThreshold.map((club) => getClubDisplayKey(club)))
    const backfillResults = results.filter((club) => !selectedKeys.has(getClubDisplayKey(club)))
    const combinedResults = [...resultsInThreshold]

    for (const club of backfillResults) {
      if (combinedResults.length >= minClub && combinedResults.length >= maxClub) {
        break
      }

      combinedResults.push(club)
    }

    return combinedResults.slice(0, maxClub)

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
      const data = await response.json().catch(() => ({}))
      const role = String(data?.role || '').toLowerCase()
      activeRole = role === 'admin' || role === 'officer' ? role : 'user'

      isAuthenticated = true;
      isAuthChecking = false;
      await getResults();
      return;
    }

    isAuthenticated = false
    activeRole = 'user'
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
    resultsErrorMessage = ''

    try {
      const [detailResult, scoreResult] = await Promise.allSettled([
        APICreater('GET', '/api/results', null),
        APICreater('GET', '/api/results?includeScores=true', null),
      ])

      const detailPayload = detailResult.status === 'fulfilled' ? detailResult.value : []
      const scorePayload = scoreResult.status === 'fulfilled' ? scoreResult.value : []

      let detailList = extractListPayload(detailPayload)
      let scoreList = extractListPayload(scorePayload)

      const detailError = extractErrorMessage(detailPayload)
      const scoreError = extractErrorMessage(scorePayload)
      const requestErrors = [
        detailResult.status === 'rejected' ? 'Unable to load /api/results.' : '',
        scoreResult.status === 'rejected' ? 'Unable to load /api/results?includeScores=true.' : '',
        detailError,
        scoreError,
      ].filter((entry) => entry)

      if (detailList.length === 0 && scoreList.length === 0 && requestErrors.length > 0) {
        resultsErrorMessage = requestErrors[0]
      }

      if (detailList.length === 0 && scoreList.length === 0 && !resultsErrorMessage && (activeRole === 'admin' || activeRole === 'officer')) {
        const fallbackEndpoint = activeRole === 'admin' ? '/api/admin/orgs' : '/api/officer/orgs'
        const fallbackPayload = await APICreater('GET', fallbackEndpoint, null)
        detailList = extractListPayload(fallbackPayload)
        scoreList = []
      }

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
          activitiesDescription: getVisibleActivitiesDescription(mergedClub),
          imagePath: getClubImagePath(mergedClub),
          externalLink: getVisibleExternalLink(mergedClub?.externalLink ?? mergedClub?.ExternalLink),
          contactInfo: getVisibleContactInfo(mergedClub?.contactInfo ?? mergedClub?.ContactInfo),
          includeOfficerEmails: Boolean(mergedClub?.includeOfficerEmails ?? mergedClub?.IncludeOfficerEmails),
          updatedAt: typeof (mergedClub?.updatedAt ?? mergedClub?.UpdatedAt) === 'string'
            ? (mergedClub?.updatedAt ?? mergedClub?.UpdatedAt).trim()
            : '',
        }
      })

      pageNum = 1
    } catch (error) {
      console.error('Unable to load results', error)
      results = []
      resultsErrorMessage = 'Unable to load results right now. Please refresh and try again.'
    }
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
      {#if resultsErrorMessage}
        <p>{resultsErrorMessage}</p>
      {:else if getVisibleResults().length === 0}
        <p>Please fill out the sorting quiz to see organization results.</p>
      {:else}
        {#each getVisibleResults() as club, index (club.id || `${club.clubName}-${index}`)}
          <article class="club-card">
            <h1>{club.clubName}</h1>
            <p class="match-score">Match: {Math.round(getClubScore(club))}%</p>

            <img
              class="club-hero-image"
              src={club.imagePath || defaultResultImage}
              alt={`Image for ${club.clubName}`}
            />

            <h2>{club.description || 'No description available yet.'}</h2>

            {#if club.meetingTime}
              <h3>Meeting Information</h3>
              <p>{club.meetingTime}</p>
            {/if}

            {#if club.activitiesDescription}
              <h3>Activities</h3>
              <p>{club.activitiesDescription}</p>
            {/if}

            <div class="club-meta">
              {#if club.externalLink}
                <a href={club.externalLink} target="_blank" rel="noopener noreferrer">Visit website</a>
              {/if}

              {#if club.contactInfo}
                <p><strong>Officer Emails:</strong> {club.contactInfo}</p>
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
    background: #ffffff;
  }

  .results-page {
    --page-bg: linear-gradient(180deg, #ffffff 0%, #f3f3f3 100%);
    --card-bg: #ffffff;
    --card-border: #828282;
    --text-main: #000000;
    --text-subtle: #4a4a4a;
    --action: #000000;
    --action-hover: #828282;
    --focus-ring: #ffcd00;

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
    box-shadow: 0 8px 18px rgba(0, 0, 0, 0.12);
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
    background: rgba(130, 130, 130, 0.2);
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
    :global(html),
    :global(body) {
      background: #000000;
    }

    .results-page {
      --page-bg: linear-gradient(180deg, #000000 0%, #1e1e1e 100%);
      --card-bg: #121212;
      --card-border: #828282;
      --text-main: #ffffff;
      --text-subtle: #d5d5d5;
      --action: #ffcd00;
      --action-hover: #e5b800;
      --focus-ring: #ffffff;
    }

    .status-card,
    .result-card {
      box-shadow: 0 14px 28px rgba(0, 0, 0, 0.45);
    }

    .club-meta a {
      color: #ffcd00;
    }

    button {
      color: #000000;
    }

    button:disabled {
      opacity: 0.5;
      background: #828282;
      color: #ffffff;
    }
  }
</style>