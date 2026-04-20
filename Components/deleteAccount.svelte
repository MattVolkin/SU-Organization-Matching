<script>
    // This component will handle the delete account functionality. It will be a simple page with a button that will delete the user's account when clicked. We will also add a confirmation popup to prevent accidental deletions. 
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
            window.location.href = '/'; // go back to home page (will most likely prompt you to login/make an account again)
        }
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
    <button on:click={handleDelete}>Delete Information</button>
</div>


<style>

:global(body) {
        margin: 0;
        padding: 0;
        min-height: 100%;
        background: #ffffff;
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
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    button:hover {
        background-color: darkred;
    }
</style>



