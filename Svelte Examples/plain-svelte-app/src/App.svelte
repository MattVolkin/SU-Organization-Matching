<script>
  import { onDestroy, onMount } from 'svelte';
  import Header from '../../../Components/header.svelte'
  import Footer from '../../../Components/footer.svelte'

  let results = $state(["Computer Science Club"]) // for now this is just a placeholder, in the future this will be an array of clubs that match the user's interests and demographics, fetched from the backend server when the page loads
  let pageNum = $state(1) // for keeping track of what page of results the user is on. Not currently being used but will be helpful for future implementation 

  async function getResults() {
    const response = await fetch('/results')
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
    


  //setup for the questions
  let genderoptions = ["Man", "Woman", "Non-Binary", "Other", "Prefer not to say"]
  let raceoptions = ["American Native/Alaska Native", "Asian", "Black or African American", "Hispanic or Latino", "Middle Eastern or North African", "Native Hawaiian or Pacific Islander", "White", "Prefer not to say"]
  let religionoptions = ["Protestantism", "Catholocism", "Judaism", "Islam", "Buddhism", "Hinduism", "No religion", "Other", "Prefer not to say"]
  let majoroptions = ["Anthropology", "Applied Physics", "Art (Studio)", "Art History", "Biochemistry", "Biology", "Business", "Chemistry", "Classics", "Communication Studies", "Computational Mathematics", "Computer Science", "Pre-Dentistry", "Economics", "Education", "Pre-Engineering", "English", "Environmental Studies", "Feminist Studies", "French", "German", "Greek", "Health Professions", "History", "International Studies", "Kinesiology", "Latin", "Latin American & Border Studies", "Pre-Law", "Mathematics", "Pre-Med", "Pre-Ministry", "Music", "Pre-Nursing", "Pre-Occupational Therapy", "Philosophy", "Physics", "Political Science", "Psychology", "Pre-Physician Assistant", "Pre-Physical Therapy", "Religion", "Sociology", "Spanish", "Theatre", "Undecided"]
  //feilds in the form
  let name = $state('')
  let gender = $state('')
  let race = $state([])
  let religion = $state('')
  let major = $state([])
  let nameError = $state('')
  let genderError = $state('')
  let raceError = $state('')
  let majorError = $state('')
  let religionError = $state('')
  let submitError = $state('')
  const showSubmitError = true // Set to false to hide submit/auth errors for testing.

  const getFieldValues = {
    name: () => name,
    gender: () => gender,
    religion: () => religion,
    race: () => race,
    major: () => major,
  }

  const setFieldValues = {
    name: (value) => { name = value },
    gender: (value) => { gender = value },
    religion: (value) => { religion = value },
    race: (value) => { race = Array.isArray(value) ? value : race },
    major: (value) => { major = Array.isArray(value) ? value : major },
  }

  function isFieldEmpty(value) {
    if (Array.isArray(value)) {
      return value.length === 0
    }
    return !String(value ?? '').trim()
  }

  function applyPrefillFields(fields) {
    for (const [key, value] of Object.entries(fields || {})) {
      const getCurrent = getFieldValues[key]
      const setValue = setFieldValues[key]
      if (!getCurrent || !setValue) {
        continue
      }
      if (!isFieldEmpty(getCurrent())) {
        continue
      }
      setValue(value)
    }
  }

  async function loadPrefillFields() {
    try {
      const tokenFromStorage = localStorage.getItem('authToken') || ''
      const headers = tokenFromStorage
        ? { Authorization: `Bearer ${tokenFromStorage}` }
        : {}

      const response = await fetch('/api/prefill', {
        method: 'GET',
        credentials: 'include',
        headers,
      })
      if (!response.ok) {
        return
      }
      const data = await response.json().catch(() => ({}))
      applyPrefillFields(data?.fields || {})
    } catch (error) {
      console.error('Prefill fetch failed:', error)
    }
  }

  function onAuthMessage(event) {
    if (event.origin !== window.location.origin) {
      return
    }
    if (event.data?.type !== 'google-auth-success') {
      return
    }
    loadPrefillFields()
  }
/**@function submit - Handles form submission */
  async function submitDemographics() {
    nameError = ''
    genderError = ''
    raceError = ''
    majorError = ''
    religionError = ''
    submitError = ''

    if (!name.trim()) {
      nameError = 'Please enter your name.'
    }

    if (!gender) {
      genderError = 'Please select a gender option.'
    }

    if (race.length === 0) {
      raceError = 'Please select at least one race/ethnicity option.'
    }
    
    if (major.length === 0) {
      majorError = 'Please select at least one intended major.'
    }

    if (!religion) {
      religionError = 'Please select a religion option.'
    }

    if (nameError || genderError || raceError || majorError || religionError) {
      return
    }

    console.log(`Name: ${name}, Gender: ${gender}, Race: ${race}, Religion: ${religion}, Major: ${major}`)//for debuging logs the form data

    try {
      const response = await fetch('/submit', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name, gender, race, religion, major, })//sends the form data to the backend server as a JSON object
      })

      const data = await response.json().catch(() => ({}))

      if (!response.ok) {
        if (response.status === 401) {
          submitError = 'Please log in with Google before submitting this form.'
          return
        }
        submitError = data?.error || 'Unable to submit the form right now. Please try again.'
        return
      }

      console.log('Success:', data)
    } catch (error) {
      console.error('Error:', error)
      submitError = 'Network error while submitting. Please try again.'
    }
  }

  onMount(() => {
    window.addEventListener('message', onAuthMessage)
    loadPrefillFields()
  })

  onDestroy(() => {
    window.removeEventListener('message', onAuthMessage)
  })
