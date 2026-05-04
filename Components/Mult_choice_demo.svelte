<!-- @component Demographic quiz form component collecting user information for organization matching (gender, race, religion, major, etc.). -->
<script>
/**
 * @type {state} name - User's name input field
 * @type {state} gender - Selected gender preference
 * @type {state} race - Array of selected racial/ethnic identities
 * @type {state} religion - Selected religious affiliation
 * @type {state} major - Array of selected academic majors
 * @type {state} LGBTQ - User's LGBTQ identification status
 * @type {state} disability - User's disability status
 * @type {state} majorRequiredProxy - Helper state to track if major selection is required
 * @type {constant} genderOptions - Array of gender presentation options
 * @type {constant} raceOptions - Array of racial and ethnic identity options
 * @type {constant} religionOptions - Array of religious affiliation options
 * @type {constant} majorOptions - Comprehensive list of academic majors available
 * @type {constant} yes_no_options - Standard yes/no/prefer not to say options
 * @function submitDemographics - Validates and submits demographic form to backend API
 */
import { APICreater } from './APIHandler.svelte';
import MultiSelectDropdown from './MultiSelectDropdown.svelte';

  let name = $state('')
	let gender = $state('')
  let race = $state([])
	let religion = $state('')
  let major = $state([])
  let LGBTQ = $state("")
  let disability = $state("")
  let majorRequiredProxy = $state(null)

  const genderOptions = [
    'Man',
    'Woman',
    'Non-Binary',
    'Other',
    'Prefer not to say',
  ]

  const raceOptions = [
    'American Native/Alaska Native',
    'Asian',
    'Black or African American',
    'Hispanic or Latino',
    'Middle Eastern or North African',
    'Native Hawaiian or Pacific Islander',
    'White',
    'Prefer not to say',
  ]
  let yes_no_options = ['Yes', 'No', 'Prefer not to say']

  const religionOptions = [
    'Buddhism',
    'Catholicism',
    'Hinduism',
    'Islam',
    'Judaism',
    'Protestantism',    
    'No Religion',
    'Other',
    'Prefer not to say',
  ]

  const majorOptions = [
    'Anthropology',
    'Applied Physics',
    'Art (Studio)',
    'Art History',
    'Biochemistry',
    'Biology',
    'Business',
    'Chemistry',
    'Classics',
    'Communication Studies',
    'Computational Mathematics',
    'Computer Science',
    'Pre-Dentistry',
    'Economics',
    'Education',
    'Pre-Engineering',
    'English',
    'Environmental Studies',
    'Feminist Studies',
    'French',
    'German',
    'Greek',
    'Health Professions',
    'History',
    'International Studies',
    'Kinesiology',
    'Latin',
    'Latin American & Border Studies',
    'Pre-Law',
    'Mathematics',
    'Pre-Med',
    'Pre-Ministry',
    'Music',
    'Pre-Nursing',
    'Pre-Occupational Therapy',
    'Philosophy',
    'Physics',
    'Political Science',
    'Psychology',
    'Pre-Physician Assistant',
    'Pre-Physical Therapy',
    'Religion',
    'Sociology',
    'Spanish',
    'Theatre',
    'Undecided',
  ]


  async function submit(event) {
    event.preventDefault()

    if (majorRequiredProxy) {
      majorRequiredProxy.setCustomValidity('')
    }

    if (!Array.isArray(major) || major.length === 0) {
      if (majorRequiredProxy) {
        majorRequiredProxy.setCustomValidity('Please select one of these options.')
        majorRequiredProxy.reportValidity()
      }
      return
    }

    await APICreater('POST', '/submit', {
      "name": name,
      "gender": gender,
      "race": race,
      "religion": religion,
      "major": major,
      "lgbtq": LGBTQ,
     "disabilities": disability
    })
     window.location.href = '/swiping.html';
  }

  $effect(() => {
    if (majorRequiredProxy && Array.isArray(major) && major.length > 0) {
      majorRequiredProxy.setCustomValidity('')
    }
  })
</script>


