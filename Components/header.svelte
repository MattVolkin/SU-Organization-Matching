<script>
 let {userType= "user"} =$props() ;
</script>
<div class="header">
  <h1>SU Organization Matching Tool</h1>
  <div class="login_button">
    <button type="button">Login</button> 
  </div>
  <nav>
    <a href="/">Home</a>
    <a href="/about">About this Project</a>
    <a href="/howto">How To Use this tool</a>
    {#if userType === 'admin'}
      <a href="/create">Create New Club</a>
      <a href="/change-user-type">Change User Type</a> <!-- TODO: Implement change user type functionality -->
    {:else if userType === 'officer'}
      <a href="/manage-club">Manage Club</a> <!-- TODO: Implement manage club functionality -->
    {/if}
    <!-- Add more links as needed -->
  </nav>
</div>
<style>
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background-color: #f8f8f8;
  }

  nav a {
    margin-right: 1rem;
    text-decoration: none;
    color: #333;
  }

  nav a:hover {
    text-decoration: underline;
  }

  .login_button button {
    padding: 0.5rem 1rem;
    font-size: 1rem;
    cursor: pointer;
  }
</style>