</script>

<Header />

<main>
  <h1>Demographic Form</h1>

  <form on:submit|preventDefault={submitDemographics}>
    <label for="name">Name:</label>
    <input id="name" type="text" bind:value={name} aria-invalid={nameError ? 'true' : 'false'} /> <!-- Text input for the user's name, bound to the 'name' variable -->
    {#if nameError}
      <p class="error-message">{nameError}</p>
    {/if}

    <fieldset>
      <legend>Gender</legend>
      {#each genderoptions as option} <!-- Loop through gender options to create radio buttons -->
        <label>
          <input type="radio" name="gender" value={option} bind:group={gender} />
          {option}
        </label>
      {/each}
      {#if genderError}
        <p class="error-message">{genderError}</p>
      {/if}
    </fieldset>

    <fieldset>
      <legend>Race/Ethnicity (Select all that apply)</legend> <!-- Fieldset for race/ethnicity options -->
      {#each raceoptions as option} <!-- Loop through race options to create checkboxes -->
        <label>
          <input
            type="checkbox"
            name="race"
            value={option}
            bind:group={race}
          />
          {option}
        </label>
      {/each}
      {#if raceError}
        <p class="error-message">{raceError}</p>
      {/if}
    </fieldset>

    <fieldset>
      <legend>Intended Major(s)/Program of Study</legend>
      <p class="field-help">Select all that apply.</p>
      <div class="multi-select-grid">
        {#each majoroptions as option}
          <label>
            <input type="checkbox" name="major" value={option} bind:group={major} />
            {option}
          </label>
        {/each}
      </div>
      {#if majorError}
        <p class="error-message">{majorError}</p>
      {/if}
    </fieldset>

    <fieldset>
      <legend>Religion</legend>
      {#each religionoptions as option} <!-- Loop through religion options to create radio buttons -->
        <label>
          <input type="radio" name="religion" value={option} bind:group={religion} />
          {option}
        </label>
      {/each}
      {#if religionError}
        <p class="error-message">{religionError}</p>
      {/if}
    </fieldset>




    <button type="submit">Submit</button>
    {#if showSubmitError && submitError}
      <p class="error-message">{submitError}</p>
    {/if}
  </form>
</main>

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
    color: inherit;
  }

  @media (prefers-color-scheme: dark) {
    h1 {
      color: #e0e0e0;
    }
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
    color: inherit;
  }

  @media (prefers-color-scheme: dark) {
    label {
      color: #e0e0e0;
    }
  }

  fieldset {
    border: none;
    border-top: 1px solid #ddd;
    padding: 1.5rem 0;
    margin: 0;
  }

  @media (prefers-color-scheme: dark) {
    fieldset {
      border-top-color: #444;
    }
  }

  legend {
    font-weight: 600;
    font-size: 1rem;
    padding: 0;
    margin-bottom: 1rem;
    color: #333;
    text-align: left;
  }

  @media (prefers-color-scheme: dark) {
    legend {
      color: #e0e0e0;
    }
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

  .field-help {
    margin: 0 0 0.75rem;
    color: #555;
    font-size: 0.95rem;
  }

  @media (prefers-color-scheme: dark) {
    .field-help {
      color: #aaa;
    }
  }

  .multi-select-grid {
    display: grid;
    grid-template-columns: 1fr;
    max-height: 18rem;
    overflow-y: auto;
    border: 1px solid #ddd;
    border-radius: 6px;
    padding: 0.5rem 0.75rem;
  }

  @media (prefers-color-scheme: dark) {
    .multi-select-grid {
      border-color: #444;
    }
  }

  .error-message {
    margin: 0.75rem 0 0;
    padding: 0.75rem 1rem;
    background-color: #ffebee;
    color: #b00020;
    font-size: 0.92rem;
    border-radius: 20px;
    border: 1px solid #ffcdd2;
    display: inline-block;
    font-weight: 500;
  }

  @media (prefers-color-scheme: dark) {
    .error-message {
      background-color: #3a1a1a;
      color: #ff6b6b;
      border-color: #662222;
    }
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

    .multi-select-grid {
      grid-template-columns: 1fr 1fr;
    }
  }
</style>
