<!-- @component
 Creates the admin home page 
 
 **TODO**:
  - Replace edit and delete actions with actual API calls -->
<script>
/**
 * @type {statew} adminPreviewType - Allows admin to preview the page as a regular user, officer or admin. Defaults to admin view. FOR TESTING PURPOSES ONLY
 * @type {state} maxClubsPerPage - Number of clubs to show per page in the club management table
 * @type {state} pageNum - Current page number for club management pagination
 * @type {state} clubs - List of clubs to display in the club management table, should be fetched from the API in a real implementation
 * @function getAllClubs - Placeholder function to fetch all clubs from the API and store in state
 * @function nextPage - Increments pageNum to show the next page of clubs, disabled if on the last page 
 * @function prevPage - Decrements pageNum to show the previous page of clubs, disabled if on the first page
 * @function getClubName - Helper function to get the club name from a club object or string, returns 'Unknown Club' if name is not available 
 * @function editClub - Placeholder function to handle editing a club, should be replaced with actual implementation to edit club details
 * @function deleteClub - Placeholder function to handle deleting a club, should be replaced with actual implementation to delete the club from the database and update the UI
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
    const totalPages = $derived(Math.max(1, Math.ceil(clubs.length / maxClubsPerPage)))
    const startIndex = $derived((pageNum - 1) * maxClubsPerPage)
    const paginatedClubs = $derived(clubs.slice(startIndex, startIndex + maxClubsPerPage))

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

    //Todo: add API to fetch all clubs and store in state
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
       await APICreater('DELETE', '/api/admin/orgs', { id: club.id })
       clubs = clubs.filter((existingClub) => existingClub?.id !== club.id)
       if (pageNum > totalPages) {
            pageNum = totalPages
       }
    }

</script>
<div class="admin-home">
    <Header userType="admin" previewAs={adminPreviewType} onPreviewChange={(nextView) => adminPreviewType = nextView} />    
    {#if adminPreviewType === 'admin'}
        <h1>Admin Home</h1>
        <div class="club-management">
            <p>Welcome to the Admin Home! Here you can manage clubs and what they post on the website.</p>
            <!-- Create the club management table  EX CS Club        Edit    Delete-->
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
        <HomePage previewAs={adminPreviewType} showChrome={false} />
    {/if}
    <section class="quiz-quick-action" aria-label="Admin quiz shortcut">
        <div>
            <h2>Want to Take the Quiz?</h2>
        </div>
        <button class="quiz-action" type="button" onclick={goToQuiz}>Take The Quiz</button>
    </section>
</div>
<Footer />

<style>
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

    .quiz-quick-action p {
        margin: 0.3rem 0 0;
        color: var(--muted);
        line-height: 1.35;
        font-size: 0.95rem;
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
