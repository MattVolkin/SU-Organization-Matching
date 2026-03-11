<script>
  let name = ''
  let club = ''
  let contact = ''
  let officer = ""
  let email = ''
	let gender = ''
	let race = ''
	let religion = ''
  let major = ''


  function submit() {
    console.log(`Name: ${name}, Club: ${club}, Contact: ${contact}, Officer : ${officer}, Email: ${email}`)
    fetch('/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name, club, contact, officer, email })
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

<main>
  <h1>Club Form</h1>

  <form on:submit|preventDefault={submit}>
    
    <fieldset>
      <legend>Gender</legend>
      <label>
        <input type="radio" name="gender" value="Man" bind:group={gender} required />
        Man <br>
      </label>
      <label>
        <input type="radio" name="gender" value="Woman" bind:group={gender} required />
        Woman <br>
      </label>
			<label>
        <input type="radio" name="gender" value="Non-Binary" bind:group={gender} required />
        Non-Binary <br>
      </label>
			<label>
        <input type="radio" name="gender" value="Other" bind:group={gender} required />
        Other <br>
      </label>
			<label>
        <input type="radio" name="gender" value="Prefer not to say" bind:group={gender} required />
        Prefer not to say
      </label>
    </fieldset>
		
    <fieldset>
      <legend>Race/Ethnicity (Select all that apply)</legend>
      <label>
        <input type="checkbox" name="race" value="American Native/Alaska Native" bind:group={race} required />
        American Native/Alaska Native <br>
      </label>
			<label>
        <input type="checkbox" name="race" value="Asian" bind:group={race} required />
        Asian <br>
      </label>
			<label>
        <input type="checkbox" name="race" value="Black or African American" bind:group={race} required />
        Black or African American <br>
      </label>
      <label>
        <input type="checkbox" name="race" value="Hispanic or Latino" bind:group={race} required />
        Hispanic or Latino <br>
      </label>
			<label>
        <input type="checkbox" name="race" value="Middle Eastern or North African" bind:group={race} required />
        Middle Eastern or North African <br>
      </label>
			<label>
        <input type="checkbox" name="race" value="Native Hawaiian or Pacific Islander" bind:group={race} required />
        Native Hawaiian or Pacific Islander <br>
      </label>
			<label>
        <input type="checkbox" name="race" value="White" bind:group={race} required />
        White <br>
      </label>
			<label>
        <input type="radio" name="race" value="Prefer not to say" bind:group={race} required />
        Prefer not to say <br>
      </label>
    </fieldset>

    <fieldset>
      <legend>Religion</legend>
      <label>
        <input type="radio" name="religion" value="Protestantism" bind:group={religion} required />
        Protestantism <br>
      </label>
      <label>
        <input type="radio" name="religion" value="Catholocism" bind:group={religion} required />
        Catholocism <br>
      </label>
			<label>
        <input type="radio" name="religion" value="Judaism" bind:group={religion} required />
        Judaism <br>
      </label>
				<label>
        <input type="radio" name="religion" value="Islam" bind:group={religion} required />
        Islam <br>
				</label>
			<label>
        <input type="radio" name="religion" value="Buddhism" bind:group={religion} required />
        Buddhism <br>
      </label>
			<label>
        <input type="radio" name="religion" value="Hinduism" bind:group={religion} required />
        Hinduism <br>
      </label>
			<label>
        <input type="radio" name="religion" value="No religion" bind:group={religion} required />
        No religion <br>
      </label>
			<label>
        <input type="radio" name="religion" value="Other" bind:group={religion} required />
        Other <br>
      </label>
			<label>
        <input type="radio" name="religion" value="Prefer not to say" bind:group={religion} required />
        Prefer not to say <br>
      </label>
    </fieldset>

    <fieldset>
			<legend>
				Intended Major(s)/Program of Study
			</legend>
			<label>
        <input type="checkbox" name="major" value="Anthropology" bind:group={major} required />
        Anthropology <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Applied Physics" bind:group={major} required />
        Applied Physics<br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Art (Studio)" bind:group={major} required />
        Art (Studio)<br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Art History" bind:group={major} required />
        Art History <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Biochemistry" bind:group={major} required />
        Biochemistry <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Biology" bind:group={major} required />
        Biology <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Business" bind:group={major} required />
        Business <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Chemistry" bind:group={major} required />
        Chemistry <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Classics" bind:group={major} required />
        Classics <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Communication Studies" bind:group={major} required />
        Communication Studies <br>
      </label>
      			<label>
        <input type="checkbox" name="major" value="Computational Mathematics" bind:group={major} required />
        Computational Mathematics <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Computer Science" bind:group={major} required />
        Computer Science <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Dentistry" bind:group={major} required />
        Pre-Dentistry <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Economics" bind:group={major} required />
        Economics <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Education" bind:group={major} required />
        Education <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Engineering" bind:group={major} required />
        Pre-Engineering <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="English" bind:group={major} required />
        English <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Environmental Studies" bind:group={major} required />
        Environmental Studies <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Feminist Studies" bind:group={major} required />
        Feminist Studies <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="French" bind:group={major} required />
        French <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="German" bind:group={major} required />
        German <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Greek" bind:group={major} required />
        Greek <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Health Professions" bind:group={major} required />
        Health Professions <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="History" bind:group={major} required />
        History <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="International Studies" bind:group={major} required />
        International Studies <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Kinesiology" bind:group={major} required />
        Kinesiology <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Latin" bind:group={major} required />
        Latin <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Latin American & Border Studies" bind:group={major} required />
        Latin American & Border Studies <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Law" bind:group={major} required />
        Pre-Law <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Mathematics" bind:group={major} required />
        Mathematics <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Med" bind:group={major} required />
        Pre-Med <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Ministry" bind:group={major} required />
        Pre-Ministry <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Music" bind:group={major} required />
        Music <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Nursing" bind:group={major} required />
        Pre-Nursing <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Occupational Therapy" bind:group={major} required />
        Pre-Occupational Therapy <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Philosophy" bind:group={major} required />
        Philosophy <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Physics" bind:group={major} required />
        Physics <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Political Science" bind:group={major} required />
        Political Science <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Psychology" bind:group={major} required />
        Psychology <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Physician Assistant" bind:group={major} required />
        Pre-Physician Assistant <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Pre-Physical Therapy" bind:group={major} required />
        Pre-Physical Therapy <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Religion" bind:group={major} required />
        Religion <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Sociology" bind:group={major} required />
        Sociology <br>
      </label>
			<label>
        <input type="checkbox" name="major" value="Spanish" bind:group={major} required />
        Spanish <br>
      </label>
    	<label>
        <input type="checkbox" name="major" value="Theatre" bind:group={major} required />
        Theatre <br>
      </label>
    	<label>
        <input type="radio" name="major" value="Undecided" bind:group={major} required />
        Undecided <br>
      </label>
		</fieldset>

    <button type="submit">Submit</button>
  </form>
</main>

<style>
  main {
    max-width: 500px;
    margin: 2rem auto;
    font-family: system-ui, sans-serif;
  }

  form {
    display: grid;
    gap: 0.75rem;
  }

  fieldset {
    border: 1px solid #ccc;
    padding: 0.75rem;
  }

  button {
    width: fit-content;
  }
</style>