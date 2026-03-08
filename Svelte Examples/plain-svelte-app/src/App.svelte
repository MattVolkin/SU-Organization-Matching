<script>

  import Header from '../../../Components/header.svelte'
  import Footer from '../../../Components/footer.svelte'
  //setup for the questions
  let genderoptions = ["Man", "Woman", "Non-Binary", "Other", "Prefer not to say"]
  let raceoptions = ["American Native/Alaska Native", "Asian", "Black or African American", "Hispanic or Latino", "Middle Eastern or North African", "Native Hawaiian or Pacific Islander", "White", "Prefer not to say"]
  let religionoptions = ["Protestantism", "Catholocism", "Judaism", "Islam", "Buddhism", "Hinduism", "No religion", "Other", "Prefer not to say"]
  let majoroptions = ["Anthropology", "Applied Physics", "Art (Studio)", "Art History", "Biochemistry", "Biology", "Business", "Chemistry", "Classics", "Communication Studies", "Computational Mathematics", "Computer Science", "Pre-Dentistry", "Economics", "Education", "Pre-Engineering", "English", "Environmental Studies", "Feminist Studies", "French", "German", "Greek", "Health Professions", "History", "International Studies", "Kinesiology", "Latin", "Latin American & Border Studies", "Pre-Law", "Mathematics", "Pre-Med", "Pre-Ministry", "Music", "Pre-Nursing", "Pre-Occupational Therapy", "Philosophy", "Physics", "Political Science", "Psychology", "Pre-Physician Assistant", "Pre-Physical Therapy", "Religion", "Sociology", "Spanish", "Theatre", "Undecided"]
  //feilds in the form
  let name = ''
  let club = ''
  let contact = ''
  let officer = ""
  let email = ''
  let gender = ''
  let race = ''
  let religion = ''
  let major = ''
/**@function submit - Handles form submission */
  function submit() {
    console.log(`Name: ${name}, Club: ${club}, Contact: ${contact}, Officer : ${officer}, Email: ${email}, Gender: ${gender}, Race: ${race}, Religion: ${religion}, Major: ${major}`)//for debuging logs the form data
    fetch('http://localhost:8080/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name, club, contact, officer, email, gender, race, religion, major })//sends the form data to the backend server as a JSON object
    })
    .then(response => response.json())
    .then(data => {
      console.log('Success:', data)
    })
    .catch((error) => {
      console.error('Error:', error)
    })
  }
</script>

<Header />
<main>
  <h1>Demographic Form</h1>

  <form on:submit|preventDefault={submit}>
    <label for="name">Name:</label>
    <input id="name" type="text" bind:value={name} required /> <!-- Text input for the user's name, bound to the 'name' variable -->

    <fieldset>
      <legend>Gender</legend>
      {#each genderoptions as option} <!-- Loop through gender options to create radio buttons -->
        <label>
          <input type="radio" name="gender" value={option} bind:group={gender} required />
          {option}
        </label>
      {/each}
    </fieldset>

    <fieldset>
      <legend>Race/Ethnicity (Select all that apply)</legend> <!-- Fieldset for race/ethnicity options -->
      {#each raceoptions as option} <!-- Loop through race options to create checkboxes -->
        <label>
          <input type="checkbox" name="race" value={option} bind:group={race} required />     
          {option}
        </label>
      {/each}
    </fieldset>

    <fieldset>
      <legend>Intended Major(s)/Program of Study</legend>
      <select multiple bind:value={major} required>
        <option value="">Select a major...</option>
        {#each majoroptions as option} <!-- Loop through major options to create dropdown options -->
          <option value={option}>{option}</option>
        {/each}
      </select>
    </fieldset>

    <fieldset>
      <legend>Religion</legend>
      {#each religionoptions as option} <!-- Loop through religion options to create radio buttons -->
        <label>
          <input type="radio" name="religion" value={option} bind:group={religion} required />
          {option}
        </label>
      {/each}
    </fieldset>




    <button type="submit">Submit</button>
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
