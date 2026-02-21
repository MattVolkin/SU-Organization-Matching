<script>
  let name = ''
  let club = ''
  let contact = ''
  let officer = ""

  function submit() {
    console.log(`Name: ${name}, Club: ${club}, Contact: ${contact}, Officer : ${officer}`)
    fetch('http://localhost:8080/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name, club, contact, officer })
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
    <label for="name">Name:</label>
    <input id="name" type="text" bind:value={name} required />

    <label for="club">Club:</label>
    <input id="club" type="text" bind:value={club} required />

    <fieldset>
      <legend>Preferred contact</legend>
      <label>
        <input type="radio" name="contact" value="email" bind:group={contact} required />
        Email
      </label>
      <label>
        <input type="radio" name="contact" value="none" bind:group={contact} required />
        None
      </label>
    </fieldset>
    <fieldset>
      <legend>Officer</legend>
      <label>
        <input type="radio" name="officer" value="yes" bind:group={officer} required />
        Yes
      </label>
      <label>
        <input type="radio" name="officer" value="no" bind:group={officer} required />
        No
      </label>
    </fieldset>

    <button type="submit">Submit</button>
  </form>
</main>

<style>
  main {
    max-width: 420px;
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
