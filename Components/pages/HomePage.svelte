<script>
import { onMount } from 'svelte';
import Header from '../header.svelte';
import Footer from '../footer.svelte';
import LoginPopup from '../login_popup.svelte';

export let previewAs = '';
export let showChrome = true;

let userType = 'user';

function getActiveUserType() {
    return previewAs || userType;
}

async function loadUserType() {
    if (previewAs) {
        userType = previewAs;
        return;
    }

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

    if (userType === 'admin') {
        window.location.replace('/admin-home.html');
    }
}

onMount(() => {
    loadUserType();
});

function goToQuiz() {
    window.location.href = '/demographic-quiz.html';
}

function goToDeleteAccount() {
    window.location.href = '/delete-account.html';
}

</script>

{#if showChrome}
    <Header userType={getActiveUserType()} previewAs={previewAs} />
    <LoginPopup />
{/if}

<div class="page-shell">
    <main class="home-page">
        <!-- Hero section -->
        <section class="hero-content">
            <h1>Find The Organization That Fits You</h1>
            <p class="lead">
                Explore clubs and student organizations based on your interests, activities,
                and preferences, then get matched with communities where you can thrive.
            </p>
        </section>

        <!-- How-it-works section -->
        <section class="info-list" aria-label="How it works">
            <h2>How It Works</h2>
            <ul>
                <li><strong>Take The Quiz:</strong> Answer a short set of questions so we can understand your interests and preferences.</li>
                <li><strong>Get Matches:</strong> See organizations ranked by compatibility and discover communities that align with you.</li>
                <li><strong>Connect Faster:</strong> Use your matches to quickly find meeting times, activities, and officers to contact.</li>
            </ul>
        </section>

        <!-- Primary call-to-action section -->
        <div class="cta-row">
            <button type="button" class="cta-primary" onclick={goToQuiz}>Get Started</button>
            <a class="cta-secondary" href="/results.html">View Results</a>
        </div>
    </main>

    {#if showChrome}
        <!-- Account management + global footer -->
        <button type="button" onclick={goToDeleteAccount}>Delete Account Information</button>

        <Footer />
    {/if}
</div>

<style>

    :global(html),
    :global(body) {
        margin: 0;
        padding: 0;
        min-height: 100%;
        background: #ffffff;
    }

    :global(body) {
        overflow-x: hidden;
    }

    .page-shell {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        background: #ffffff;
    }

    .home-page {
        --page-bg: linear-gradient(180deg, #ffffff 0%, #f3f3f3 100%);
        --text-main: #000000;
        --text-subtle: #454545;
        --action: #000000;
        --action-hover: #828282;
        --cta-secondary-bg: #ffcd00;
        --cta-secondary-hover: #e5b800;
        --cta-secondary-text: #000000;
        --cta-secondary-border: #000000;
        --focus-ring: #828282;

        flex: 1;
        min-height: 0;
        background: var(--page-bg);
        color: var(--text-main);
        padding: 2rem 1rem;
        font-size: 16px;
    }

    .hero-content {
        width: min(100%, 980px);
        margin: 0 auto;
        padding: 1.4rem 1rem 1rem;
        text-align: center;
    }


    h1 {
        margin: 0.55rem 0 0.9rem 0;
        font-size: clamp(2rem, 3vw + 1rem, 3.1rem);
        line-height: 1.12;
    }

    .lead {
        margin: 0 auto;
        max-width: 54rem;
        color: var(--text-subtle);
        line-height: 1.65;
        font-size: 1.12rem;
    }

    .cta-row {
        width: min(100%, 900px);
        margin: 0 auto 1.35rem auto;
        padding: 0 1rem;
        display: flex;
        gap: 0.9rem;
        flex-wrap: wrap;
        align-items: center;
        justify-content: center;
    }

    button {
        border: none;
        border-radius: 0.65rem;
        padding: 0.72rem 1.2rem;
        font-size: 1rem;
        font-weight: 700;
        color: #ffffff;
        background: var(--action);
        cursor: pointer;
        transition: background-color 0.2s ease, transform 0.2s ease;
    }

    button:hover {
        background: var(--action-hover);
        transform: translateY(-1px);
    }

    button:focus-visible,
    a:focus-visible {
        outline: 2px solid var(--focus-ring);
        outline-offset: 2px;
    }

    .cta-row .cta-primary,
    .cta-row .cta-secondary {
        display: inline-block;
        border-radius: 0.65rem;
        padding: 0.72rem 1.2rem;
        font-weight: 700;
        font-size: 1rem;
        line-height: 1.2;
        text-decoration: none;
        transition: background-color 0.2s ease, color 0.2s ease, transform 0.2s ease;
    }

    .cta-row .cta-primary {
        color: #ffffff;
        background: var(--action);
        border: 1px solid #000000;
    }

    .cta-row .cta-primary:hover {
        background: var(--action-hover);
        transform: translateY(-1px);
    }

    .cta-row .cta-secondary {
        color: var(--cta-secondary-text);
        background: var(--cta-secondary-bg);
        border: 1px solid var(--cta-secondary-border);
    }

    .cta-row .cta-secondary:hover {
        color: var(--cta-secondary-text);
        background: var(--cta-secondary-hover);
        transform: translateY(-1px);
    }

    .info-list {
        width: min(100%, 900px);
        margin: 0.6rem auto 0 auto;
        padding: 1.1rem 1rem 1.4rem;
        text-align: center;
        border-top: 1px solid #828282;
    }

    .info-list h2 {
        margin: 0 0 0.65rem 0;
        font-size: clamp(1.15rem, 1.1vw + 0.9rem, 1.5rem);
    }

    .info-list ul {
        margin: 0 auto;
        padding-left: 1.25rem;
        width: min(100%, 50rem);
        text-align: left;
    }

    .info-list li {
        margin: 0.5rem 0;
        color: var(--text-subtle);
        line-height: 1.6;
        font-size: 1.02rem;
    }

    .info-list li::marker {
        color: var(--action);
    }

    @media (max-width: 640px) {
        .home-page {
            padding: 1rem 0.85rem;
        }

        .hero-content,
        .info-list {
            padding: 1rem 0.9rem;
        }

        h1 {
            font-size: clamp(1.65rem, 7vw, 2.25rem);
        }

        .lead {
            font-size: 1rem;
        }

        .cta-row {
            flex-direction: column;
            align-items: stretch;
        }

        .cta-row .cta-primary,
        .cta-row .cta-secondary {
            width: 100%;
            text-align: center;
        }
    }

    @media (prefers-color-scheme: dark) {
        :global(html),
        :global(body),
        .page-shell {
            background: #000000;
        }

        .home-page {
            --page-bg: linear-gradient(180deg, #000000 0%, #1e1e1e 100%);
            --text-main: #ffffff;
            --text-subtle: #d5d5d5;
            --action: #ffcd00;
            --action-hover: #e5b800;
            --cta-secondary-bg: #ffffff;
            --cta-secondary-hover: #e6e6e6;
            --cta-secondary-text: #000000;
            --cta-secondary-border: #828282;
            --focus-ring: #ffcd00;
        }

        .info-list {
            border-top-color: #828282;
        }
    }
</style>
