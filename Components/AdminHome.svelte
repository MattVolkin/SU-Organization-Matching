<!-- @component Creates the admin home page and lets admins manage clubs or preview non-admin views. -->
<script>
/**
 * @type {state} adminPreviewType - Allows admin to preview the page as a regular user, officer, or admin. Defaults to admin view.
 * @type {state} maxClubsPerPage - Number of clubs to show per page in the club management table
 * @type {state} pageNum - Current page number for club management pagination
 * @type {state} clubs - List of clubs fetched from the API for the club management table
 * @function loadClubs - Fetches all clubs from the admin API and updates loading/error state
 * @function nextPage - Increments pageNum to show the next page of clubs, disabled if on the last page 
 * @function prevPage - Decrements pageNum to show the previous page of clubs, disabled if on the first page
 * @function getClubName - Helper function to get the club name from a club object or string, returns 'Unknown Club' if name is not available 
 * @function editClub - Redirects to the club settings page for the selected club
 * @function deleteClub - Deletes a club from the backend and updates the current table view
 */
    import { onMount, onDestroy } from 'svelte'
    import Header from './header.svelte'
    import Footer from './footer.svelte'
    import { APICreater } from './APIHandler.svelte'
    import HomePage from './pages/HomePage.svelte'

    let adminPreviewType = $state('admin')
    let maxClubsPerPage = $state(8)
    let pageNum = $state(1)
    let clubs = $state([])
    let isLoading = $state(true)
    let loadError = $state('')
    let adminEmail = $state('')
    let adminRoleUpdateState = $state('idle')
    let adminRoleStatus = $state('')
    const totalPages = $derived(Math.max(1, Math.ceil(clubs.length / maxClubsPerPage)))
    const startIndex = $derived((pageNum - 1) * maxClubsPerPage)
    const sortedClubs = $derived(
        [...clubs].sort((a, b) => getClubName(a).localeCompare(getClubName(b), undefined, { sensitivity: 'base' }))
    )
    const paginatedClubs = $derived(sortedClubs.slice(startIndex, startIndex + maxClubsPerPage))

    function handleAuthLogin() {
        if (adminPreviewType === 'admin') {
            void loadClubs()
        }
    }

    function handleAuthLogout() {
        clubs = []
        loadError = ''
    }

    onMount(() => {
        window.addEventListener('auth-login', handleAuthLogin)
        window.addEventListener('auth-logout', handleAuthLogout)
        void loadClubs()
    })

    onDestroy(() => {
        window.removeEventListener('auth-login', handleAuthLogin)
        window.removeEventListener('auth-logout', handleAuthLogout)
    })

    // Loads current clubs for admin management mode.
    async function loadClubs() {
        isLoading = true
        loadError = ''
        try {
            const fetchedClubs = await APICreater('GET', '/api/admin/orgs', null)
            clubs = Array.isArray(fetchedClubs) ? fetchedClubs : []
            if (pageNum > totalPages) {
                pageNum = totalPages
            }
        } catch (error) {
            loadError = 'Unable to load clubs. Please refresh and try again.'
            clubs = []
        } finally {
            isLoading = false
        }
    }
    function nextPage() {
        if (pageNum < totalPages) {
            pageNum += 1
        }
    }
    function prevPage() {
        if (pageNum > 1) {
            pageNum -= 1
        }
    }

    function getClubName(club) {
        return typeof club === 'string' ? club : (club?.clubName || club?.name || 'Unknown Club')
    }

    function editClub(club) {
        window.location.href = `/settings.html?club=${encodeURIComponent(getClubName(club))}`
    }

    function goToQuiz() {
        window.location.href = '/demographic-quiz.html'
    }

    async function deleteClub(club) {
       if (!club?.id) {
            return
       }

          const clubName = getClubName(club)
          const confirmed = window.confirm(`Delete "${clubName}"? This action cannot be undone.`)
          if (!confirmed) {
              return
          }

       await APICreater('DELETE', '/api/admin/orgs', { id: club.id })
       clubs = clubs.filter((existingClub) => existingClub?.id !== club.id)
       if (pageNum > totalPages) {
            pageNum = totalPages
       }
    }

    async function updateAdminRole(nextRole) {
        const trimmedEmail = String(adminEmail || '').trim().toLowerCase()
        if (!trimmedEmail) {
            adminRoleUpdateState = 'error'
            adminRoleStatus = 'Enter an email address first.'
            return
        }

        adminRoleUpdateState = 'saving'
        adminRoleStatus = ''

        try {
            const response = await APICreater('PATCH', '/api/admin/users', {
                email: trimmedEmail,
                role: nextRole,
            })

            if (response && typeof response === 'object' && response.error) {
                throw new Error(String(response.error))
            }

            adminRoleUpdateState = 'saved'
            adminRoleStatus = nextRole === 'admin'
                ? `${trimmedEmail} is now an admin.`
                : `${trimmedEmail} is now a member.`
        } catch (error) {
            adminRoleUpdateState = 'error'
            adminRoleStatus = `Unable to update role. ${error instanceof Error ? error.message : 'Please try again.'}`
        }
    }

