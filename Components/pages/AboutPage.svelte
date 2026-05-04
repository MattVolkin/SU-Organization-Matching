<!-- @component About page displaying project information, goals, and team credits for the organization matching system. -->
<script>
/**
 * @type {state} userType - Current user's role loaded from authentication token
 * @function loadUserType - Fetches and sets user's role from backend
 * @lifecycle onMount - Loads user type on component initialization
 */
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
</script>

<Header userType={userType} />
<LoginPopup />
<div class="page-shell">
    <main class="about-page">
        <section class="hero-section">
            <h1>About This Project</h1>
            <p class="intro">
                This project was created as part of a computer science capstone project for Southwestern University.
                Our goal is to make it easier for students to find organizations where they feel community, while also
                helping campus organizations connect with new members.
            </p>
        </section>

        <section class="content-grid" aria-label="Project overview and goals">
            <article class="content-section">
                <h2>Goals</h2>
                <ul>
                    <li>Create a short quiz that narrows down choices so students do not need to browse all 75+ organizations manually.</li>
                    <li>Keep the experience easy to use and accessible for students.</li>
                    <li>Help students find organizations where they can build a strong sense of community.</li>
                    <li>Provide officers with tools to update organization information so results stay current.</li>
                </ul>
            </article>
        </section>
    </main>

    <Footer />
</div>

<style>
    :global(html),
    :global(body) {
        margin: 0;
        padding: 0;
        min-height: 100%;
        background: #ffffff;
    }

    .page-shell {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        background: #ffffff;
    }

    .about-page {
        --page-bg: linear-gradient(180deg, #ffffff 0%, #f3f3f3 100%);
        --text-main: #000000;
        --text-subtle: #4a4a4a;
        --accent: #ffcd00;
        --divider: #828282;

        flex: 1;
        min-height: 0;
        background: var(--page-bg);
        color: var(--text-main);
        padding: 2rem 1rem 3rem;
        font-size: 16px;
        font-family: system-ui, -apple-system, Segoe UI, sans-serif;
    }

    .hero-section {
        width: min(100%, 980px);
        margin: 0 auto 1.2rem;
        padding: 0 1rem;
        text-align: center;
    }


    h1 {
        margin: 0.55rem 0 0.9rem;
        font-size: clamp(2rem, 3vw + 1rem, 3.1rem);
        line-height: 1.12;
    }

    .intro {
        max-width: 70ch;
        margin: 0 auto;
        color: var(--text-subtle);
        line-height: 1.65;
        font-size: 1.02rem;
    }

    .content-grid {
        width: min(100%, 760px);
        margin: 0 auto;
        padding: 0 1rem;
    }

    .content-section {
        border-top: 1px solid var(--divider);
        padding-top: 1rem;
    }

    .content-section h2 {
        margin: 0 0 0.55rem;
        font-size: 1.22rem;
        color: var(--text-main);
    }

    .content-section ul {
        margin: 0;
        padding-left: 1.15rem;
    }

    .content-section li {
        margin: 0.55rem 0;
        color: var(--text-subtle);
        line-height: 1.6;
    }

    .content-section li::marker {
        color: #000000;
    }

    @media (max-width: 840px) {
        .about-page {
            padding: 1rem 0.85rem 2rem;
        }

        .hero-section,
        .content-grid {
            padding: 0 0.3rem;
        }

        h1 {
            font-size: clamp(1.7rem, 7vw, 2.3rem);
        }
    }

    @media (prefers-color-scheme: dark) {
        :global(html),
        :global(body),
        .page-shell {
            background: #000000;
        }

        .about-page {
            --page-bg: linear-gradient(180deg, #000000 0%, #1e1e1e 100%);
            --text-main: #ffffff;
            --text-subtle: #d5d5d5;
            --accent: #ffcd00;
            --divider: #828282;
        }
    }
</style>

