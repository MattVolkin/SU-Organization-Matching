<script>
/**
 * @type {state} clubName - Name entered for the new club
 * @type {state} PresidentEmail - Email entered for the first club officer
 * @type {state} statusMessage - Confirmation message shown after successful submission
 * @function createClub - Submits the new club to the admin API, then clears the form and shows a status message
 */
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

<div class="create-club-shell">
    <!-- Top navigation -->
    <Header userType="admin" />

    <!-- Main admin layout -->
    <main class="create-club-page">
        <!-- Club creation form card -->
        <section class="create-club-card" aria-labelledby="create-club-title">
            <div class="card-head">
                <h1 id="create-club-title">Create a New Club</h1>
                <p class="subtitle">To add a new club, fill out the form below.</p>
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
    </main>

    <!-- Global footer -->
    <Footer />
</div>

<style>
    :global(html),
    :global(body) {
        min-height: 100%;
        margin: 0;
        background: #ffffff;
        color: #000000;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }

    .create-club-shell {
        --page-bg: linear-gradient(180deg, #ffffff 0%, #f3f3f3 100%);
        --text-main: #000000;
        --text-subtle: #4a4a4a;
        --card-bg: #ffffff;
        --card-border: #828282;
        --input-bg: #ffffff;
        --input-border: #828282;
        --action: #000000;
        --action-hover: #828282;
        --status-bg: #f3f3f3;
        --status-border: #828282;
        --status-text: #000000;

        min-height: 100vh;
        display: flex;
        flex-direction: column;
        background: var(--page-bg);
        color: var(--text-main);
    }

    .create-club-page {
        flex: 1;
        width: min(100%, 1040px);
        margin: 1rem auto 1.25rem auto;
        padding: 0 1rem;
        box-sizing: border-box;
    }

    .create-club-card {
        border-radius: 1rem;
        background: var(--card-bg);
        border: 1px solid var(--card-border);
        box-shadow: 0 10px 24px rgba(0, 0, 0, 0.1);
        padding: 1rem;
    }

    .card-head {
        margin-bottom: 0.85rem;
    }


    h1 {
        margin: 0.3rem 0 0.35rem;
        font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
        line-height: 1.2;
        color: var(--text-main);
    }

    .subtitle {
        margin: 0;
        color: var(--text-subtle);
        font-size: 0.95rem;
    }

    .create-form {
        display: grid;
        gap: 0.55rem;
    }

    label {
        font-weight: 700;
        font-size: 0.9rem;
        color: var(--text-main);
    }

    input {
        width: min(100%, 36rem);
        box-sizing: border-box;
        border: 1px solid var(--input-border);
        border-radius: 0.55rem;
        padding: 0.56rem 0.68rem;
        font-size: 0.94rem;
        background: var(--input-bg);
        color: var(--text-main);
        transition: border-color 140ms ease, box-shadow 140ms ease;
    }

    input:focus {
        outline: none;
        border-color: var(--action);
        box-shadow: 0 0 0 3px rgba(255, 205, 0, 0.24);
    }

    button {
        margin-top: 0.35rem;
        border: 1px solid transparent;
        border-radius: 0.6rem;
        padding: 0.62rem 0.9rem;
        font-weight: 700;
        font-size: 0.94rem;
        cursor: pointer;
        color: #ffffff;
        background: var(--action);
        transition: background-color 140ms ease, transform 140ms ease, box-shadow 140ms ease;
    }

    button:hover {
        background: var(--action-hover);
        transform: translateY(-1px);
        box-shadow: 0 8px 18px rgba(0, 0, 0, 0.18);
    }

    button:active {
        transform: translateY(0);
    }

    .status-pill {
        margin: 0.8rem 0 0;
        border-radius: 0.6rem;
        padding: 0.5rem 0.7rem;
        width: fit-content;
        background: var(--status-bg);
        border: 1px solid var(--status-border);
        color: var(--status-text);
        font-size: 0.89rem;
        font-weight: 600;
    }

    @media (max-width: 900px) {
        .create-club-page {
            margin-top: 0.8rem;
            margin-bottom: 1rem;
            padding: 0 0.8rem;
        }
    }

    @media (prefers-color-scheme: dark) {
        :global(html),
        :global(body) {
            background: #000000;
            color: #ffffff;
        }

        .create-club-shell {
            --page-bg: linear-gradient(180deg, #000000 0%, #1e1e1e 100%);
            --text-main: #ffffff;
            --text-subtle: #d5d5d5;
            --card-bg: #121212;
            --card-border: #828282;
            --input-bg: #1e1e1e;
            --input-border: #828282;
            --action: #ffcd00;
            --action-hover: #e5b800;
            --status-bg: #1e1e1e;
            --status-border: #828282;
            --status-text: #ffffff;
        }

        .create-club-card,
        input,
        .status-pill {
            background: var(--card-bg);
            color: var(--text-main);
        }

        button {
            color: #000000;
        }
    }
</style>