</script>
<!-- Admin page shell with role-aware preview mode -->
<div class="admin-home">
    <!-- Header supports admin preview switching -->
    <Header userType="admin" previewAs={adminPreviewType} onPreviewChange={(nextView) => adminPreviewType = nextView} />    
    <!-- Admin management view vs user/officer preview rendering -->
    {#if adminPreviewType === 'admin'}
        <h1>Admin Home</h1>
        <div class="club-management">
            <p>Welcome to the Admin Home! Here you can manage clubs and what they post on the website.</p>
            <!-- Club management table with loading, error, and empty states -->
            <table>
                <colgroup>
                    <col class="col-name" />
                    <col class="col-actions" />
                </colgroup>
                <thead>
                    <tr>
                        <th>Club Name</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#if isLoading}
                        <tr>
                            <td colspan="2">Loading clubs...</td>
                        </tr>
                    {:else if loadError}
                        <tr>
                            <td colspan="2">{loadError}</td>
                        </tr>
                    {:else if paginatedClubs.length === 0}
                        <tr>
                            <td colspan="2">No clubs found.</td>
                        </tr>
                    {:else}
                        {#each paginatedClubs as club}
                            <tr>
                                <td>{getClubName(club)}</td>
                                <td>
                                    <button onclick={() => editClub(club)}>Edit</button>
                                    <button onclick={() => deleteClub(club)}>Delete</button>
                               </td>
                            </tr>
                        {/each}
                    {/if}
                </tbody>
            </table>
            <div class="pager">
                <button onclick={prevPage} disabled={pageNum === 1}>Previous</button>
                <span>Page {pageNum}</span>
                <button onclick={nextPage} disabled={pageNum === totalPages}>Next</button>
            </div>
        </div>
    {:else}
        <!-- Reuses the standard homepage in preview-only mode -->
        <HomePage previewAs={adminPreviewType} showChrome={false} />
    {/if}
    {#if adminPreviewType === 'admin'}
        <section class="admin-users" aria-label="Admin user management">
            <h2>Admin Access</h2>
            <p>Add or remove admin access by email.</p>
            <div class="admin-users-controls">
                <label for="admin-email">User email</label>
                <input
                    id="admin-email"
                    type="email"
                    placeholder="name@example.edu"
                    bind:value={adminEmail}
                />
                <div class="admin-users-actions">
                    <button
                        type="button"
                        onclick={() => updateAdminRole('admin')}
                        disabled={adminRoleUpdateState === 'saving'}
                    >
                        Add Admin
                    </button>
                    <button
                        type="button"
                        class="secondary"
                        onclick={() => updateAdminRole('member')}
                        disabled={adminRoleUpdateState === 'saving'}
                    >
                        Remove Admin
                    </button>
                </div>
            </div>
            {#if adminRoleStatus}
                <p class={`admin-users-status ${adminRoleUpdateState}`}>{adminRoleStatus}</p>
            {/if}
        </section>
    {/if}
    {#if adminPreviewType === 'admin'}
        <!-- Quick action to jump directly to the quiz flow -->
        <section class="quiz-quick-action" aria-label="Admin quiz shortcut">
            <div>
                <h2>Want to Take the Quiz?</h2>
            </div>
            <button class="quiz-action" type="button" onclick={goToQuiz}>Take The Quiz</button>
        </section>
    {/if}
</div>
<!-- Global footer -->
<Footer />

<style>
    /* Admin page layout, table, and theme styles */
    .admin-home {
        --bg: linear-gradient(180deg, #edf4fb 0%, #f6f9fd 100%);
        --card: #ffffff;
        --text: #132c45;
        --muted: #4f6781;
        --border: #d4e0ec;
        --accent: #0f6d8c;
        --accent-hover: #0b5972;
        --danger: #a5294a;
        --danger-hover: #8a1f3d;

        background: var(--bg);
        color: var(--text);
        min-height: 100vh;
        padding-bottom: 0.75rem;
    }

    .admin-home h1 {
        width: min(100%, 1040px);
        margin: 1rem auto 0.6rem auto;
        padding: 0 1rem;
        font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
        line-height: 1.2;
    }

    .club-management {
        width: min(100%, 1040px);
        margin: 0 auto 1rem auto;
        padding: 1rem;
        border: none;
        border-radius: 1rem;
        background: var(--card);
        box-shadow: 0 10px 24px rgba(13, 37, 62, 0.1);
    }

    .club-management > p {
        margin: 0 0 0.8rem 0;
        color: var(--muted);
        line-height: 1.45;
    }

    .admin-users {
        width: min(100%, 1040px);
        border: 1px solid var(--border);
        border-radius: 0.8rem;
        padding: 0.9rem;
        margin: 0 auto 0.95rem auto;
        background: #f9fcff;
    }

    .admin-users h2 {
        margin: 0;
        font-size: 1rem;
    }

    .admin-users p {
        margin: 0.4rem 0 0.7rem 0;
        color: var(--muted);
    }

    .admin-users-controls {
        display: grid;
        gap: 0.45rem;
    }

    .admin-users-controls label {
        font-weight: 700;
        font-size: 0.9rem;
    }

    .admin-users-controls input {
        width: min(100%, 28rem);
        border: 1px solid var(--border);
        border-radius: 0.55rem;
        padding: 0.56rem 0.68rem;
        font-size: 0.94rem;
    }

    .admin-users-actions {
        display: flex;
        flex-wrap: wrap;
        gap: 0.6rem;
    }

    .admin-users-actions .secondary {
        background: #4f6781;
    }

    .admin-users-actions .secondary:hover {
        background: #3f556d;
    }

    .admin-users-status {
        margin: 0.7rem 0 0 0;
        font-size: 0.9rem;
        font-weight: 700;
    }

    .admin-users-status.saved {
        color: #1f6f3f;
    }

    .admin-users-status.error {
        color: #9d1f3d;
    }

    .quiz-quick-action {
        width: min(100%, 1040px);
        margin: 0 auto 1rem auto;
        padding: 1rem;
        border: 1px solid var(--border);
        border-radius: 1rem;
        background: linear-gradient(135deg, rgba(15, 109, 140, 0.08), rgba(47, 74, 102, 0.08));
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.9rem;
    }

    .quiz-quick-action h2 {
        margin: 0;
        font-size: 1.05rem;
        color: var(--text);
    }

    .quiz-action {
        min-width: 10.5rem;
        padding: 0.62rem 1.05rem;
        border-radius: 0.6rem;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        border: 1px solid var(--border);
        border-radius: 0.8rem;
        overflow: hidden;
        background: #fff;
    }

    th,
    td {
        text-align: left;
        padding: 0.85rem 1.1rem;
        border-bottom: 1px solid #e8eef4;
    }

    .col-name {
        width: 68%;
    }

    .col-actions {
        width: 32%;
    }

    th {
        background: #f5f9ff;
        text-transform: uppercase;
        letter-spacing: 0.03em;
        font-size: 0.8rem;
        color: #2f4a66;
    }

    th:nth-child(2) {
        text-align: right;
    }

    tbody tr:hover {
        background: #f9fcff;
    }

    td:nth-child(2) {
        display: flex;
        gap: 0.65rem;
        flex-wrap: wrap;
        justify-content: flex-end;
    }

    button {
        border: none;
        border-radius: 0.55rem;
        padding: 0.42rem 0.78rem;
        font-size: 0.84rem;
        font-weight: 700;
        color: #fff;
        background: var(--accent);
        cursor: pointer;
        transition: background-color 0.2s ease;
    }

    td button:last-child {
        background: var(--danger);
    }

    button:hover {
        background: var(--accent-hover);
    }

    td button:last-child:hover {
        background: var(--danger-hover);
    }

    button:disabled {
        opacity: 0.55;
        cursor: not-allowed;
    }

    .pager {
        margin-top: 0.9rem;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.6rem;
    }

    .pager span {
        color: var(--muted);
        font-weight: 600;
    }

    @media (max-width: 760px) {
        .admin-home h1,
        .club-management,
        .quiz-quick-action {
            width: calc(100% - 1rem);
        }

        .club-management {
            padding: 0.8rem;
            border-radius: 0.85rem;
        }

        table,
        thead,
        tbody,
        th,
        td,
        tr {
            display: block;
            width: 100%;
        }

        thead {
            display: none;
        }

        tr {
            border-bottom: 1px solid #e6edf5;
            padding: 0.45rem 0;
        }

        td {
            border: none;
            padding: 0.35rem 0.2rem;
        }

        td:nth-child(2) {
            justify-content: flex-start;
        }

        .pager {
            flex-direction: column;
            align-items: stretch;
        }

        .pager button {
            width: 100%;
        }

        .quiz-quick-action {
            flex-direction: column;
            align-items: stretch;
        }

        .quiz-action {
            width: 100%;
        }
    }

    @media (prefers-color-scheme: dark) {
        .admin-home {
            --bg: linear-gradient(180deg, #0c1725 0%, #0a111b 100%);
            --card: #0f1c2d;
            --text: #deebfb;
            --muted: #b2c6df;
            --border: #2a3b53;
            --accent: #2b8fb5;
            --accent-hover: #3ca3cb;
            --danger: #b84267;
            --danger-hover: #d05079;
        }

        .club-management {
            box-shadow: 0 12px 28px rgba(0, 0, 0, 0.45);
        }

        table {
            background: #102133;
            border-color: #2a3b53;
        }

        th {
            background: #13273e;
            color: #c9d9ef;
        }

        td,
        tr {
            border-bottom-color: #253952;
        }

        tbody tr:hover {
            background: #14263d;
        }
    }
</style>
