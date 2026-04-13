<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { APICreater } from './APIHandler.svelte';

  type ClubValue = string | {
    id?: number;
    ID?: number;
    clubName?: string;
    ClubName?: string;
    description?: string;
    Description?: string;
    meetingTime?: string;
    MeetingTime?: string;
    externalLink?: string;
    ExternalLink?: string;
    contactInfo?: string;
    ContactInfo?: string;
    includeOfficerEmails?: boolean;
    IncludeOfficerEmails?: boolean;
    personality?: string[];
    officers?: string[];
    personalityTraits?: string[];
    activities?: string[];
    genders?: string[];
    ethnicities?: string[];
    religions?: string[];
    strict_genders?: boolean;
    dedicated_majors?: string[];
    other?: string[];
  };

  const uploadedImagesStorageKey = 'club-uploaded-images';
  const selectedImageStorageKey = 'club-selected-images';

  let userType = 'user';
  let pageNotice = $state('');
  let accessibleClubs = $state<string[]>(["Test Club 1"]);
  let selectedClubFromUrl = $state('');
  let generalMeetingTime = $state('');
  let generalSocialMedia = $state('');
  let resultContactInfo = $state('');
  let includeOfficerEmailsInResults = $state(false);
  let resultDescription = $state('');
  let resultActivitiesText = $state('');
  let selectedAdjectives = $state<string[]>([]);
  let allAdjectives = $state<string[]>([]);
  let allClubActivities = $state<string[]>([]);
  let combinedTraitOptions = $derived(getTraitSelectOptions());
  let trendsGender = $state('');
  let trendsEthnicities = $state('');
  let trendsReligions = $state('');
  let trendsDedicatedMajors = $state('');
  let trendsOther = $state('');
  let trendsStrictGenders = $state(false);
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
    allClubActivities = csvToList('Board Games, Study Group, Movies, Guest Speakers, Community Service, Leadership');
    selectedAdjectives = ['Welcoming', 'Community-focused'];
    generalMeetingTime = 'Thursdays at 6:30 PM';
    generalSocialMedia = 'https://instagram.com/yourclub';
    resultContactInfo = '';
    includeOfficerEmailsInResults = false;
    resultDescription = 'Add your club description for the results page here.';
    resultActivitiesText = 'Board Games, Study Group, Movies, Guest Speakers, Community Service';
    trendsGender = 'Any';
    trendsEthnicities = 'Any';
    trendsReligions = 'Any';
    trendsDedicatedMajors = 'Computer Science';
    trendsOther = 'No experience required';
    trendsStrictGenders = false;

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
    const rawID = (club as { ID?: number; id?: number }).ID ?? (club as { ID?: number; id?: number }).id;
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

  function getClubPersonalityTraits(club: ClubValue) {
    if (typeof club === 'string') {
      return [];
    }
    if (Array.isArray(club.personality)) {
      return club.personality;
    }
    return Array.isArray(club.personalityTraits) ? club.personalityTraits : [];
  }

  function listToCSV(value: unknown) {
    return Array.isArray(value)
      ? value
          .map((entry) => String(entry || '').trim())
          .filter((entry) => entry.length > 0)
          .join(', ')
      : '';
  }

  function csvToList(value: string) {
    return value
      .split(',')
      .map((entry) => entry.trim())
      .filter((entry) => entry.length > 0);
  }

  function buildContactInfoForSave(club: ClubValue) {
    const baseContact = resultContactInfo.trim();
    if (!includeOfficerEmailsInResults) {
      return baseContact;
    }

    const officerEmails = getExistingClubOfficers(club)
      .map((email) => String(email || '').trim())
      .filter((email) => email.length > 0);

    if (officerEmails.length === 0) {
      return baseContact;
    }

    const existingText = baseContact.toLowerCase();
    const missingEmails = officerEmails.filter((email) => !existingText.includes(email.toLowerCase()));
    if (missingEmails.length === 0) {
      return baseContact;
    }

    const officerLine = `Officer emails: ${missingEmails.join(', ')}`;
    return baseContact ? `${baseContact}\n${officerLine}` : officerLine;
  }

  function getTraitSelectOptions() {
    const seen = new Set<string>();
    const merged: string[] = [];

    for (const item of allAdjectives) {
      const normalized = String(item || '').trim();
      const key = normalized.toLowerCase();
      if (normalized && !seen.has(key)) {
        seen.add(key);
        merged.push(normalized);
      }
    }

    for (const activity of csvToList(resultActivitiesText)) {
      const key = activity.toLowerCase();
      if (!seen.has(key)) {
        seen.add(key);
        merged.push(activity);
      }
    }

    for (const activity of allClubActivities) {
      const normalized = String(activity || '').trim();
      const key = normalized.toLowerCase();
      if (normalized && !seen.has(key)) {
        seen.add(key);
        merged.push(normalized);
      }
    }

    return merged;
  }

  function toClubSlug(clubName: string) {
    return clubName.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
  }

  function normalizeClubList(response: unknown): ClubValue[] {
    return Array.isArray(response) ? response : [];
  }

  function normalizeAdjectiveLabel(value: unknown) {
    if (typeof value !== 'string') {
      return '';
    }

    const normalized = value.trim();
    if (!normalized) {
      return '';
    }

    const lower = normalized.toLowerCase();
    if (lower === 'personality_traits' || lower === 'adjective') {
      return '';
    }

    return normalized;
  }

  function firstStringFromValue(value: unknown) {
    if (typeof value === 'string') {
      return normalizeAdjectiveLabel(value);
    }

    if (Array.isArray(value)) {
      for (const entry of value) {
        if (typeof entry === 'string') {
          const normalized = normalizeAdjectiveLabel(entry);
          if (normalized) {
            return normalized;
          }
        }
      }
    }

    return '';
  }

  function extractAdjectiveLabel(item: unknown) {
    if (typeof item === 'string') {
      return normalizeAdjectiveLabel(item);
    }

    if (!item || typeof item !== 'object') {
      return '';
    }

    const record = item as Record<string, unknown>;
    const directKeys = ['label', 'name', 'value', 'text', 'en'];

    for (const key of directKeys) {
      const directValue = firstStringFromValue(record[key]);
      if (directValue) {
        return directValue;
      }
    }

    const translations = record.translations;
    if (translations && typeof translations === 'object' && !Array.isArray(translations)) {
      const translationMap = translations as Record<string, unknown>;
      const preferredLocales = ['en', 'en-US', 'english'];

      for (const locale of preferredLocales) {
        const preferredValue = firstStringFromValue(translationMap[locale]);
        if (preferredValue) {
          return preferredValue;
        }
      }

      for (const value of Object.values(translationMap)) {
        const fallbackValue = firstStringFromValue(value);
        if (fallbackValue) {
          return fallbackValue;
        }
      }
    }

    return '';
  }

  async function createTraitsFromAdjectivesAPI() {
    const response = await APICreater('GET', '/api/adjectives', null);
    const source = Array.isArray(response)
      ? response
      : (Array.isArray((response as { adjectives?: unknown[] } | null)?.adjectives)
        ? (response as { adjectives: unknown[] }).adjectives
        : []);

    const seen = new Set<string>();
    const labels: string[] = [];

    for (const item of source) {
      const label = extractAdjectiveLabel(item);
      const dedupeKey = label.toLowerCase();
      if (!label || seen.has(dedupeKey)) {
        continue;
      }
      seen.add(dedupeKey);
      labels.push(label);
    }

    allAdjectives = labels;
  }

  async function loadAdjectives() {
    await createTraitsFromAdjectivesAPI();
  }

  async function getClubOfficers() {
    const response = await APICreater('GET', getOrgApiPath(), null);
    officerClubs = normalizeClubList(response);
  }

  async function loadAllClubActivities() {
    const response = await APICreater('GET', '/api/results', null);
    const clubs = Array.isArray(response) ? response : [];
    const seen = new Set<string>();
    const merged: string[] = [];

    for (const club of clubs) {
      const activities = Array.isArray(club?.activities) ? club.activities : [];
      for (const activity of activities) {
        const normalized = String(activity || '').trim();
        const key = normalized.toLowerCase();
        if (normalized && !seen.has(key)) {
          seen.add(key);
          merged.push(normalized);
        }
      }
    }

    allClubActivities = merged;
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
      description: resultDescription.trim(),
      meetingTime: generalMeetingTime.trim(),
      externalLink: generalSocialMedia.trim(),
      contactInfo: buildContactInfoForSave(clubInfo || ''),
      includeOfficerEmails: includeOfficerEmailsInResults,
      personality: adjectives,
      activities: csvToList(resultActivitiesText),
      genders: csvToList(trendsGender),
      ethnicities: csvToList(trendsEthnicities),
      religions: csvToList(trendsReligions),
      strict_genders: trendsStrictGenders,
      dedicated_majors: csvToList(trendsDedicatedMajors),
      other: csvToList(trendsOther),
    });
  }

  async function initializeSettingsPage() {
    loadSavedClubMedia();
    loading = true;

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

      await Promise.all([createTraitsFromAdjectivesAPI(), getClubOfficers(), loadAllClubActivities()]);

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
      selectedAdjectives = getClubPersonalityTraits(requestedClubInfo || '');
      generalMeetingTime = requestedClubInfo && typeof requestedClubInfo !== 'string'
        ? String(requestedClubInfo.meetingTime ?? requestedClubInfo.MeetingTime ?? '').trim()
        : '';
      generalSocialMedia = requestedClubInfo && typeof requestedClubInfo !== 'string'
        ? String(requestedClubInfo.externalLink ?? requestedClubInfo.ExternalLink ?? '').trim()
        : '';
      resultContactInfo = requestedClubInfo && typeof requestedClubInfo !== 'string'
        ? String(requestedClubInfo.contactInfo ?? requestedClubInfo.ContactInfo ?? '').trim()
        : '';
      includeOfficerEmailsInResults = requestedClubInfo && typeof requestedClubInfo !== 'string'
        ? Boolean(requestedClubInfo.includeOfficerEmails ?? requestedClubInfo.IncludeOfficerEmails)
        : false;
      resultDescription = requestedClubInfo && typeof requestedClubInfo !== 'string'
        ? String(requestedClubInfo.description ?? requestedClubInfo.Description ?? '').trim()
        : '';
      resultActivitiesText = listToCSV(getClubActivities(requestedClubInfo || ''));

      if (requestedClubInfo && typeof requestedClubInfo !== 'string') {
        trendsGender = listToCSV(requestedClubInfo.genders);
        trendsEthnicities = listToCSV(requestedClubInfo.ethnicities);
        trendsReligions = listToCSV(requestedClubInfo.religions);
        trendsDedicatedMajors = listToCSV(requestedClubInfo.dedicated_majors);
        trendsOther = listToCSV(requestedClubInfo.other);
        trendsStrictGenders = Boolean(requestedClubInfo.strict_genders);
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
  }

  function handleAuthLogin() {
    if (isTestMode) {
      return;
    }
    void initializeSettingsPage();
  }

  function handleAuthLogout() {
    if (isTestMode) {
      return;
    }
    void initializeSettingsPage();
  }

  onMount(() => {
    window.addEventListener('auth-login', handleAuthLogin);
    window.addEventListener('auth-logout', handleAuthLogout);
    void initializeSettingsPage();
  });

  onDestroy(() => {
    window.removeEventListener('auth-login', handleAuthLogin);
    window.removeEventListener('auth-logout', handleAuthLogout);
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
  {#if accessibleClubs.length === 0}
    <p class="club-warning">No editable club loaded.</p>
  {:else}
    <section class="club-settings-card">

      <h3>Results page content</h3>
      <div class="results-content-grid">
        <label>
          Description
          <textarea rows="4" bind:value={resultDescription}></textarea>
        </label>

        <label>
          Activities (comma-separated)
          <textarea rows="3" bind:value={resultActivitiesText}></textarea>
        </label>

        <div class="general-info-grid">
          <label>
            Meeting time
            <input type="text" bind:value={generalMeetingTime} />
          </label>

          <label>
            Social media / website
            <input type="text" bind:value={generalSocialMedia} placeholder="https://instagram.com/yourclub" />
          </label>
        </div>

        <label class="checkbox-label">
          <input type="checkbox" bind:checked={includeOfficerEmailsInResults} />
          Include officer emails in contact info
        </label>

        <div>
          <h3>Upload club images for results page</h3>
          <input type="file" accept="image/*" multiple onchange={(event) => handleImageUpload(selectedClubFromUrl, event)} />

          {#if (clubImageLibrary[selectedClubFromUrl] || []).length > 0}
            <p>Select one image to use on the results page:</p>
            <div class="image-grid">
              {#each clubImageLibrary[selectedClubFromUrl] as imageUrl}
                <button
                  class:selected={selectedImageByClub[selectedClubFromUrl] === imageUrl}
                  class="image-choice"
                  onclick={() => chooseResultImage(selectedClubFromUrl, imageUrl)}
                  type="button"
                >
                  <img src={imageUrl} alt={`Uploaded image for ${selectedClubFromUrl}`} />
                </button>
              {/each}
            </div>
          {/if}
        </div>
      </div>

      <h3>Personality + activities trait select</h3>
      <div class="adjective-row">
        <div class="traits-column">
          <label for="personality-select"></label>
          <select id="personality-select" multiple bind:value={selectedAdjectives}>
            {#each combinedTraitOptions as traitOption}
              <option value={traitOption}>{traitOption}</option>
            {/each}
          </select>
        </div>
      </div>

      <button class="save-button action-button" type="button" onclick={() => saveClubAdjectives(selectedClubFromUrl, selectedAdjectives)}>Save Adjectives</button>

      <h3>Trends</h3>
      <div class="trends-grid">
        <label>
          Genders (comma-separated)
          <input type="text" bind:value={trendsGender} />
        </label>

        <label>
          Ethnicities (comma-separated)
          <input type="text" bind:value={trendsEthnicities} />
        </label>

        <label>
          Religions (comma-separated)
          <input type="text" bind:value={trendsReligions} />
        </label>

        <label>
          Dedicated majors (comma-separated)
          <input type="text" bind:value={trendsDedicatedMajors} />
        </label>

        <label class="trends-full-row">
          Other (comma-separated)
          <input type="text" bind:value={trendsOther} />
        </label>

        <label class="checkbox-label trends-full-row">
          <input type="checkbox" bind:checked={trendsStrictGenders} />
          Strict genders matching
        </label>
      </div>

      <h3>Add new officer</h3>
      <div class="officer-row">
        <input
          type="email"
          placeholder="southwestern.edu"
          value={pendingOfficerEmailByClub[selectedClubFromUrl] || ''}
          oninput={(event) => setPendingOfficerEmail(selectedClubFromUrl, (event.currentTarget as HTMLInputElement).value)}
        />
        <button type="button" onclick={() => addOfficer(selectedClubFromUrl)}>Add Officer</button>
      </div>
      {#if officerStatusByClub[selectedClubFromUrl]}
        <p class="officer-status">{officerStatusByClub[selectedClubFromUrl]}</p>
      {/if}
    </section>
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
    border: none;
    color: #24485e;
    font-weight: 600;
  }

  .club-settings-card {
    border: none;
    border-radius: 0;
    padding: 1rem;
    margin-bottom: 1rem;
    background: transparent;
    box-shadow: none;
    color: black;
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
  textarea,
  input[type='text'],
  input[type='email'],
  input[type='file'] {
    width: 100%;
    max-width: 42rem;
    margin-top: 0.35rem;
  }

  input[type='email'],
  input[type='text'],
  textarea {
    border: 1px solid #c9d6e5;
    border-radius: 0.5rem;
    padding: 0.55rem 0.65rem;
    font-size: 0.98rem;
    background: #fff;
  }

  textarea {
    resize: vertical;
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
    flex: 1 1 18rem;
    max-width: 43rem;
  }

  .traits-column {
    flex: 1 1 18rem;
    display: grid;
    gap: 0.35rem;
  }

  .traits-column label {
    font-weight: 700;
  }

  .results-content-grid {
    display: grid;
    gap: 0.7rem;
    margin-bottom: 0.9rem;
  }

  .general-info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 0.65rem;
    margin-bottom: 0.9rem;
  }

  .general-info-grid label {
    display: grid;
    gap: 0.35rem;
    font-weight: 700;
  }

  .results-content-grid label {
    display: grid;
    gap: 0.35rem;
    font-weight: 700;
  }

  .trends-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0.8rem;
  }

  .trends-grid label {
    display: grid;
    gap: 0.35rem;
    font-weight: 700;
    min-width: 0;
  }

  .trends-grid input[type='text'] {
    width: min(75%, 460px);
    max-width: 460px;
  }

  .trends-full-row input[type='text'] {
    width: min(75%, 520px);
    max-width: 520px;
  }

  .trends-full-row {
    grid-column: 1 / -1;
  }

  .checkbox-label {
    display: flex !important;
    align-items: center;
    gap: 0.45rem;
  }

  .checkbox-label input[type='checkbox'] {
    width: 1.15rem;
    height: 1.15rem;
    margin: 0;
    accent-color: #0f6d8c;
    flex: 0 0 auto;
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

  @media (max-width: 860px) {
    .trends-grid {
      grid-template-columns: 1fr;
    }

    .trends-grid input[type='text'],
    .trends-full-row input[type='text'] {
      width: 100%;
      max-width: none;
    }
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