<script>
import { onMount } from 'svelte';
import Header from '../header.svelte';
import Footer from '../footer.svelte';
import LoginPopup from '../login_popup.svelte';

let userType = 'user';

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

onMount(() => {
    loadUserType();
});

function goToQuiz() {
    window.location.href = '/demographicQuiz';
}
</script>

<Header userType={userType} />
<LoginPopup />

<main class="home-page">
    <section class="hero-content">
        <p class="kicker">Southwestern University</p>
        <h1>Find The Organization That Fits You</h1>
        <p class="lead">
            Explore clubs and student organizations based on your interests, activities,
            and preferences, then get matched with communities where you can thrive.
        </p>

        <div class="cta-row">
            <button type="button" onclick={goToQuiz}>Get Started</button>
            <a href="/Results">View Results</a>
        </div>
    </section>

    <section class="info-list" aria-label="How it works">
        <h2>How It Works</h2>
        <ul>
            <li><strong>Take The Quiz:</strong> Answer a short set of questions so we can understand your interests and preferences.</li>
            <li><strong>Get Matches:</strong> See organizations ranked by compatibility and discover communities that align with you.</li>
            <li><strong>Connect Faster:</strong> Use your matches to quickly find meeting times, activities, and officers to contact.</li>
        </ul>
    </section>
</main>

<Footer />

<style>
    .home-page {
        --page-bg: linear-gradient(180deg, #edf4fb 0%, #f6f9fd 100%);
        --text-main: #132c45;
        --text-subtle: #4f6781;
        --action: #0f6d8c;
        --action-hover: #0b5972;
        --focus-ring: #60a5fa;

        min-height: calc(100vh - 220px);
        background: var(--page-bg);
        color: var(--text-main);
        padding: 1rem;
    }

    .hero-content {
        width: min(100%, 1040px);
        margin: 0 auto;
        padding: 1.25rem;
        text-align: center;
    }

    .kicker {
        margin: 0;
        color: #3f5f80;
        font-weight: 700;
        letter-spacing: 0.04em;
        text-transform: uppercase;
        font-size: 0.78rem;
    }

    h1 {
        margin: 0.45rem 0 0.75rem 0;
        font-size: clamp(1.5rem, 2vw + 1rem, 2.4rem);
        line-height: 1.15;
    }

    .lead {
        margin: 0;
        max-width: 70ch;
        color: var(--text-subtle);
        line-height: 1.55;
    }

    .cta-row {
        margin-top: 1rem;
        display: flex;
        gap: 0.7rem;
        flex-wrap: wrap;
        align-items: center;
    }

    button {
        border: none;
        border-radius: 0.55rem;
        padding: 0.55rem 0.95rem;
        font-size: 0.95rem;
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

    .cta-row a {
        display: inline-block;
        color: #1b4f73;
        font-weight: 600;
        text-decoration: none;
    }

    .cta-row a:hover {
        text-decoration: underline;
    }

    .info-list {
        width: min(100%, 1040px);
        margin: 1rem auto 0 auto;
        padding: 0 1.25rem 1.25rem;
        text-align: center;
    }

    .info-list h2 {
        margin: 0 0 0.45rem 0;
        font-size: 1.1rem;
    }

    .info-list ul {
        margin: 0;
        padding-left: 0;
        list-style-position: inside;
    }

    .info-list li {
        margin: 0.35rem 0;
        color: var(--text-subtle);
        line-height: 1.5;
    }

    @media (max-width: 640px) {
        .home-page {
            padding: 0.85rem;
        }

        .hero-content,
        .info-list {
            padding: 0.9rem;
        }

        .cta-row {
            flex-direction: column;
            align-items: stretch;
        }

        button {
            width: 100%;
        }
    }
</style>
