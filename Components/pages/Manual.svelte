<!-- @component Comprehensive documentation page explaining how to use the organization matching tool for students, officers, and admins. -->
<script>
/**
 * @type {state} userType - Current user's role loaded from authentication token
 * @type {state} isOfficer - Boolean indicating if current user is an officer
 * @type {state} isAdmin - Boolean indicating if current user is an admin
 * @function loadUserType - Fetches and sets user's role from backend, updates officer and admin flags
 * @lifecycle onMount - Loads user type on component initialization
 */
import { onMount } from 'svelte';
import Header from '../header.svelte';
import Footer from '../footer.svelte';
import LoginPopup from '../login_popup.svelte';

let userType = 'user';
let isOfficer = false;
let isAdmin = false;

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
	isOfficer = userType === 'officer';
	isAdmin = userType === 'admin';
}

onMount(() => {
    loadUserType();
});
</script>

<Header userType={userType} />
<LoginPopup />

<div class="page-shell">
<main class="manual-page">
	<section class="hero-section">
	
		<h1>How to use the SU Organization Matching Tool</h1>
		<p class="intro">
			Disclaimer: This tool was made to give suggestions based on the information provided. Please do not take it as the
			end-all be-all, and do not be afraid to explore clubs that it does not provide.
		</p>
	</section>

	<section class="content-grid">
		<article class="content-section">
			<h2>Users</h2>
			<h3>Login</h3>
			<p>
				When you first visit our website, you will be prompted to log in with a Google account. This is used to save your results
				so you can return to view them later.
			</p>

			<h3>Taking the quiz</h3>
			<p>
				After you log in, click on the Take the Quiz button, which
				will take you to the first section of the quiz.
			</p>
			<p>
				This part is a standard demographics form used to gauge your academic interests, as well as your gender,
				religion, ethnicity, and other details. We collect this data because it is important to certain organizations
				on campus.
			</p>
			<p>
				It is important to note that answering “Prefer not to say” for any of the questions is functionally the same for sorting as answering with nothing. 
				This is particularly important for organizations with strict gender matching requirements, as you will not see any of them with this option selected.
			</p>
			<p>
				The next part is a personality quiz that asks about the types of activities you enjoy and what
				personality traits you align with.
			</p>
			<p>
				This section is done in a dating-app style, where you can swipe on your phone or, on a desktop, use the
				mouse or arrow keys. After you are done with this section, you will be automatically redirected to the
				results page to view the results.
			</p>

			<h3>Viewing your results</h3>
			<p>
				You can access your results by clicking the My Results tab on the home page.
			</p>
			<p>
				If you have not taken the quiz at least once, no results will appear. If you have 
				taken the quiz before and want to view the results from your last quiz, do 
				so by clicking the My Results tab. The quiz is not cumulative, so retaking the 
				quiz will cause the website to recalculate your results, not add to your previous scores.

			</p>
			<p>
				On the results page, you will find a list of student organizations you matched with, 
				each with a short description, a list of the organization's activities, any links to 
				social media, and possibly a contact within the club. 

			</p>
		</article>
		{#if isOfficer || isAdmin}
		<article class="card">
			<h2>Officers</h2>
			<h3>Editing a club's information</h3>
			<p>
				In addition to the regular quiz and matching capabilities of a regular user, a student who is an officer
				of one or more clubs will be able to edit information associated with the club(s).
			</p>
			<p>
				You can do this by clicking on the Manage Club tab on the header. If you are an officer of multiple clubs,
				then you can hover over the button to view all of them, then click on the one you want to manage.
			</p>

			<h3>Editing options</h3>
			<p>Below are all the settings our website allows officers to change, along with descriptions of what they are and how they affect the sorting algorithm.</p>
			<ul>
				<li><strong>Description</strong> - Information about the club that is displayed on the results page.</li>
				<li><strong>Activities Description </strong> -A discription  of activities displayed on the results page.</li>
				<li><strong>Meeting Information</strong> - The organization's meeting time, day, and location to be displayed on the results page.</li>
				<li><strong>Social Media/Website</strong> - A link to your organization's social media or website to be displayed on the results page.</li>
				<li><strong>Image</strong> - You can upload an image file to be displayed on the results page to represent your club.</li>
				<li><strong>Include Officer Emails</strong> - A contact email address to be displayed on the results page for students to reach out to with questions about the club on the Results page if checked</li>
				<li><strong>Personality Traits & Activities</strong> - A list of personality traits and activities that you can select to change how the system sorts users into your club.</li>
				<li><strong>Genders</strong> -  Genders that are linked to the purpose of the organization (slightly affects sorting). Please only include this field if it is important to the purpose of your organization.</li>
				<li><strong>Ethnicities</strong> - Races/Ethicities that are linked to the purpose of the organization (affects sorting). Please only include this field if it is important to the purpose of your organization.</li>
				<li><strong>Religions</strong> - Religions that are linked to the purpose of the organization (affects sorting). Please only include this field if it is important to the purpose of your organization.</li>
				<li><strong>Dedicated Majors</strong> - Majors that are linked to the purpose of your organization (affects sorting).</li>
				<li><strong>Associated Majors</strong> - Trends in majors that officers see in their club. (slightly affects sorting).</li>
				<li><strong>Other</strong> - Miscellaneous trends that officers see in their club that affect sorting. Please only include items in this field if they are important to the purpose of your organization.</li>
				<li><strong>Strict Gender Matching</strong> - If checked, makes the Genders selection absolute. Please only include this field if it is important to the purpose of your organization.</li>
				<li><strong>Current Officers</strong> - A list of current officers in the club. You can remove officers by selecting the check box next to their email address, then clicking the Delete Selected Officers button.</li>
				<li><strong>Add New Officer</strong> - Type in an officer's email address and click the Add Officer button to add them as a new officer.</li>
			</ul>
		</article>
		{/if}
		{#if isAdmin}
		<article class="card">
			<h2>Admins</h2>
			<h3>Editing clubs</h3>
			<p>
				As admins, you will be able to edit all clubs on the website as well as create and delete clubs. 
			</p>
			<p>
				You do this by clicking the edit or delete button associated with a club on the home page. 
				The page you see when you click on the edit club button is the same as the officer view for the club, 
				so refer to the Editing Options section above for information on what each part does.
			</p>

			<h3>Creating a New Club</h3>
			<p>
				Creating a new club is really simple: click the Create New Club button in the header, 
				then enter the club name and the president's email address. The system does not send a notification to the new officer, so inform them separately.
			</p>
			<p>
				The president of the new club will then have access to edit it on the site. 
				They just need to fill out the settings page so the website knows more about the club.

			</p>

			<h3>Switching Account Types</h3>
			<p>
				As admins, you may want to see what different views of the website look like. 
				To do this, look at the top of the website, where you will see a banner labeled Change User Type, 
				with the options Regular User, Officer, and Admin to represent the three account types we have.
			</p>
			<p>
				By default, it should have Admin selected, but to change the account type, 
				select the Regular User or Officer button. To go back to an admin account, just reselect the Admin button. 
			</p>
		</article>
		{/if}
	</section>

</main>
</div>
<Footer />

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

	.manual-page {
		--page-bg: linear-gradient(180deg, #f7fbff 0%, #eef6ff 100%);
		--text-main: #10243a;
		--text-subtle: #31506e;

		flex: 1;
		min-height: 0;
		background: var(--page-bg);
		color: var(--text-main);
		padding: 2rem 1rem 3rem;
		font-size: 16px;
	}

	.hero-section {
		width: min(100%, 980px);
		margin: 0 auto 1.5rem;
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
		width: min(100%, 980px);
		margin: 0 auto;
		padding: 0 1rem;
		display: grid;
		grid-template-columns: 1fr;
		gap: 1.5rem;
	}

	.content-section {
		padding: 1.2rem 0;
	}

	.content-section h2 {
		margin: 0 0 0.8rem;
		font-size: 1.35rem;
		color: var(--text-main);
	}

	.content-section h3 {
		margin: 1.2rem 0 0.4rem;
		font-size: 1.05rem;
		color: var(--text-main);
	}

	.content-section p {
		margin: 0.4rem 0 0;
		line-height: 1.7;
		color: var(--text-subtle);
	}

	ul {
		margin: 0;
		padding-left: 1.2rem;
		line-height: 1.7;
		color: var(--text-subtle);
	}

	li + li {
		margin-top: 0.4rem;
	}

	strong {
		color: var(--text-main);
	}

	@media (prefers-color-scheme: dark) {
		.manual-page {
			--page-bg: linear-gradient(180deg, #0b1220 0%, #111827 100%);
			--text-main: #e5edf8;
			--text-subtle: #b6c7df;
		}
	}
</style>