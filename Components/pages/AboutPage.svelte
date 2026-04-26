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
</script>

<Header userType={userType} />
<LoginPopup />
<div class="page-shell">
    <main class="about-page">
        <section class="hero-section">
            <p class="eyebrow">About</p>
            <h1>About This Project</h1>
            <p class="intro">
                This project was created as part of a computer science capstone project for Southwestern University.
                Our goal is to make it easier for students to find organizations where they feel community, while also
                helping campus organizations connect with new members.
            </p>
        </section>

        <section class="content-grid" aria-label="Project overview and goals">
            <article class="card">
                <h2>Project Purpose</h2>
                <p>
                    SU Organization Matching helps students discover organizations that align with their interests,
                    identity, and campus involvement preferences.
                </p>
                <p>
                    The platform also gives officers and admins the tools they need to keep club information accurate and up to date.
                </p>
            </article>

            <article class="card">
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
        background: #0b1220;
    }

    .page-shell {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        background: #0b1220;
    }

    .about-page {
        --page-bg: linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
        --text-main: #10243a;
        --text-subtle: #31506e;
        --card-border: #dbe7f3;

        flex: 1;
        min-height: 0;
        background: var(--page-bg);
        color: var(--text-main);
        padding: 2rem 1rem 3rem;
        font-size: 16px;
    }

    .hero-section {
        width: min(100%, 980px);
        margin: 0 auto 1.4rem;
        padding: 0 1rem;
        text-align: center;
    }

    .eyebrow {
        margin: 0 0 0.5rem;
        color: #0f6d8c;
        font-size: 0.9rem;
        text-transform: uppercase;
        letter-spacing: 0.12em;
        font-weight: 700;
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
        width: min(100%, 980px);
        margin: 0 auto;
        padding: 0 1rem;
        display: grid;
        gap: 1rem;
        grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .card {
        background: #ffffff;
        border: 1px solid var(--card-border);
        border-radius: 0.95rem;
        padding: 1.05rem;
        box-shadow: 0 10px 24px rgba(13, 37, 62, 0.08);
    }

    .card h2 {
        margin: 0 0 0.55rem;
        font-size: 1.22rem;
        color: var(--text-main);
    }

    .card p {
        margin: 0.5rem 0 0;
        line-height: 1.65;
        color: var(--text-subtle);
    }

    .card ul {
        margin: 0;
        padding-left: 1.2rem;
    }

    .card li {
        margin: 0.5rem 0;
        color: var(--text-subtle);
        line-height: 1.55;
    }

    .card li::marker {
        color: #0f6d8c;
    }

    @media (max-width: 840px) {
        .about-page {
            padding: 1rem 0.85rem 2rem;
        }

        .hero-section,
        .content-grid {
            padding: 0 0.3rem;
        }

        .content-grid {
            grid-template-columns: 1fr;
        }

        h1 {
            font-size: clamp(1.7rem, 7vw, 2.3rem);
        }
    }
</style>

