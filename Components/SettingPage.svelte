<script lang="ts">
  import { onMount } from 'svelte';
  import Header from './header.svelte';
  import Footer from './footer.svelte';
  import { APICreater } from './APIHandler.svelte';

  type ClubValue = string | { clubName?: string; ClubName?: string };

  const uploadedImagesStorageKey = 'club-uploaded-images';
  const selectedImageStorageKey = 'club-selected-images';

  let pageNotice = $state('');
  let accessibleClubs = $state<string[]>([]);
  let selectedClubFromUrl = $state('');
  let selectedAdjectives = $state<string[]>([]);
  let allAdjectives = $state<string[]>([]);
  let officerClubs = $state<ClubValue[]>([]);
  let clubImageLibrary = $state<Record<string, string[]>>({});
  let selectedImageByClub = $state<Record<string, string>>({});
  let loading = $state(true);

  function getClubName(club: ClubValue) {
    return typeof club === 'string' ? club : (club?.clubName || club?.ClubName || 'Unknown Club');
  }

  function toClubSlug(clubName: string) {
    return clubName.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
  }

  function normalizeClubList(response: unknown): ClubValue[] {
    return Array.isArray(response) ? response : [];
  }

  async function loadAdjectives() {
    const response = await APICreater('GET', '/api/admin/adjectives', null);
    allAdjectives = Array.isArray(response) ? response : [];
  }

  async function loadOfficerClubs() {
    const response = await APICreater('GET', '/api/officer/orgs', null);
    officerClubs = normalizeClubList(response);
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

  function saveClubAdjectives(club: string, adjectives: string[]) {
    return APICreater('POST', '/api/officer/update', { ClubName: club, Adjectives: adjectives });
  }

  onMount(async () => {
    loadSavedClubMedia();

    try {
      await Promise.all([loadAdjectives(), loadOfficerClubs()]);

      const pathParts = window.location.pathname.split('/').filter(Boolean);
      const slugFromPath = pathParts[0] === 'manage-club' ? (pathParts[1] || '') : '';
      const clubFromQuery = new URLSearchParams(window.location.search).get('club') || '';
      const clubFromSlug = officerClubs.find((club) => toClubSlug(getClubName(club)) === slugFromPath);
      const requestedClub = clubFromQuery || getClubName(clubFromSlug || '');

      selectedClubFromUrl = requestedClub;

      if (!requestedClub) {
        pageNotice = 'Open a specific club settings page from the Manage Club header menu.';
        accessibleClubs = [];
        return;
      }

      if (slugFromPath && clubFromQuery && toClubSlug(clubFromQuery) !== slugFromPath) {
        pageNotice = 'The selected club URL is invalid. Open the page again from the Manage Club menu.';
        accessibleClubs = [];
        return;
      }

      const canManageRequestedClub = officerClubs.some((club) => getClubName(club) === requestedClub);
      if (!canManageRequestedClub) {
        pageNotice = `You do not have permission to manage "${requestedClub}".`;
        accessibleClubs = [];
        return;
      }

      pageNotice = '';
      accessibleClubs = [requestedClub];
    } catch (error) {
      console.error('Unable to load club settings', error);
      pageNotice = 'Unable to load club settings right now.';
      accessibleClubs = [];
    } finally {
      loading = false;
    }
  });
</script>

<Header userType="officer" />

<main class="settings-page">
  <h2>{selectedClubFromUrl ? `${selectedClubFromUrl} Settings` : 'Club Settings'}</h2>
  <p class="warning-banner">All changes on this page are saved immediately.</p>

  {#if loading}
    <p class="club-warning">Loading club settings...</p>
  {/if}

  {#if pageNotice}
    <p class="club-warning">{pageNotice}</p>
  {/if}

  {#if accessibleClubs.length === 0}
    <p class="club-warning">No editable club loaded.</p>
  {:else}
    {#each accessibleClubs as club}
      <section class="club-settings-card">
        <p class="club-title">{club} Settings</p>
        <p>Change adjectives to describe the club:</p>

        <select multiple bind:value={selectedAdjectives}>
          {#each allAdjectives as adjective}
            <option value={adjective}>{adjective}</option>
          {/each}
        </select>

        <button type="button" onclick={() => saveClubAdjectives(club, selectedAdjectives)}>Save Adjectives</button>

        <h3>Upload club images for results page</h3>
        <input type="file" accept="image/*" multiple onchange={(event) => handleImageUpload(club, event)} />

        {#if (clubImageLibrary[club] || []).length > 0}
          <p>Select one image to use on the results page:</p>
          <div class="image-grid">
            {#each clubImageLibrary[club] as imageUrl}
              <button
                class:selected={selectedImageByClub[club] === imageUrl}
                class="image-choice"
                onclick={() => chooseResultImage(club, imageUrl)}
                type="button"
              >
                <img src={imageUrl} alt={`Uploaded image for ${club}`} />
              </button>
            {/each}
          </div>
        {/if}
      </section>
    {/each}
  {/if}
</main>

<Footer />

<style>
  .settings-page {
    width: min(100%, 1040px);
    margin: 0 auto;
    padding: 1rem;
  }

  .settings-page h2 {
    margin: 0.75rem 0 0.85rem 0;
    font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
    color: #132c45;
  }

  .warning-banner {
    margin: 0 0 1rem 0;
    padding: 0.8rem 0.95rem;
    border-radius: 0.75rem;
    background: #edf7fb;
    border: 1px solid #bfdde8;
    color: #24485e;
    font-weight: 600;
  }

  .club-settings-card {
    border: 1px solid #d7dee8;
    border-radius: 0.85rem;
    padding: 1rem;
    margin-bottom: 1rem;
    background: #fafcff;
    box-shadow: 0 10px 24px rgba(13, 37, 62, 0.08);
  }

  .club-title {
    margin-top: 0;
    font-weight: 700;
  }

  .club-warning {
    border: 1px solid #f4c27a;
    background: #fff8ee;
    color: #7a4d12;
    border-radius: 0.65rem;
    padding: 0.75rem 0.9rem;
    margin: 0.8rem 0;
  }

  select,
  input[type='file'] {
    width: 100%;
    max-width: 42rem;
    margin-top: 0.35rem;
  }

  select {
    min-height: 9rem;
    padding: 0.5rem;
  }

  button {
    margin-top: 0.8rem;
    border: none;
    border-radius: 0.55rem;
    padding: 0.6rem 0.95rem;
    font-weight: 700;
    color: #fff;
    background: #0f6d8c;
    cursor: pointer;
  }

  h3 {
    margin-top: 1rem;
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