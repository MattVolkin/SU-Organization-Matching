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
        background: linear-gradient(180deg, #edf4fb 0%, #f6f9fd 100%);
        color: #132c45;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }

    .create-club-shell {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
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
        background: #ffffff;
        border: 1px solid #d4e0ec;
        box-shadow: 0 10px 24px rgba(13, 37, 62, 0.1);
        padding: 1rem;
    }

    .card-head {
        margin-bottom: 0.85rem;
    }


    h1 {
        margin: 0.3rem 0 0.35rem;
        font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
        line-height: 1.2;
        color: #132c45;
    }

    .subtitle {
        margin: 0;
        color: #4f6781;
        font-size: 0.95rem;
    }

    .create-form {
        display: grid;
        gap: 0.55rem;
    }

    label {
        font-weight: 700;
        font-size: 0.9rem;
        color: #132c45;
    }

    input {
        width: min(100%, 36rem);
        box-sizing: border-box;
        border: 1px solid #d4e0ec;
        border-radius: 0.55rem;
        padding: 0.56rem 0.68rem;
        font-size: 0.94rem;
        background: #ffffff;
        transition: border-color 140ms ease, box-shadow 140ms ease;
    }

    input:focus {
        outline: none;
        border-color: #0f6d8c;
        box-shadow: 0 0 0 3px rgba(15, 109, 140, 0.16);
    }

    button {
        margin-top: 0.35rem;
        border: none;
        border-radius: 0.6rem;
        padding: 0.62rem 0.9rem;
        font-weight: 700;
        font-size: 0.94rem;
        cursor: pointer;
        color: #ffffff;
        background: #0f6d8c;
        transition: background-color 140ms ease, transform 140ms ease, box-shadow 140ms ease;
    }

    button:hover {
        background: #0b5972;
        transform: translateY(-1px);
        box-shadow: 0 8px 18px rgba(13, 37, 62, 0.18);
    }

    button:active {
        transform: translateY(0);
    }

    .status-pill {
        margin: 0.8rem 0 0;
        border-radius: 0.6rem;
        padding: 0.5rem 0.7rem;
        width: fit-content;
        background: #eef8ff;
        border: 1px solid #b9dff0;
        color: #0f4660;
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
</style>

