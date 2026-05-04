<!-- @component Provides a delete account interface with confirmation checkbox and dual confirmation dialogs for account removal. -->
<script>
/**
 * @type {state} clickedCheckBox - Boolean to track whether user has confirmed they understand account deletion is permanent
 * @function handleDelete - Handles account deletion request after confirmation, calls logout on success
 * @function logout - Clears local auth state, removes token from localStorage, and dispatches logout event
 * @function toggleCheckBox - Toggles the confirmation checkbox state
 */
    let clickedCheckBox = false;

    import API, { APICreater } from './APIHandler.svelte' // import the API handler to make calls to the backend

  

    async function handleDelete() {
        if (!clickedCheckBox) {
            alert("Please check the box to confirm you want to delete your account.");
            return;
        }

        const confirmation = confirm("Are you sure you want to delete your account? This action cannot be undone.");
        if (confirmation) {

            try {
				await APICreater('POST', '/api/delete', {});
				console.log('tried to delete account information', payload.length);
			} catch (error) {
				console.error('Unable to send delete request', error);
			}
			

            
            alert("Your account has been deleted."); 

            logout(); // log the user out after deleting their account

        }
    }

    async function logout() {
        console.log('Logging out user');
    const headers = authToken
      ? { Authorization: `Bearer ${authToken}` }
      : {};
    await fetch('/logout', {
      method: 'POST',
      credentials: 'include',
      headers,
    });
    userEmail = '';
    authToken = '';
    resolvedUserType = 'user';
    officerClubs = [];
    localStorage.removeItem('authToken');
    window.dispatchEvent(new CustomEvent('auth-logout'));
  }

    function toggleCheckBox() {
        clickedCheckBox = !clickedCheckBox;
    }
    


</script>


<div class="delete-container">
    <h1>Delete stored information</h1>
    <p>Are you sure you want to delete your information? This action cannot be undone.</p>
    <label>
        <input type="checkbox" on:change={toggleCheckBox} />
        I understand that this action cannot be undone.
    </label>
    <button on:click={logout}>Delete Information</button>
</div>


<style>

@media (prefers-color-scheme: dark) {
		:root {

	--in-darkmode: 255;
	--text-color: white;
	--background-color: #1a1a1a;

	}

}

@media (prefers-color-scheme: light) {
		:root {

	--in-darkmode: 0;
	--text-color: #1a1a1a;
	--background-color: white;

	}

}

:global(body) {
        margin: 0;
        padding: 0;
        min-height: 100%;
        background: var(--background-color);
        color: var(--text-color);
    }

    .delete-container {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100vh;
        text-align: center;
    }

    button {
        margin: auto;
        margin-top: 20px;
        padding: 10px 20px;
        background-color: red;
        color: var(--text-color);
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    button:hover {
        background-color: darkred;
    }
</style>



