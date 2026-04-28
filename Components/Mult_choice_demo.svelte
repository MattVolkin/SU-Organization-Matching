<script>
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
  }

  main {
    width: 100%;
    margin: 0;
    padding: 0;
    font-family: system-ui, sans-serif;
    background: #ffffff;
    color: #10243a;
  }

  .settings-page {
    width: min(100%, 1040px);
    margin: 0 auto;
    padding: 1rem;
  }

  .settings-page h2 {
    margin: 0.75rem 0 0.85rem 0;
    font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
    color: #132c45;
  }

 
  .club-settings-card {
    border: 1px solid #d7dee8;
    border-radius: 0.85rem;
    padding: 1rem;
    margin-bottom: 1rem;
    background: #fafcff;
    box-shadow: 0 10px 24px rgba(13, 37, 62, 0.08);
  }

  form {
    display: grid;
    gap: 1rem;
  }

  fieldset {
    margin: 0;
    border: 1px solid #d7dee8;
    border-radius: 1rem;
    padding: 1rem;
    background: #fafcff;
  }

  legend {
    padding: 0 0.45rem;
    font-weight: 800;
    color: #132c45;
  }

  .hint {
    margin: 0.35rem 0 0.75rem;
    color: #31506e;
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
    background: #ffffff;
    border: 1px solid #dbe7f3;
    box-shadow: 0 8px 16px rgba(16, 36, 58, 0.04);
    cursor: pointer;
    transition: transform 0.18s ease, border-color 0.18s ease, box-shadow 0.18s ease;
  }

  label:hover {
    transform: translateY(-1px);
    border-color: #bfdde8;
    box-shadow: 0 12px 22px rgba(16, 36, 58, 0.08);
  }

  label span {
    line-height: 1.35;
  }

  input[type='radio'],
  input[type='checkbox'] {
    margin-top: 0.2rem;
    accent-color: #1f6f8b;
    flex: 0 0 auto;
  }

  button {
    width: fit-content;
    margin-top: 0.25rem;
    border: none;
    border-radius: 0.55rem;
    padding: 0.6rem 0.95rem;
    font-weight: 700;
    color: #ffffff;
    background: #0f6d8c;
    box-shadow: 0 10px 24px rgba(13, 37, 62, 0.08);
    cursor: pointer;
    transition: transform 0.18s ease, box-shadow 0.18s ease;
  }

  button:hover {
    transform: translateY(-1px);
    background: #0d5f79;
    box-shadow: 0 12px 24px rgba(13, 37, 62, 0.12);
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
</style>