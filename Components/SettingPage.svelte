<script lang="ts">
  import { onMount } from 'svelte';
  import { APICreater } from './APIHandler.svelte';

  type ClubValue = string | {
    ID?: number;
    clubName?: string;
    ClubName?: string;
    officers?: string[];
    personalityTraits?: string[];
    activities?: string[];
  };

  const uploadedImagesStorageKey = 'club-uploaded-images';
  const selectedImageStorageKey = 'club-selected-images';

  let userType = 'user';
  let pageNotice = $state('');
  let accessibleClubs = $state<string[]>(["Test Club 1"]);
  let selectedClubFromUrl = $state('');
  let selectedAdjectives = $state<string[]>([]);
  let allAdjectives = $state<string[]>([]);
  let officerClubs = $state<ClubValue[]>([]);
  let clubImageLibrary = $state<Record<string, string[]>>({});
  let selectedImageByClub = $state<Record<string, string>>({});
  let pendingOfficerEmailByClub = $state<Record<string, string>>({});
  let officerStatusByClub = $state<Record<string, string>>({});
  let loading = $state(true);
  let isTestMode = $state(false);
  let previewUrl = $state('');

  function getOrgApiPath() {
    return userType === 'admin' ? '/api/admin/orgs' : '/api/officer/orgs';
  }

  async function loadUserType() {
    const tokenFromStorage = localStorage.getItem('authToken') || '';
    const headers = tokenFromStorage
      ? { Authorization: `Bearer ${tokenFromStorage}` }
      : {};

    const response = await fetch('/api/user', {
      method: 'GET',
      credentials: 'include',
      headers,
    });

    if (!response.ok) {
      userType = 'user';
      return;
    }

    const data = await response.json().catch(() => ({}));
    const role = String(data?.role || '').toLowerCase();
    userType = role === 'admin' || role === 'officer' ? role : 'user';
  }
  
  const testAdjectives = [
    'Welcoming',
    'Community Service',
    'Study Group',
    'Greek Life',
    'Board Games',
    'Nerdy',
  ];

  function makePreviewImage(label: string, color: string) {
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 720"><defs><linearGradient id="g" x1="0" x2="1" y1="0" y2="1"><stop offset="0%" stop-color="${color}"/><stop offset="100%" stop-color="#0f172a"/></linearGradient></defs><rect width="1200" height="720" fill="url(#g)"/><circle cx="280" cy="280" r="180" fill="rgba(255,255,255,0.15)"/><circle cx="910" cy="470" r="240" fill="rgba(255,255,255,0.12)"/><text x="70" y="640" fill="white" font-family="Arial" font-size="62" font-weight="700">${label}</text></svg>`;
    return `data:image/svg+xml;utf8,${encodeURIComponent(svg)}`;
  }

  function enablePreviewMode(requestedClubName: string, noticeMessage = '') {
    const requestedTestClub = requestedClubName || 'Test Club';
    isTestMode = true;
    pageNotice = noticeMessage;
    selectedClubFromUrl = requestedTestClub;
    accessibleClubs = [requestedTestClub];
    allAdjectives = testAdjectives;
    selectedAdjectives = ['Welcoming', 'Community-focused'];

    if (!(clubImageLibrary[requestedTestClub] || []).length) {
      const previewImages = [
        makePreviewImage(`${requestedTestClub} Activity`, '#1d4ed8'),
        makePreviewImage(`${requestedTestClub} Event`, '#0891b2'),
      ];

      clubImageLibrary = {
        ...clubImageLibrary,
        [requestedTestClub]: previewImages,
      };
      selectedImageByClub = {
        ...selectedImageByClub,
        [requestedTestClub]: previewImages[0],
      };
      saveClubMedia();
    }

    loading = false;
  }

  function getClubName(club: ClubValue) {
    return typeof club === 'string' ? club : (club?.clubName || club?.ClubName || 'Unknown Club');
  }

  function getClubID(club: ClubValue) {
    if (typeof club === 'string') {
      return 0;
    }
    const rawID = (club as { ID?: number }).ID;
    return typeof rawID === 'number' ? rawID : 0;
  }

  function getExistingClubOfficers(club: ClubValue) {
    if (typeof club === 'string') {
      return [];
    }
    return Array.isArray(club.officers) ? club.officers : [];
  }

  function getClubActivities(club: ClubValue) {
    if (typeof club === 'string') {
      return [];
    }
    return Array.isArray(club.activities) ? club.activities : [];
  }

  function toClubSlug(clubName: string) {
    return clubName.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
  }

  function normalizeClubList(response: unknown): ClubValue[] {
    return Array.isArray(response) ? response : [];
  }

  async function loadAdjectives() {
    const response = await APICreater('GET', '/api/adjectives', null);
    allAdjectives = Array.isArray(response)
      ? response
          .map((item) => (typeof item === 'string' ? item : (item?.en || item?.question_type || '')))
          .filter((item) => typeof item === 'string' && item.trim().length > 0)
      : [];
  }

  async function getClubOfficers() {
    const response = await APICreater('GET', getOrgApiPath(), null);
    officerClubs = normalizeClubList(response);
  }

  async function saveClubOfficer(club: string, officerEmail: string) {
    const clubInfo = officerClubs.find((entry) => getClubName(entry) === club);
    const clubID = clubInfo ? getClubID(clubInfo) : 0;
    if (clubID <= 0) {
      throw new Error('Missing valid club ID for officer update');
    }

    const mergedOfficers = [...getExistingClubOfficers(clubInfo || '')];
    if (!mergedOfficers.some((existing) => existing.toLowerCase() === officerEmail.toLowerCase())) {
      mergedOfficers.push(officerEmail);
    }

    await APICreater('PATCH', getOrgApiPath(), {
      id: clubID,
      officers: mergedOfficers,
    });

    officerClubs = officerClubs.map((entry) => {
      if (getClubName(entry) !== club || typeof entry === 'string') {
        return entry;
      }
      return {
        ...entry,
        officers: mergedOfficers,
      };
    });
  }

  function setPendingOfficerEmail(club: string, value: string) {
    pendingOfficerEmailByClub = {
      ...pendingOfficerEmailByClub,
      [club]: value,
    };
  }

  function setOfficerStatus(club: string, message: string) {
    officerStatusByClub = {
      ...officerStatusByClub,
      [club]: message,
    };
  }

  function isValidEmail(email: string) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  }

  async function addOfficer(club: string) {
    const officerEmail = (pendingOfficerEmailByClub[club] || '').trim();
    if (!officerEmail) {
      setOfficerStatus(club, 'Enter an email address first.');
      return;
    }

    if (!isValidEmail(officerEmail)) {
      setOfficerStatus(club, 'Enter a valid email address.');
      return;
    }

    try {
      if (!isTestMode) {
        await saveClubOfficer(club, officerEmail);
      }

      setPendingOfficerEmail(club, '');
      setOfficerStatus(club, `Added ${officerEmail} to ${club}.`);
    } catch (error) {
      console.error('Unable to add officer', error);
      setOfficerStatus(club, 'Unable to add officer right now.');
    }
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
    const clubInfo = officerClubs.find((entry) => getClubName(entry) === club);
    const clubID = clubInfo ? getClubID(clubInfo) : 0;
    if (clubID <= 0) {
      return Promise.reject(new Error('Missing valid club ID for adjective update'));
    }

    return APICreater('PATCH', getOrgApiPath(), {
      id: clubID,
      personalityTraits: adjectives,
      activities: getClubActivities(clubInfo || ''),
    });
  }

  onMount(async () => {
    loadSavedClubMedia();

    const searchParams = new URLSearchParams(window.location.search);
    const previewParams = new URLSearchParams(window.location.search);
    previewParams.set('preview', '1');
    previewUrl = `${window.location.pathname}?${previewParams.toString()}`;
    const requestedClubFromParams = searchParams.get('testClub') || searchParams.get('club') || '';
    const previewModeRequested = searchParams.get('preview') === '1';
    const shouldUseTestMode =
      previewModeRequested ||
      searchParams.get('test') === '1' ||
      requestedClubFromParams.trim().toLowerCase() === 'test club';

    if (shouldUseTestMode) {
      enablePreviewMode(
        requestedClubFromParams,
        previewModeRequested
          ? 'Preview mode enabled. You are viewing mock content for layout testing.'
          : 'Test mode enabled. API checks are bypassed for this page.'
      );
      return;
    }

    try {
      await loadUserType();
      if (userType !== 'admin' && userType !== 'officer') {
        pageNotice = 'Only officers and admins can manage club settings.';
        accessibleClubs = [];
        return;
      }

      await Promise.all([loadAdjectives(), getClubOfficers()]);

      const pathParts = window.location.pathname.split('/').filter(Boolean);
      const slugFromPath = pathParts[0] === 'manage-club' ? (pathParts[1] || '') : '';
      const clubFromQuery = new URLSearchParams(window.location.search).get('club') || '';
      const clubFromSlug = officerClubs.find((club) => toClubSlug(getClubName(club)) === slugFromPath);
      const requestedClub = clubFromQuery || (clubFromSlug ? getClubName(clubFromSlug) : '');

      selectedClubFromUrl = requestedClub;

      if (!requestedClub) {
        enablePreviewMode(
          'Test Club')
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

      const requestedClubInfo = officerClubs.find((club) => getClubName(club) === requestedClub);
      if (requestedClubInfo && typeof requestedClubInfo !== 'string' && Array.isArray(requestedClubInfo.personalityTraits)) {
        selectedAdjectives = requestedClubInfo.personalityTraits;
      }

      pageNotice = '';
      accessibleClubs = [requestedClub];
    } catch (error) {
      console.error('Unable to load club settings', error);
      enablePreviewMode(
        requestedClubFromParams,
        'Unable to load club settings from the API. Showing preview mode instead.'
      );
    } finally {
      if (!isTestMode) {
        loading = false;
      }
    }
  });
</script>

<main class="settings-page">
  <h2>{selectedClubFromUrl ? `${selectedClubFromUrl} Settings` : 'Club Settings'}</h2>
  <p class="warning-banner">All changes on this page are saved immediately after clicking the "Save" button.</p>

  {#if loading}
    <!-- <p class="club-warning">Loading club settings...</p> -->
  {/if}

  {#if pageNotice}
    <!-- <p class="club-warning">{pageNotice}</p> -->
  {/if}

  {#if !isTestMode}
    <p class="preview-help">
      Want to quickly test this layout without API access?
      <a href={previewUrl}>Open preview mode</a>
    </p>
  {/if}

  {#if accessibleClubs.length === 0}
    <p class="club-warning">No editable club loaded.</p>
  {:else}
    {#each accessibleClubs as club}
      <section class="club-settings-card">
        <p class="club-title">{club} Settings</p>
        <p>Change adjectives to describe the club:</p>

        <div class="adjective-row">
          <select multiple bind:value={selectedAdjectives}>
            {#each allAdjectives as adjective}
              <option value={adjective}>{adjective}</option>
            {/each}
          </select>

          <button class="save-button action-button" type="button" onclick={() => saveClubAdjectives(club, selectedAdjectives)}>Save Adjectives</button>
        </div>

        <h3>Add new officer</h3>
        <div class="officer-row">
          <input
            type="email"
            placeholder="southwestern.edu"
            value={pendingOfficerEmailByClub[club] || ''}
            oninput={(event) => setPendingOfficerEmail(club, (event.currentTarget as HTMLInputElement).value)}
          />
          <button type="button" onclick={() => addOfficer(club)}>Add Officer</button>
        </div>
        {#if officerStatusByClub[club]}
          <p class="officer-status">{officerStatusByClub[club]}</p>
        {/if}

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
    color: black;
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

  .preview-help {
    margin: 0.25rem 0 1rem 0;
    padding: 0.7rem 0.9rem;
    border-radius: 0.65rem;
    background: #fff7ed;
    border: 1px solid #fed7aa;
    color: #9a3412;
    font-weight: 600;
  }

  .preview-help a {
    margin-left: 0.35rem;
    color: #c2410c;
  }

  select,
  input[type='email'],
  input[type='file'] {
    width: 100%;
    max-width: 42rem;
    margin-top: 0.35rem;
  }

  input[type='email'] {
    border: 1px solid #c9d6e5;
    border-radius: 0.5rem;
    padding: 0.55rem 0.65rem;
    font-size: 0.98rem;
    background: #fff;
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

  .adjective-row {
    display: flex;
    flex-wrap: wrap;
    align-items: flex-end;
    gap: 0.6rem;
  }

  .adjective-row select {
    flex: 1 1 30rem;
    max-width: 43rem;
  }

  .officer-row {
    display: flex;
    flex-wrap: wrap;
    gap: 0.6rem;
    align-items: center;
  }

  .officer-row button {
    margin-top: 0.35rem;
  }

  .action-button {
    margin-top: 0.35rem;
    min-height: 2.35rem;
    padding: 0.6rem 1.05rem;
    border-radius: 0.6rem;
    font-weight: 700;
  }

  .save-button {
    background: #0f6d8c;
    color: #fff;
    border: none;
  }

  .save-button:hover {
    background: #0b5a74;
  }

  .officer-status {
    margin: 0.5rem 0 0.2rem;
    color: #0f5132;
    font-weight: 600;
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