<script>
  import { onMount } from 'svelte';
  import Header from '../../../Components/header.svelte'
  import Footer from '../../../Components/footer.svelte'
  import LoginPopup from '../../../Components/login_popup.svelte'
  let results = $state(["Computer Science Club"]) // for now this is just a placeholder, in the future this will be an array of clubs that match the user's interests and demographics, fetched from the backend server when the page loads
  let pageNum = $state(1) // for keeping track of what page of results the user is on. Not currently being used but will be helpful for future implementation 
  let isAuthChecking = $state(true)
  let isAuthenticated = $state(false)
  let shouldPromptLogin = $state(false)

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
      shouldPromptLogin = false;
      return;
    }

    isAuthenticated = false;
    isAuthChecking = false;
    shouldPromptLogin = true;
  }

  function handleLoginSuccess(data) {
    if (data?.token) {
      localStorage.setItem('authToken', data.token);
    }
    isAuthenticated = true;
    isAuthChecking = false;
    shouldPromptLogin = false;
  }

  function handleLoginBlocked() {
    window.location.assign('/login');
  }

  async function getResults() {
    const response = await fetch('http://localhost:8080/results')
    results = await response.json() // in backend please make the list of results an array of clubs (sorted in decending order of match percentage) with each club having a name, description and match percentage field
  }

  async function nextPage() {
    pageNum += 1
    await getResults() // for future implementation when we have more results than we want to show on one page, this will fetch the next page of results from the backend
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
    promptLoginIfNeeded();
  });
    
</script>
  
<Header />
  <LoginPopup autoOpen={shouldPromptLogin} onSuccess={handleLoginSuccess} onBlocked={handleLoginBlocked} />
  {#if isAuthChecking}
    <p>Checking sign-in status...</p>
  {:else if !isAuthenticated}
    <p>Please complete sign-in to view organization information.</p>
  {:else}
    <h1>{results[pageNum-1]}</h1>
    <h2>Hi we are {results[pageNum-1]} and we are commited to to provide a safe space to play games and hang out with other computer nerds </h2>

    <button onclick={prevPage} disabled={pageNum === 1}>Previous</button>
    <button onclick={nextPage}>Next</button>
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
  {/if}

<Footer />