<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { APICreater } from './APIHandler.svelte';
  import MultiSelectDropdown from './MultiSelectDropdown.svelte';

  type ClubValue = string | {
    id?: number;
    ID?: number;
    clubName?: string;
    ClubName?: string;
    imagePath?: string;
    ImagePath?: string;
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
    activitiesDescreption?: string;
    activitiesDescription?: string;
    ActivitiesDescreption?: string;
    ActivitiesDescription?: string;
    genders?: string[];
    ethnicities?: string[];
    religions?: string[];
    strict_genders?: boolean;
    dedicated_majors?: string[];
    associated_majors?: string[];
    other?: string[];
  };

  const uploadedImagesStorageKey = 'club-uploaded-images';
  const selectedImageStorageKey = 'club-selected-images';

  let userType = 'user';
  let pageNotice = $state('');
  let accessibleClubs = $state<string[]>([]);
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

  // Trend options
  const genderOptions = ['Man', 'Woman', 'Non-Binary', 'Other'];
 const raceOptions = [
    'American Native/Alaska Native',
    'Asian',
    'Black or African American',
    'Hispanic or Latino',
    'Middle Eastern or North African',
    'Native Hawaiian or Pacific Islander',
    'White']
  const religionOptions = ['Buddhism', 'Catholicism', 'Hinduism', 'Islam', 'Judaism', 'Protestantism', 'No Religion'];
  const majorOptions = [
    'Anthropology',
    'Applied Physics',
    'Art (Studio)',
    'Art History',
    'Biochemistry',
    'Biology',
    'Business',
    'Chemistry',
    'Classics',
    'Communication Studies',
    'Computational Mathematics',
    'Computer Science',
    'Pre-Dentistry',
    'Economics',
    'Education',
    'Pre-Engineering',
    'English',
    'Environmental Studies',
    'Feminist Studies',
    'French',
    'German',
    'Greek',
    'Health Professions',
    'History',
    'International Studies',
    'Kinesiology',
    'Latin',
    'Latin American & Border Studies',
    'Pre-Law',
    'Mathematics',
    'Pre-Med',
    'Pre-Ministry',
    'Music',
    'Pre-Nursing',
    'Pre-Occupational Therapy',
    'Philosophy',
    'Physics',
    'Political Science',
    'Psychology',
    'Pre-Physician Assistant',
    'Pre-Physical Therapy',
    'Religion',
    'Sociology',
    'Spanish',
    'Theatre',
    'Undecided',
  ]


  const otherOptions = ['LGBTQ', 'Greek Life', "Disabilities"]

  let trendsGender = $state<string[]>([]);
  let trendsEthnicities = $state<string[]>([]);
  let trendsReligions = $state<string[]>([]);
  let trendsDedicatedMajors = $state<string[]>([]);
  let trendsAssociatedMajors = $state<string[]>([]);
  let trendsOther = $state<string[]>([]);
  let trendsStrictGenders = $state(false);
  let officerClubs = $state<ClubValue[]>([]);
  let clubImageLibrary = $state<Record<string, string[]>>({});
  let selectedImageByClub = $state<Record<string, string>>({});
  let pendingOfficerEmailByClub = $state<Record<string, string>>({});
  let selectedOfficerEmailsByClub = $state<Record<string, string[]>>({});
  let officerStatusByClub = $state<Record<string, string>>({});
  let loading = $state(true);
  let saveStateResults = $state<'idle' | 'saving' | 'saved' | 'error'>('idle');
  let saveStateAdjectives = $state<'idle' | 'saving' | 'saved' | 'error'>('idle');
  let saveStateTrends = $state<'idle' | 'saving' | 'saved' | 'error'>('idle');
  let saveStateResetTimers = {
    results: 0,
    adjectives: 0,
    trends: 0,
  };

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

  function makePreviewImage(label: string, color: string) {
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 720"><defs><linearGradient id="g" x1="0" x2="1" y1="0" y2="1"><stop offset="0%" stop-color="${color}"/><stop offset="100%" stop-color="#0f172a"/></linearGradient></defs><rect width="1200" height="720" fill="url(#g)"/><circle cx="280" cy="280" r="180" fill="rgba(255,255,255,0.15)"/><circle cx="910" cy="470" r="240" fill="rgba(255,255,255,0.12)"/><text x="70" y="640" fill="white" font-family="Arial" font-size="62" font-weight="700">${label}</text></svg>`;
    return `data:image/svg+xml;utf8,${encodeURIComponent(svg)}`;
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

  function getAuthHeaders() {
    const tokenFromStorage = localStorage.getItem('authToken') || '';
    return tokenFromStorage
      ? { Authorization: `Bearer ${tokenFromStorage}` }
      : {};
  }

  function buildClubImageApiPath(club: ClubValue) {
    const clubID = getClubID(club);
    if (clubID <= 0) {
      throw new Error('Missing valid club ID for club image request');
    }

    return `${getOrgApiPath()}/${clubID}/image`;
  }

  async function parseClubImageApiResponse(response: Response) {
    const payload = await response.json().catch(() => ({}));
    if (!response.ok) {
      const message = typeof payload?.error === 'string' && payload.error.trim()
        ? payload.error.trim()
        : 'Club image request failed';
      throw new Error(message);
    }

    return payload;
  }

  async function uploadClubImageFile(club: ClubValue, imageFile: File) {
    const requestBody = new FormData();
    requestBody.set('image', imageFile);

    const response = await fetch(buildClubImageApiPath(club), {
      method: 'POST',
      credentials: 'include',
      headers: getAuthHeaders(),
      body: requestBody,
    });

    return parseClubImageApiResponse(response);
  }

  function extractUploadedImagePath(payload: unknown) {
    if (!payload || typeof payload !== 'object') {
      return '';
    }

    const record = payload as Record<string, unknown>;
    const topLevelPath = typeof record.imagePath === 'string' ? record.imagePath.trim() : '';
    if (topLevelPath) {
      return topLevelPath;
    }

    const club = record.club;
    if (club && typeof club === 'object') {
      const clubRecord = club as Record<string, unknown>;
      const nestedPath = typeof clubRecord.imagePath === 'string' ? clubRecord.imagePath.trim() : '';
      if (nestedPath) {
        return nestedPath;
      }
    }

    return '';
  }

  function getExistingClubOfficers(club: ClubValue) {
    if (typeof club === 'string') {
      return [];
    }
    return Array.isArray(club.officers) ? club.officers : [];
  }

  function getOfficersForClub(clubName: string) {
    const club = officerClubs.find((entry) => getClubName(entry) === clubName);
    return getExistingClubOfficers(club || '');
  }

  function getSelectedOfficerEmails(clubName: string) {
    return selectedOfficerEmailsByClub[clubName] || [];
  }

  function setSelectedOfficerEmails(clubName: string, emails: string[]) {
    selectedOfficerEmailsByClub = {
      ...selectedOfficerEmailsByClub,
      [clubName]: emails,
    };
  }

  function getClubActivities(club: ClubValue) {
    if (typeof club === 'string') {
      return [];
    }
    return Array.isArray(club.activities) ? club.activities : [];
  }

  function getClubActivitiesText(club: ClubValue) {
    if (typeof club === 'string') {
      return '';
    }

    const description = String(
      club.activitiesDescreption ??
      club.activitiesDescription ??
      club.ActivitiesDescreption ??
      club.ActivitiesDescription ??
      ''
    ).trim();

    if (description) {
      return description;
    }

    return listToCSV(getClubActivities(club));
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

  function mergeUniqueStrings(...groups: string[][]) {
    const seen = new Set<string>();
    const merged: string[] = [];

    for (const group of groups) {
      for (const item of group) {
        const normalized = String(item || '').trim();
        const key = normalized.toLowerCase();
        if (normalized && !seen.has(key)) {
          seen.add(key);
          merged.push(normalized);
        }
      }
    }

    return merged;
  }

  function getTraitSelectOptions() {
    const seen = new Set<string>();
    const merged: string[] = [];

    const selectedClub = officerClubs.find((club) => getClubName(club) === selectedClubFromUrl);
    const selectedClubActivities = selectedClub ? getClubActivities(selectedClub) : [];

    for (const item of allAdjectives) {
      const normalized = String(item || '').trim();
      const key = normalized.toLowerCase();
      if (normalized && !seen.has(key)) {
        seen.add(key);
        merged.push(normalized);
      }
    }

    for (const activity of selectedClubActivities) {
      const normalized = String(activity || '').trim();
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
      await saveClubOfficer(club, officerEmail);

      setPendingOfficerEmail(club, '');
      setOfficerStatus(club, `Added ${officerEmail} to ${club}.`);
    } catch (error) {
      console.error('Unable to add officer', error);
      setOfficerStatus(club, 'Unable to add officer right now.');
    }
  }

  async function deleteSelectedOfficers(club: string) {
    const clubInfo = officerClubs.find((entry) => getClubName(entry) === club);
    const clubID = clubInfo ? getClubID(clubInfo) : 0;
    if (clubID <= 0) {
      setOfficerStatus(club, 'Missing valid club ID for officer deletion.');
      return;
    }

    const selectedEmails = getSelectedOfficerEmails(club)
      .map((email) => String(email || '').trim())
      .filter((email) => email.length > 0);

    if (selectedEmails.length === 0) {
      setOfficerStatus(club, 'Select at least one officer to delete.');
      return;
    }

    const selectedLookup = new Set(selectedEmails.map((email) => email.toLowerCase()));
    const remainingOfficers = getExistingClubOfficers(clubInfo || '')
      .map((email) => String(email || '').trim())
      .filter((email) => email.length > 0)
      .filter((email) => !selectedLookup.has(email.toLowerCase()));

    const confirmation = confirm(
      selectedEmails.length === 1
        ? ' Do you really want to delete this officer? This action cannot be undone.'
        : ' Do you really want to delete these officers? This action cannot be undone.'
    );

    if (!confirmation) {
      return;
    }

    try {
      await APICreater('PATCH', getOrgApiPath(), {
        id: clubID,
        officers: remainingOfficers,
      });

      officerClubs = officerClubs.map((entry) => {
        if (getClubName(entry) !== club || typeof entry === 'string') {
          return entry;
        }

        return {
          ...entry,
          officers: remainingOfficers,
        };
      });

      setSelectedOfficerEmails(club, []);
      setOfficerStatus(club, selectedEmails.length === 1 ? 'Deleted 1 officer.' : `Deleted ${selectedEmails.length} officers.`);
    } catch (error) {
      console.error('Unable to delete officers', error);
      setOfficerStatus(club, 'Unable to delete officers right now.');
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

    const clubInfo = officerClubs.find((entry) => getClubName(entry) === club);
    if (!clubInfo) {
      input.value = '';
      return;
    }

    try {
      const uploadedPaths: string[] = [];
      for (const file of Array.from(input.files)) {
        const responsePayload = await uploadClubImageFile(clubInfo, file);
        const uploadedPath = extractUploadedImagePath(responsePayload);
        if (uploadedPath) {
          uploadedPaths.push(uploadedPath);
        }
      }

      if (uploadedPaths.length === 0) {
        throw new Error('Upload completed without a returned image path.');
      }

      const latestImage = uploadedPaths[uploadedPaths.length - 1];

      clubImageLibrary = {
        ...clubImageLibrary,
        [club]: [latestImage],
      };

      selectedImageByClub = {
        ...selectedImageByClub,
        [club]: latestImage,
      };

      officerClubs = officerClubs.map((entry) => {
        if (typeof entry === 'string' || getClubName(entry) !== club) {
          return entry;
        }

        return {
          ...entry,
          imagePath: latestImage,
          ImagePath: latestImage,
        };
      });

      saveClubMedia();
      input.value = '';
      return;
    } catch (error) {
      console.error('Unable to upload club image', error);
      input.value = '';
      return;
    }
  }

  function chooseResultImage(club: string, imageUrl: string) {
    selectedImageByClub = {
      ...selectedImageByClub,
      [club]: imageUrl,
    };
    saveClubMedia();
  }

  function clearSaveStateResetTimer(section: 'results' | 'adjectives' | 'trends') {
    const activeTimer = saveStateResetTimers[section];
    if (activeTimer) {
      window.clearTimeout(activeTimer);
      saveStateResetTimers[section] = 0;
    }
  }

  function scheduleSaveStateReset(section: 'results' | 'adjectives' | 'trends') {
    clearSaveStateResetTimer(section);
    saveStateResetTimers[section] = window.setTimeout(() => {
      if (section === 'results') {
        saveStateResults = 'idle';
      }
      if (section === 'adjectives') {
        saveStateAdjectives = 'idle';
      }
      if (section === 'trends') {
        saveStateTrends = 'idle';
      }
      saveStateResetTimers[section] = 0;
    }, 2200);
  }

  function getSaveIndicatorLabel(state: 'idle' | 'saving' | 'saved' | 'error') {
    if (state === 'saving') {
      return 'Saving...';
    }
    if (state === 'saved') {
      return 'Saved';
    }
    if (state === 'error') {
      return 'Unable to save. Try again.';
    }
    return '';
  }

  async function handleSaveResultsPageContent() {
    saveStateResults = 'saving';
    clearSaveStateResetTimer('results');

    try {
      await saveResultsPageInfo(
        resultDescription,
        resultActivitiesText,
        generalMeetingTime,
        generalSocialMedia,
        resultContactInfo.trim(),
        includeOfficerEmailsInResults
      );
      saveStateResults = 'saved';
      scheduleSaveStateReset('results');
    } catch (error) {
      console.error('Unable to save results page content', error);
      saveStateResults = 'error';
    }
  }

  async function handleSaveAdjectives() {
    saveStateAdjectives = 'saving';
    clearSaveStateResetTimer('adjectives');

    try {
      await saveClubAdjectives(selectedClubFromUrl, selectedAdjectives);
      saveStateAdjectives = 'saved';
      scheduleSaveStateReset('adjectives');
    } catch (error) {
      console.error('Unable to save traits', error);
      saveStateAdjectives = 'error';
    }
  }

  async function handleSaveTrends() {
    saveStateTrends = 'saving';
    clearSaveStateResetTimer('trends');

    try {
      await saveTrends(
        trendsGender,
        trendsEthnicities,
        trendsReligions,
        trendsDedicatedMajors,
        trendsAssociatedMajors,
        trendsOther,
        trendsStrictGenders
      );
      saveStateTrends = 'saved';
      scheduleSaveStateReset('trends');
    } catch (error) {
      console.error('Unable to save trends', error);
      saveStateTrends = 'error';
    }
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
      contactInfo: resultContactInfo.trim(),
      includeOfficerEmails: includeOfficerEmailsInResults,
      personality: adjectives,
      activities: csvToList(resultActivitiesText),
      genders: trendsGender,
      ethnicities: trendsEthnicities,
      religions: trendsReligions,
      strict_genders: trendsStrictGenders,
      dedicated_majors: trendsDedicatedMajors,
      associated_majors: trendsAssociatedMajors,
      other: trendsOther,
    });
  }

  async function initializeSettingsPage() {
    loadSavedClubMedia();
    loading = true;

    const searchParams = new URLSearchParams(window.location.search);
    const requestedClubFromParams = searchParams.get('club') || '';

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
        pageNotice = 'No club was selected.';
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

      const requestedClubInfo = officerClubs.find((club) => getClubName(club) === requestedClub);
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
        if(APICreater)
      resultActivitiesText = getClubActivitiesText(requestedClubInfo || '');
      selectedAdjectives = mergeUniqueStrings(
        getClubPersonalityTraits(requestedClubInfo || ''),
        getClubActivities(requestedClubInfo || ''),
        csvToList(resultActivitiesText)
      );
      setSelectedOfficerEmails(requestedClub, []);

      if (requestedClubInfo && typeof requestedClubInfo !== 'string') {
        trendsGender = requestedClubInfo.genders || [];
        trendsEthnicities = requestedClubInfo.ethnicities || [];
        trendsReligions = requestedClubInfo.religions || [];
        trendsDedicatedMajors = requestedClubInfo.dedicated_majors || [];
        trendsAssociatedMajors = requestedClubInfo.associated_majors || [];
        trendsOther = requestedClubInfo.other || [];
        trendsStrictGenders = Boolean(requestedClubInfo.strict_genders);

        const existingImagePath = String(requestedClubInfo.imagePath ?? requestedClubInfo.ImagePath ?? '').trim();
        if (existingImagePath) {
          clubImageLibrary = {
            ...clubImageLibrary,
            [requestedClub]: [existingImagePath],
          };

          selectedImageByClub = {
            ...selectedImageByClub,
            [requestedClub]: existingImagePath,
          };

          saveClubMedia();
        }
      }

      pageNotice = '';
      accessibleClubs = [requestedClub];
    } catch (error) {
      console.error('Unable to load club settings', error);
      pageNotice = 'Unable to load club settings from the API.';
      accessibleClubs = [];
    } finally {
      loading = false;
    }
  }

  function handleAuthLogin() {
    void initializeSettingsPage();
  }

  function handleAuthLogout() {
    void initializeSettingsPage();
  }
  async function saveTrends(
    genders: string[],
    ethnicities: string[],
    religions: string[],
    dedicatedMajors: string[],
    associatedMajors: string[],
    other: string[],
    strictGenders: boolean
  ) {
    const clubInfo = officerClubs.find((entry) => getClubName(entry) === selectedClubFromUrl);
    const clubID = clubInfo ? getClubID(clubInfo) : 0;
    if (clubID <= 0) {
      return Promise.reject(new Error('Missing valid club ID for trends update'));
    }

    return APICreater('PATCH', getOrgApiPath(), {
      id: clubID,
      genders: genders,
      ethnicities: ethnicities,
      religions: religions,
      dedicated_majors: dedicatedMajors,
      associated_majors: associatedMajors,
      other: other,
      strict_genders: strictGenders
    });

  }
  async function saveResultsPageInfo(description: string, activities: string, meetingTime: string, socialMedia: string, contactInfo: string, includeOfficerEmails: boolean) {
    const clubInfo = officerClubs.find((entry) => getClubName(entry) === selectedClubFromUrl);
    const clubID = clubInfo ? getClubID(clubInfo) : 0;
    if (clubID <= 0) {
      return Promise.reject(new Error('Missing valid club ID for results page update'));
    }

    return APICreater('PATCH', getOrgApiPath(), {
      id: clubID,
      description: description.trim(),
      meetingTime: meetingTime.trim(),
      externalLink: socialMedia.trim(),
      contactInfo: contactInfo.trim(),
      includeOfficerEmails: includeOfficerEmails,
      activitiesDescription: activities,
    });
  }

  onMount(() => {
    window.addEventListener('auth-login', handleAuthLogin);
    window.addEventListener('auth-logout', handleAuthLogout);
    void initializeSettingsPage();
  });

  onDestroy(() => {
    window.removeEventListener('auth-login', handleAuthLogin);
    window.removeEventListener('auth-logout', handleAuthLogout);
    clearSaveStateResetTimer('results');
    clearSaveStateResetTimer('adjectives');
    clearSaveStateResetTimer('trends');
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
          Activities Description
          <textarea rows="3" bind:value={resultActivitiesText}></textarea>
        </label>

        <div class="general-info-grid">
          <label>
            Meeting Information
            <input type="text" bind:value={generalMeetingTime} />
          </label>

          <label>
            Social Media / Website
            <input type="text" bind:value={generalSocialMedia} placeholder="https://instagram.com/yourclub" />
          </label>
        </div>

        <label class="checkbox-label">
          <input type="checkbox" bind:checked={includeOfficerEmailsInResults} />
          Include Officer Emails in contact info
        </label>

        <div>
          <h3>Upload Club Images for the Results Page</h3>
          <input type="file" accept="image/*" onchange={(event) => handleImageUpload(selectedClubFromUrl, event)} />

          {#if (clubImageLibrary[selectedClubFromUrl] || []).length > 0}
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

      <button
        class="save-button action-button"
        type="button"
        onclick={() => handleSaveResultsPageContent()}
        disabled={saveStateResults === 'saving'}
      >
        Save Results Page Content
      </button>
      {#if saveStateResults !== 'idle'}
        <p class={`save-status ${saveStateResults}`} aria-live="polite">{getSaveIndicatorLabel(saveStateResults)}</p>
      {/if}

      <h3>Personality Traits & Activities </h3>
      <div class="adjective-row">
        <MultiSelectDropdown id="personality-traits" label="Select traits" options={combinedTraitOptions} bind:value={selectedAdjectives} />
      </div>

      <button class="save-button action-button" type="button" onclick={() => handleSaveAdjectives()} disabled={saveStateAdjectives === 'saving'}>Save Traits</button>
      {#if saveStateAdjectives !== 'idle'}
        <p class={`save-status ${saveStateAdjectives}`} aria-live="polite">{getSaveIndicatorLabel(saveStateAdjectives)}</p>
      {/if}

      <h3>Trends</h3>
      <div class="trends-grid">
        <MultiSelectDropdown id="trends-genders" label="Genders" options={genderOptions} bind:value={trendsGender} />
        <MultiSelectDropdown id="trends-ethnicities" label="Ethnicities" options={raceOptions} bind:value={trendsEthnicities} />
        <MultiSelectDropdown id="trends-religions" label="Religions" options={religionOptions} bind:value={trendsReligions} />
        <MultiSelectDropdown id="trends-majors" label="Dedicated Majors" options={majorOptions} bind:value={trendsDedicatedMajors} />
        <MultiSelectDropdown id="trends-associated-majors" label="Associated Majors" options={majorOptions} bind:value={trendsAssociatedMajors} />
        <div class="trends-full-row">
          <MultiSelectDropdown id="trends-other" label="Other" options={otherOptions} bind:value={trendsOther} />
        </div>

        <label class="checkbox-label trends-full-row">
          <input type="checkbox" bind:checked={trendsStrictGenders} />
          Strict Genders Matching
        </label>
      </div>

      <button
        class="save-button action-button"
        type="button"
        onclick={() => handleSaveTrends()}
        disabled={saveStateTrends === 'saving'}
      >
        Save Trends
      </button>
      {#if saveStateTrends !== 'idle'}
        <p class={`save-status ${saveStateTrends}`} aria-live="polite">{getSaveIndicatorLabel(saveStateTrends)}</p>
      {/if}

      <h3>Current officers</h3>
      <div class="officers-box">
        {#if getOfficersForClub(selectedClubFromUrl).length > 0}
          <div class="officer-list">
            {#each getOfficersForClub(selectedClubFromUrl) as officerEmail}
              <label class="officer-item">
                <input
                  type="checkbox"
                  checked={getSelectedOfficerEmails(selectedClubFromUrl).some((selectedEmail) => selectedEmail.toLowerCase() === String(officerEmail).trim().toLowerCase())}
                  onchange={(event) => {
                    const currentEmail = String(officerEmail).trim();
                    const isChecked = (event.currentTarget as HTMLInputElement).checked;
                    const currentSelection = getSelectedOfficerEmails(selectedClubFromUrl);

                    if (isChecked) {
                      if (!currentSelection.some((selectedEmail) => selectedEmail.toLowerCase() === currentEmail.toLowerCase())) {
                        setSelectedOfficerEmails(selectedClubFromUrl, [...currentSelection, currentEmail]);
                      }
                      return;
                    }

                    setSelectedOfficerEmails(
                      selectedClubFromUrl,
                      currentSelection.filter((selectedEmail) => selectedEmail.toLowerCase() !== currentEmail.toLowerCase())
                    );
                  }}
                />
                <span>{officerEmail}</span>
              </label>
            {/each}
          </div>

          <button type="button" class="delete-button action-button" onclick={() => deleteSelectedOfficers(selectedClubFromUrl)}>
            Delete selected officers
          </button>
        {:else}
          <p class="officer-empty">No officers have been added yet.</p>
        {/if}
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
  :global(html),
  :global(body) {
    margin: 0;
    padding: 0;
    background: #ffffff;
  }

  main {
    background: #ffffff;
    color: #000000;
    font-family: system-ui, sans-serif;
  }

  .settings-page {
    --surface: #ffffff;
    --surface-muted: #f3f3f3;
    --text-main: #000000;
    --text-subtle: #4a4a4a;
    --border: #828282;
    --action: #000000;
    --action-hover: #828282;
    --accent: #ffcd00;

    width: min(100%, 1040px);
    margin: 0 auto;
    padding: 1rem;
    color: var(--text-main);
  }

  .settings-page h2 {
    margin: 0.75rem 0 0.85rem 0;
    font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
    color: var(--text-main);
  }

  .warning-banner {
    margin: 0 0 1rem 0;
    padding: 0.65rem 0.9rem;
    border-radius: 0.5rem;
    background: var(--surface);
    border: 1px solid var(--border);
    color: var(--text-subtle);
    font-size: 0.92rem;
    font-weight: 500;
    box-shadow: none;
  }

  .club-settings-card {
    border: none;
    border-radius: 0;
    padding: 1rem;
    margin-bottom: 1rem;
    background: var(--surface);
    box-shadow: none;
    color: var(--text-main);
  }

  .club-warning {
    border: 1px solid var(--border);
    background: var(--surface-muted);
    color: var(--text-main);
    border-radius: 0.65rem;
    padding: 0.75rem 0.9rem;
    margin: 0.8rem 0;
  }

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
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    padding: 0.55rem 0.65rem;
    font-size: 0.98rem;
    background: var(--surface);
    color: var(--text-main);
  }

  textarea {
    resize: vertical;
  }

  button {
    margin-top: 0.8rem;
    border: 1px solid transparent;
    border-radius: 0.55rem;
    padding: 0.6rem 0.95rem;
    font-weight: 700;
    color: #fff;
    background: var(--action);
    cursor: pointer;
  }

  button:hover {
    background: var(--action-hover);
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
    min-width: 0;
  }

  .general-info-grid input[type='text'] {
    width: 75%;
    max-width: 75%;
    min-width: 0;
  }

  .results-content-grid label {
    display: grid;
    gap: 0.35rem;
    font-weight: 700;
    min-width: 0;
  }

  .trends-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1.2rem;
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
    accent-color: var(--action);
    flex: 0 0 auto;
  }

  .officer-row {
    display: flex;
    flex-wrap: wrap;
    gap: 0.6rem;
    align-items: center;
  }

  .officers-box {
    margin: 0.75rem 0 0.5rem 0;
    padding: 0.9rem;
    border: 1px solid var(--border);
    border-radius: 0.75rem;
    background: var(--surface-muted);
  }

  .officer-list {
    display: grid;
    gap: 0.55rem;
  }

  .officer-item {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    color: var(--text-main);
    font-weight: 600;
  }

  .officer-item input[type='checkbox'] {
    width: 1.15rem;
    height: 1.15rem;
    margin: 0;
    accent-color: var(--action);
    flex: 0 0 auto;
  }

  .officer-empty {
    margin: 0;
    color: var(--text-subtle);
  }

  .delete-button {
    background: #b42318;
    color: #ffffff;
  }

  .delete-button:hover {
    background: #922018;
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
    background: var(--action);
    color: #fff;
    border: 1px solid transparent;
  }

  .save-button:hover {
    background: var(--action-hover);
  }

  .save-button:disabled {
    opacity: 0.7;
    cursor: wait;
  }

  .save-status {
    margin: 0.25rem 0 0.35rem;
    font-size: 0.9rem;
    font-weight: 600;
  }

  .save-status.saving {
    color: var(--text-subtle);
  }

  .save-status.saved {
    color: var(--text-main);
  }

  .save-status.error {
    color: var(--text-main);
  }

  @media (max-width: 860px) {
    .trends-grid {
      grid-template-columns: 1fr;
    }
  }

  .officer-status {
    margin: 0.5rem 0 0.2rem;
    color: var(--text-main);
    font-weight: 600;
  }

  .image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 0.75rem;
    margin-top: 0.5rem;
  }

  .image-choice {
    border: 2px solid var(--border);
    border-radius: 0.6rem;
    padding: 0.2rem;
    background: var(--surface);
    cursor: pointer;
  }

  .image-choice.selected {
    border-color: var(--action);
    box-shadow: 0 0 0 2px rgba(0, 0, 0, 0.2);
  }

  .image-choice img {
    width: 100%;
    height: 100px;
    object-fit: cover;
    border-radius: 0.4rem;
    display: block;
  }

  @media (prefers-color-scheme: dark) {
    :global(html),
    :global(body),
    main {
      background: #000000;
      color: #ffffff;
    }

    .settings-page {
      --surface: #121212;
      --surface-muted: #1e1e1e;
      --text-main: #ffffff;
      --text-subtle: #d5d5d5;
      --border: #828282;
      --action: #ffcd00;
      --action-hover: #e5b800;
      --accent: #ffcd00;
    }

    .club-settings-card,
    .club-warning,
    .warning-banner,
    .officers-box,
    input[type='email'],
    input[type='text'],
    textarea,
    .image-choice {
      background: var(--surface);
      color: var(--text-main);
      border-color: var(--border);
    }

    .warning-banner {
      color: var(--text-subtle);
      box-shadow: none;
    }

    .save-button,
    .delete-button,
    .action-button,
    button {
      color: #000000;
    }

    .settings-page h2,
    h3,
    label,
    .hint,
    .officer-item,
    .officer-empty,
    .officer-status,
    .save-status.saving,
    .save-status.saved,
    .save-status.error {
      color: var(--text-main);
    }

    .image-choice.selected {
      box-shadow: 0 0 0 2px rgba(255, 205, 0, 0.28);
    }

    .delete-button {
      color: #ffffff;
    }

  }

</style>