<script>
  import{ onMount } from 'svelte';
  import Header from '../../../Components/header.svelte'
  import Footer from '../../../Components/footer.svelte'
  let results = $state(["Computer Science Club"]) // for now this is just a placeholder, in the future this will be an array of clubs that match the user's interests and demographics, fetched from the backend server when the page loads
  let pageNum = $state(1) // for keeping track of what page of results the user is on. Not currently being used but will be helpful for future implementation 

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
    
</script>
  
<Header />
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

<Footer />
<style>
  main {
    flex: 1;
    max-width: 100%;
    margin: 0;
    padding: 0.5rem 1rem;
    font-family: system-ui, sans-serif;
    text-align: left;
  }

  h1 {
    font-size: 1.5rem;
    margin-top: 0;
    margin-bottom: 1.5rem;
    text-align: left;
  }

  form {
    display: grid;
    gap: 1rem;
    max-width: 600px;
    margin: 0 auto;
    text-align: left;
  }

  label {
    font-weight: 500;
    font-size: 1rem;
    margin-bottom: 0.25rem;
    display: block;
    text-align: left;
  }

  fieldset {
    border: none;
    border-top: 1px solid #ddd;
    padding: 1.5rem 0;
    margin: 0;
  }

  legend {
    font-weight: 600;
    font-size: 1rem;
    padding: 0;
    margin-bottom: 1rem;
    color: #333;
    text-align: left;
  }

  fieldset label {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin: 0.75rem 0;
    font-weight: normal;
    cursor: pointer;
  }

  input[type="radio"] {
    width: 1.25rem;
    height: 1.25rem;
    cursor: pointer;
  }

  input[type="checkbox"] {
    width: 1.25rem;
    height: 1.25rem;
    cursor: pointer;
  }

  button {
    width: 100%;
    padding: 1rem;
    font-size: 1.125rem;
    font-weight: 600;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    touch-action: manipulation;
  }

  button:hover {
    background-color: #2980b9;
  }

  button:active {
    transform: scale(0.98);
  }

  /* Desktop styles */
  @media (min-width: 768px) {
    main {
      margin: 0 auto;
      padding: 2rem;
    }

    h1 {
      font-size: 2rem;
    }

    button {
      width: fit-content;
      min-width: 200px;
    }
  }
</style>
