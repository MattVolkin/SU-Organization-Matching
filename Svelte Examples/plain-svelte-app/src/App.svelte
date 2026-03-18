<script>
  import { onDestroy, onMount } from 'svelte';
  import Header from '../../../Components/header.svelte'
  import Footer from '../../../Components/footer.svelte'


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

  const getDemographics = {
    name: () => name,
    gender: () => gender,
    religion: () => religion,
    race: () => race,
    major: () => major,
  }

  const setDemographicsFields = {
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
    console.log('Applying prefill fields:', fields)
    for (const [key, value] of Object.entries(fields || {})) {
      const getCurrent = getDemographics[key]
      const setValue = setDemographicsFields[key]
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
        body: JSON.stringify({ name, gender, race, religion, major})//sends the form data to the backend server as a JSON object
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

  <form onsubmit={submitDemographics}>
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
