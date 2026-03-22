<script lang="ts">
function getAllAdjectives() {
  // Placeholder function - replace with API call to fetch adjectives from the backend to help with the scalibilty
  return ["Creative", "Friendly", "Organized", "Passionate", "Innovative", "Collaborative", "Supportive", "Ambitious", "Inclusive", "Dynamic"];}
function getAllClubsByOfficer(user){
    // Placeholder function - replace with API call to fetch clubs that the given user is an officer of, this will help with scalability as well as making sure that officers can only edit the clubs they are in charge of
    return ["Computer Science Club"];
}
function saveClubAdjectives(club, adjectives) {
    // Placeholder function - replace with API call to save the given adjectives for the given club in the backend database
    console.log(`Saving adjectives for ${club}: ${adjectives.join(", ")}`);
}
function saveClubContent(club, contentFiles) {
    // Placeholder function - replace with API call to save the given content files for the given club in the backend database
    console.log(`Saving content for ${club}: ${contentFiles.join(", ")}`);
}
function saveNewOfficers(club, officersEmails) {
    // Placeholder function - replace with API call to save the given officers for the given club in the backend database
    console.log(`Saving officers for ${club}: ${officersEmails.join(", ")}`);
}
const allAdjectives = getAllAdjectives();
let selectedAdjectives = [];

</script>   

{#each getAllClubsByOfficer("currentUser") as club}
    <p>{club} Settings</p>
    <p>Change adjectives to describe the club:</p>
  <select multiple bind:value={selectedAdjectives}>
    {#each allAdjectives as adjective}
      <option value={adjective}>{adjective}</option>
    {/each}
  </select>
  <button on:click={() => saveClubAdjectives(club, selectedAdjectives)}>Save Adjectives</button>

  <p>Selected: {selectedAdjectives.join(", ") || "None"}</p>

{/each}