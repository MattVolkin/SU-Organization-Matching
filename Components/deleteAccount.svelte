<script>
    // This component will handle the delete account functionality. It will be a simple page with a button that will delete the user's account when clicked. We will also add a confirmation popup to prevent accidental deletions. 
    let clickedCheckBox = false;

    function handleDelete() {
        if (!clickedCheckBox) {
            alert("Please check the box to confirm you want to delete your account.");
            return;
        }

        const confirmation = confirm("Are you sure you want to delete your account? This action cannot be undone.");
        if (confirmation) {

            //actually do the api call to delete the account here. This is just a placeholder for now.

            
            alert("Your account has been deleted."); 
            window.location.href = '/'; // go back to home page (will most likely prompt you to login/make an account again)
        }
    }

    function toggleCheckBox() {
        clickedCheckBox = !clickedCheckBox;
    }
    


</script>


<div class="delete-container">
    <h1>Delete Account</h1>
    <p>Are you sure you want to delete your account? This action cannot be undone.</p>
    <label>
        <input type="checkbox" on:change={toggleCheckBox} />
        I understand that this action cannot be undone.
    </label>
    <button on:click={handleDelete}>Delete Account</button>
</div>






