<script>
  export let userType = "user";
</script>

<header class="header">
  <div class="header-content">
    <h1>SU Organization Matching Tool</h1>
    <button class="login-button" type="button">Login</button>
  </div>
  <nav class="nav">
    <a href="/">Home</a>
    <a href="/about">About this Project</a>
    <a href="/howto">How To Use this tool</a>
    {#if userType === 'admin'}
      <a href="/create">Create New Club</a>
      <a href="/change-user-type">Change User Type</a>
    {:else if userType === 'officer'}
      <a href="/manage-club">Manage Club</a>
    {/if}
  </nav>
</header>

<style>
  .header {
    background-color: #2c3e50;
    color: white;
    padding: 1rem 2rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  h1 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: 600;
  }

  .login-button {
    padding: 0.6rem 1.5rem;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }

  .login-button:hover {
    background-color: #2980b9;
  }

  .nav {
    display: flex;
    gap: 2rem;
    flex-wrap: wrap;
  }

  .nav a {
    color: white;
    text-decoration: none;
    font-size: 1rem;
    transition: color 0.3s ease;
  }

  .nav a:hover {
    color: #3498db;
  }
</style>


