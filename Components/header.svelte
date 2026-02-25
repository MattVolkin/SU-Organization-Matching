<script>
export let userType = $state("");
</script>

<header>
  <h1>PLACEHOLDER	<!-- TODO: Replace with actual title --></h1>
  <nav>
    <a href="/">Home</a>
    <a href="/about">About this Project</a>
    <a href="/howto">How To Use this tool</a>
    {#if userType === 'admin'}
      <a href="/create">Create New Club</a>
      <a href="/change-user-type">Change User Type</a> <!-- TODO: Implement change user type functionality -->
    {:else if userType === 'officer'}
      <a href="/manage-club">Manage Club</a> <!-- TODO: Implement manage club functionality -->>
    {/if}
    <!-- Add more links as needed -->
  </nav>
</header>
