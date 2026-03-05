<!-- creates a header for the pages with a login button and navigation links. The links shown depend on the user type (admin, officer, or regular user) -->
<script>
   let {userType = "user"} = $props();
</script>

<header class="header">
  <div class="header-content">
    <h1>SU Organization Matching Tool</h1>
    <!-- TODO :Work on login logic-->
    <button class="login-button" type="button">Login</button>
  </div>
  <!-- TODO: replace hrefs with actual links to pages once they are created -->
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
     background-color: #2c3e50;  /*background header color */
    color: white;
    padding: 1rem 2rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		text-align: center;
    /* width: 100%; */
  }

     .header-content {
    position: relative;
    display: flex;
    justify-content: center;   /* center title container */
    align-items: center;
    margin-bottom: 1rem;
    min-height: 2.75rem;       /* keeps space for button height */
  }

  h1 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: 600;
    text-align: center;
    width: 100%;
  }

  .login-button {
    position: absolute;
    right: 0;
    top: 50%;
    transform: translateY(-50%);
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
    justify-content: space-around;
    align-items: center;
    width: 100%;
    gap: 0;
  }

  .nav a {
    color: white;
    text-decoration: none;
    font-size: 1rem;
    transition: color 0.3s ease;
    flex: 1;
    text-align: center;
    padding: 0.5rem;
  }


  .nav a:hover {
    color: #3498db;
  }
</style>