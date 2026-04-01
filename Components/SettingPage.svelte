<script lang="ts">
import { onMount } from 'svelte';
import Header from './header.svelte';
import Footer from './footer.svelte';
import APIHandler from './APIHandler.svelte';

const api = new APIHandler();

function getAllAdjectives() {
  // Placeholder function - replace with API call to fetch adjectives from the backend to help with the scalibilty
  return ["Creative", "Friendly", "Organized", "Passionate", "Innovative", "Collaborative", "Supportive", "Ambitious", "Inclusive", "Dynamic"];}
function getAllClubsByOfficer(user){
    // Placeholder function - replace with API call to fetch clubs that the given user is an officer of, this will help with scalability as well as making sure that officers can only edit the clubs they are in charge of
  return ["Club 1", "Club 2"];
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
const officerClubs = getAllClubsByOfficer("currentUser");
let selectedClubFromUrl = '';
let accessibleClubs: string[] = [];
let pageNotice = '';
let selectedAdjectives = [];
let clubImageLibrary: Record<string, string[]> = {};
let selectedImageByClub: Record<string, string> = {};

const uploadedImagesStorageKey = 'clubUploadedImagesByClub';
const selectedImageStorageKey = 'selectedClubResultImages';

function toClubSlug(clubName: string) {
  return clubName.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
}

function loadSavedClubMedia() {
  try {
    const savedUploadedImages = localStorage.getItem(uploadedImagesStorageKey);
    const savedSelectedImages = localStorage.getItem(selectedImageStorageKey);

    clubImageLibrary = savedUploadedImages ? JSON.parse(savedUploadedImages) : {};
    selectedImageByClub = savedSelectedImages ? JSON.parse(savedSelectedImages) : {};
  } catch (error) {
    console.error('Unable to load saved club images', error);
    clubImageLibrary = {};
    selectedImageByClub = {};
  }
}

function saveClubMedia() {
  localStorage.setItem(uploadedImagesStorageKey, JSON.stringify(clubImageLibrary));
  localStorage.setItem(selectedImageStorageKey, JSON.stringify(selectedImageByClub));
}

function toDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result));
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
}

async function handleImageUpload(club: string, event: Event) {
  const input = event.currentTarget as HTMLInputElement;

  if (!input.files || input.files.length === 0) {
    return;
  }

  const encodedImages = await Promise.all(Array.from(input.files).map(toDataUrl));
  const existingImages = clubImageLibrary[club] || [];

  clubImageLibrary = {
    ...clubImageLibrary,
    [club]: [...existingImages, ...encodedImages],
  };

  if (!selectedImageByClub[club] && encodedImages.length > 0) {
    selectedImageByClub = {
      ...selectedImageByClub,
      [club]: encodedImages[0],
    };
  }

  saveClubMedia();
  input.value = '';
}

function chooseResultImage(club: string, imageUrl: string) {
  selectedImageByClub = {
    ...selectedImageByClub,
    [club]: imageUrl,
  };
  saveClubMedia();
}

onMount(() => {
  loadSavedClubMedia();

  const pathParts = window.location.pathname.split('/').filter(Boolean);
  const slugFromPath = pathParts[0] === 'manage-club' ? (pathParts[1] || '') : '';
  const clubFromQuery = new URLSearchParams(window.location.search).get('club') || '';
  const clubFromSlug = officerClubs.find((club) => toClubSlug(club) === slugFromPath) || '';
  const requestedClub = clubFromQuery || clubFromSlug;
  selectedClubFromUrl = requestedClub;

  if (!slugFromPath) {
    pageNotice = 'Open a specific club settings page from the Manage Club header menu.';
    accessibleClubs = [];
    return;
  }

  if (clubFromQuery && toClubSlug(clubFromQuery) !== slugFromPath) {
    pageNotice = 'The selected club URL is invalid. Open the page again from the Manage Club menu.';
    accessibleClubs = [];
    return;
  }

  if (!requestedClub) {
    pageNotice = 'That club settings page does not exist for your account.';
    accessibleClubs = [];
    return;
  }

  const canManageRequestedClub = officerClubs.includes(requestedClub);
  if (!canManageRequestedClub) {
    pageNotice = `You do not have permission to manage "${requestedClub}".`;
    accessibleClubs = [];
    return;
  }

  pageNotice = '';
  accessibleClubs = [requestedClub];
});

</script>   
<h2> WARNING all of these changes will be saved immediately </h2>
{#if pageNotice}
  <p class="club-warning">{pageNotice}</p>
{/if}

{#if accessibleClubs.length === 0}
  <p class="club-warning">No editable club loaded.</p>
{/if}

{#each accessibleClubs as club}
    <section class="club-settings-card">
    <p>{club} Settings</p>
    <p>Change adjectives to describe the club:</p>
  <select multiple bind:value={selectedAdjectives}>
    {#each allAdjectives as adjective}
      <option value={adjective}>{adjective}</option>
    {/each}
  </select>
  <button on:click={() => saveClubAdjectives(club, selectedAdjectives)}>Save Adjectives</button>

  <h3>Upload club images for results page</h3>
  <input type="file" accept="image/*" multiple on:change={(event) => handleImageUpload(club, event)} />

  {#if (clubImageLibrary[club] || []).length > 0}
    <p>Select one image to use on the results page:</p>
    <div class="image-grid">
      {#each clubImageLibrary[club] as imageUrl}
        <button
          class:selected={selectedImageByClub[club] === imageUrl}
          class="image-choice"
          on:click={() => chooseResultImage(club, imageUrl)}
          type="button"
        >
          <img src={imageUrl} alt={`Uploaded image for ${club}`} />
        </button>
      {/each}
    </div>
  {/if}
  </section>



{/each}

<style>
  .club-settings-card {
    border: 1px solid #d7dee8;
    border-radius: 0.75rem;
    padding: 1rem;
    margin-bottom: 1rem;
    background: #fafcff;
  }

  .club-warning {
    border: 1px solid #f4c27a;
    background: #fff8ee;
    color: #7a4d12;
    border-radius: 0.65rem;
    padding: 0.75rem 0.9rem;
    margin: 0.8rem 0;
  }

  .image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 0.75rem;
    margin-top: 0.5rem;
  }

  .image-choice {
    border: 2px solid #d1d5db;
    border-radius: 0.6rem;
    padding: 0.2rem;
    background: #fff;
    cursor: pointer;
  }

  .image-choice.selected {
    border-color: #1f6f8b;
    box-shadow: 0 0 0 2px rgba(31, 111, 139, 0.2);
  }

  .image-choice img {
    width: 100%;
    height: 100px;
    object-fit: cover;
    border-radius: 0.4rem;
    display: block;
  }
</style>