<main>
  <section class="settings-page">
    <h2>Demographics</h2>
   

    <section class="club-settings-card">


      <form onsubmit={submit}>
        <fieldset>
          <legend>Gender Identity</legend>
          <div class="option-stack">
            {#each genderOptions as option}
              <label>
                <input type="radio" name="gender" value={option} bind:group={gender} required />
                <span>{option}</span>
              </label>
            {/each}
          </div>
        </fieldset>

        <fieldset>
          <legend>Race / Ethnicity</legend>
          <p class="hint">Select all that apply.</p>
          <div class="option-grid">
            {#each raceOptions as option}
              <label>
                <input type="checkbox" name="race" value={option} bind:group={race} />
                <span>{option}</span>
              </label>
            {/each}
          </div>
        </fieldset>

        <fieldset>
          <legend>Religion</legend>
          <div class="option-grid compact">
            {#each religionOptions as option}
              <label>
                <input type="radio" name="religion" value={option} bind:group={religion} required />
                <span>{option}</span>
              </label>
            {/each}
          </div>
        </fieldset>
        <fieldset>
        <legend> Do you Identify as LGBTQ+?</legend>
        <div class="option-grid compact">
          {#each yes_no_options as option}
            <label>
              <input type="radio" name="LGBTQ" value={option} bind:group={LGBTQ} required />
              <span>{option}</span>
            </label>  
          {/each}
        </div>
        </fieldset>
    
        <fieldset>
        <legend>Do you have any Disabilities?</legend>
        <div class="option-grid compact">
          {#each yes_no_options as option}
            <label>
              <input type="radio" name="disability" value={option} bind:group={disability} required />
              <span>{option}</span>
            </label>  
          {/each}
        </div>
        </fieldset>

        <fieldset>
			<legend>Intended Major(s) / Program of Study</legend>
			<p class="hint">Select all that apply.</p>
      <div class="major-required-wrapper">
			<MultiSelectDropdown id="majors" label="Choose majors" options={majorOptions} bind:value={major} required />
        <input
          class="major-required-proxy"
          type="text"
          tabindex="-1"
          bind:this={majorRequiredProxy}
          value={major.join(',')}
          required
          aria-label="Major"
        />
      </div>
		</fieldset>

        <button type="submit">Submit and Continue</button>
      </form>
    </section>
  </section>
</main>



<style>
  :global(html),
  :global(body) {
    margin: 0;
    padding: 0;
    background: #ffffff;
  }

  main {
    --page-bg: linear-gradient(180deg, #ffffff 0%, #f3f3f3 100%);
    --text-main: #000000;
    --text-subtle: #454545;
    --card-bg: #ffffff;
    --card-border: #828282;
    --accent: #000000;
    --accent-hover: #828282;
    --input-bg: #ffffff;
    --input-border: #828282;
    --input-accent: #000000;
    --shadow: rgba(0, 0, 0, 0.1);

    width: 100%;
    margin: 0;
    padding: 0;
    font-family: system-ui, sans-serif;
    background: var(--page-bg);
    color: var(--text-main);
  }

  .settings-page {
    width: min(100%, 1040px);
    margin: 0 auto;
    padding: 1rem;
  }

  .settings-page h2 {
    margin: 0.75rem 0 0.85rem 0;
    font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
    color: var(--text-main);
  }

 
  .club-settings-card {
    border: 1px solid var(--card-border);
    border-radius: 0.85rem;
    padding: 1rem;
    margin-bottom: 1rem;
    background: var(--card-bg);
    box-shadow: 0 10px 24px var(--shadow);
    color: var(--text-main);
  }

  form {
    display: grid;
    gap: 1rem;
  }

  fieldset {
    margin: 0;
    border: 1px solid var(--card-border);
    border-radius: 1rem;
    padding: 1rem;
    background: var(--card-bg);
    color: var(--text-main);
  }

  legend {
    padding: 0 0.45rem;
    font-weight: 800;
    color: var(--text-main);
  }

  .hint {
    margin: 0.35rem 0 0.75rem;
    color: var(--text-subtle);
    font-size: 0.92rem;
  }

  .major-required-wrapper {
    position: relative;
  }

  .major-required-proxy {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: 0;
    border: 0;
    opacity: 0;
    pointer-events: none;
  }

  .option-stack,
  .option-grid {
    display: grid;
    gap: 0.55rem;
  }

  .option-grid {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }

  .option-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  label {
    display: flex;
    align-items: flex-start;
    gap: 0.65rem;
    padding: 0.72rem 0.8rem;
    border-radius: 0.85rem;
    background: var(--input-bg);
    border: 1px solid var(--input-border);
    box-shadow: 0 8px 16px var(--shadow);
    cursor: pointer;
    transition: transform 0.18s ease, border-color 0.18s ease, box-shadow 0.18s ease;
    color: var(--text-main);
  }

  label:hover {
    transform: translateY(-1px);
    border-color: #ffcd00;
    box-shadow: 0 12px 22px rgba(0, 0, 0, 0.15);
  }

  label span {
    line-height: 1.35;
  }

  input[type='radio'],
  input[type='checkbox'] {
    margin-top: 0.2rem;
    accent-color: var(--input-accent);
    flex: 0 0 auto;
  }

  button {
    width: fit-content;
    margin-top: 0.25rem;
    border: 1px solid transparent;
    border-radius: 0.55rem;
    padding: 0.6rem 0.95rem;
    font-weight: 700;
    color: #ffffff;
    background: var(--accent);
    box-shadow: 0 10px 24px var(--shadow);
    cursor: pointer;
    transition: transform 0.18s ease, box-shadow 0.18s ease, background-color 0.18s ease;
  }

  button:hover {
    transform: translateY(-1px);
    background: var(--accent-hover);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
  }

  

  @media (max-width: 880px) {
    main {
      padding: 0;
    }

    .settings-page {
      padding: 0.75rem;
    }

    .option-grid,
    .option-grid.compact {
      grid-template-columns: 1fr;
    }
  }

  @media (prefers-color-scheme: dark) {
    :global(html),
    :global(body) {
      background: #000000;
    }

    main {
      --page-bg: linear-gradient(180deg, #000000 0%, #1e1e1e 100%);
      --text-main: #ffffff;
      --text-subtle: #d5d5d5;
      --card-bg: #121212;
      --card-border: #828282;
      --accent: #ffcd00;
      --accent-hover: #e5b800;
      --input-bg: #1e1e1e;
      --input-border: #828282;
      --input-accent: #ffcd00;
      --shadow: rgba(0, 0, 0, 0.3);
    }

    label {
      color: var(--text-main);
    }

    button {
      color: #000000;
    }
  }
</style>