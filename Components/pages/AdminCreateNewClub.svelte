<script>
/**
 * @type {state} clubName - Name entered for the new club
 * @type {state} PresidentEmail - Email entered for the first club officer
 * @type {state} statusMessage - Confirmation message shown after successful submission
 * @function createClub - Submits the new club to the admin API, then clears the form and shows a status message
 */
    import AdminSwitch from '../AdminSwitch.svelte';
    import Header from '../header.svelte';
    import Footer from '../footer.svelte';
    import { APICreater } from '../APIHandler.svelte';

    let clubName = $state('');
    let PresidentEmail = $state('');
    let statusMessage = $state('');

    async function createClub(event) {
        event.preventDefault();
        await APICreater('POST', '/api/admin/orgs', {
            clubName,
            officers: [PresidentEmail],
        });

        statusMessage = `Submitted ${clubName} for creation.`;
        clubName = '';
        PresidentEmail = '';
    }
</script>

<!-- Top navigation -->
<Header userType="admin" />

<!-- Main admin layout: create-club form + view switch panel -->
<main class="create-club-page">
    <div class="backdrop-shape shape-left"></div>
    <div class="backdrop-shape shape-right"></div>

    <!-- Club creation form card -->
    <section class="create-club-card" aria-labelledby="create-club-title">
        <div class="card-head">
            <p class="eyebrow">Admin Controls</p>
            <h1 id="create-club-title">Create a New Club</h1>
            <p class="subtitle">Add a club and assign its first officer in one step.</p>
        </div>

        <form class="create-form" onsubmit={createClub}>
            <label for="clubName">Club name</label>
            <input
                id="clubName"
                type="text"
                bind:value={clubName}
                placeholder="New Club Name"
                required
            />

            <label for="presidentEmail">Club president email</label>
            <input
                id="presidentEmail"
                type="email"
                bind:value={PresidentEmail}
                placeholder="Club president email"
                required
            />

            <button type="submit">Create Club</button>
        </form>

        {#if statusMessage}
            <p class="status-pill" role="status">{statusMessage}</p>
        {/if}
    </section>

    <!-- Admin role preview / switch controls -->
    <aside class="admin-switch-panel">
        <AdminSwitch />
    </aside>
</main>

<!-- Global footer -->
<Footer />

<style>
    /* Page theme + component layout styles */
    :global(body) {
        margin: 0;
        background:
            radial-gradient(circle at 10% 15%, rgba(255, 193, 119, 0.22), transparent 38%),
            radial-gradient(circle at 88% 78%, rgba(64, 145, 108, 0.18), transparent 42%),
            linear-gradient(160deg, #f9f6ef 0%, #eef4f7 100%);
        color: #1b2c34;
        font-family: 'Avenir Next', 'Segoe UI', sans-serif;
    }

    .create-club-page {
        position: relative;
        display: grid;
        grid-template-columns: minmax(0, 2fr) minmax(240px, 1fr);
        gap: 1.2rem;
        width: min(1080px, calc(100% - 2rem));
        margin: 1.4rem auto 2rem;
        padding: 0.35rem;
        isolation: isolate;
    }

    .backdrop-shape {
        position: absolute;
        border-radius: 999px;
        z-index: -1;
        filter: blur(2px);
    }

    .shape-left {
        width: 220px;
        height: 220px;
        left: -60px;
        top: -20px;
        background: rgba(217, 119, 6, 0.16);
    }

    .shape-right {
        width: 190px;
        height: 190px;
        right: -45px;
        bottom: -40px;
        background: rgba(14, 116, 144, 0.14);
    }

    .create-club-card {
        border-radius: 1rem;
        border: 1px solid #d6dee2;
        background: rgba(255, 255, 255, 0.9);
        box-shadow: 0 20px 40px rgba(10, 36, 50, 0.12);
        padding: 1.15rem;
        animation: card-enter 360ms ease-out;
    }

    .card-head {
        margin-bottom: 0.95rem;
    }

    .eyebrow {
        margin: 0;
        text-transform: uppercase;
        letter-spacing: 0.09em;
        font-size: 0.75rem;
        color: #9a3412;
        font-weight: 700;
    }

    h1 {
        margin: 0.25rem 0 0.4rem;
        font-size: clamp(1.4rem, 1.35vw + 1rem, 2rem);
        line-height: 1.1;
        color: #102a3a;
    }

    .subtitle {
        margin: 0;
        color: #3e5b68;
        font-size: 0.98rem;
    }

    .create-form {
        display: grid;
        gap: 0.65rem;
    }

    label {
        font-weight: 700;
        font-size: 0.92rem;
        color: #1f3a47;
    }

    input {
        border: 1px solid #bdd0d9;
        border-radius: 0.72rem;
        padding: 0.72rem 0.8rem;
        background: #fbfdfd;
        transition: border-color 140ms ease, box-shadow 140ms ease, transform 140ms ease;
    }

    input:focus {
        outline: none;
        border-color: #0f766e;
        box-shadow: 0 0 0 3px rgba(15, 118, 110, 0.16);
        transform: translateY(-1px);
    }

    button {
        margin-top: 0.25rem;
        border: none;
        border-radius: 0.74rem;
        padding: 0.8rem 1rem;
        font-weight: 700;
        cursor: pointer;
        color: #ffffff;
        background: linear-gradient(120deg, #0f766e 0%, #0a5e8a 100%);
        transition: transform 140ms ease, box-shadow 140ms ease, filter 140ms ease;
    }

    button:hover {
        transform: translateY(-1px);
        box-shadow: 0 10px 24px rgba(15, 76, 97, 0.28);
        filter: saturate(1.03);
    }

    button:active {
        transform: translateY(0);
    }

    .status-pill {
        margin: 0.85rem 0 0;
        border-radius: 999px;
        padding: 0.48rem 0.82rem;
        width: fit-content;
        background: #e4f7ef;
        border: 1px solid #8fcea8;
        color: #14532d;
        font-size: 0.89rem;
        font-weight: 600;
    }

    .admin-switch-panel {
        border-radius: 1rem;
        border: 1px solid #dce5ea;
        background: rgba(255, 255, 255, 0.82);
        padding: 0.9rem;
        box-shadow: 0 16px 30px rgba(18, 46, 58, 0.1);
        backdrop-filter: blur(2px);
        animation: card-enter 500ms ease-out;
    }

    @keyframes card-enter {
        from {
            opacity: 0;
            transform: translateY(12px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @media (max-width: 900px) {
        .create-club-page {
            grid-template-columns: 1fr;
        }

        .admin-switch-panel {
            order: -1;
        }
    }
</style>

