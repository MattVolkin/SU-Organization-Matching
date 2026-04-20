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
<main class="manual-page">
	<section class="hero-content">
		<h1>How to use the SU Organization Matching tool</h1>
		<p class="intro">
			Disclaimer: This tool was made to give suggestions based on the information provided. Please do not take it as the
			end-all be-all, and do not be afraid to explore clubs that it does not provide.
		</p>
	</section>

	<section class="grid">
		<article>
			<h2>Users</h2>
			<h3>Login</h3>
			<p>
				When you first visit our website, you will be prompted to log in with a Google account.
			</p>

			<h3>Taking the quiz</h3>
			<p>
				After you log in, if you are going to take the matching quiz, click on the get started button, which
				will take you to the first part of the quiz.
			</p>
			<p>
				This part is a standard demographics form used to gauge academic interest in the departmental clubs on
				campus, as well as your gender, religion, and other details. This is done because some clubs are gender
				specific and have other restrictions.
			</p>
			<p>
				The next part is a personality quiz that asks you about the types of activities you like and certain
				personality traits that club leaders look for in new members.
			</p>
			<p>
				This section is done in a dating-app style, where you can swipe on your phone or, on a desktop, use the
				mouse or arrow keys. After you are done with this section, you will be automatically redirected to the
				results page to view the results.
			</p>

			<h3>Viewing your results</h3>
			<p>
				You can access your results by clicking the Results tab on the home page.
			</p>
			<p>
				If you have not taken the quiz at least once, no results will appear. If you have taken the quiz before
				and just want to view the results from your last quiz, the quiz is not cumulative, so retaking the quiz
				will cause the website to recalculate your results.
			</p>
			<p>
				On the page itself for each club in your results, you will find a small description about the club, a
				list of activities the club does, any links to social media, and possibly a contact within the club.
			</p>
		</article>

		<article>
			<h2>Officers</h2>
			<h3>Editing a club's information</h3>
			<p>
				In addition to the regular quiz and matching capabilities of a regular user, a student who is an officer
				of one or more clubs will be able to edit information associated with the club.
			</p>
			<p>
				They do this by clicking on the manage club tab on the header. If they are an officer of multiple clubs,
				then they will be able to hover and select the club they want to edit.
			</p>

			<h3>Editing options</h3>
			<p>Here are all the settings our website allows officers to change, along with descriptions of what they are and whether they change the sorting algorithm at all.</p>
			<ul>
				<li><strong>Description</strong> - what the club is about for the result page</li>
				<li><strong>Activities (comma-separated)</strong> - a list of activities for the results page</li>
				<li><strong>Meeting time</strong> - for the results page</li>
				<li><strong>Social media/website</strong> - for the results page</li>
				<li><strong>Image</strong> - you can upload a file for the result page</li>
				<li><strong>Personality + activities trait select</strong> - this is how you update the activities and personality traits that the system sorts users into your club</li>
				<li><strong>Genders (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Ethnicities (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Religions (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Dedicated majors (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Other (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Strict gender matching</strong> - if checked, makes the Genders section absolute</li>
				<li><strong>Adding and removing an officer</strong> - how officers edit who can edit their clubs</li>
			</ul>
		</article>

		<article>
			<h2>Admins</h2>
			<h3>Editing clubs</h3>
			<p>
				As admins, you will be able to edit all clubs on the website and delete a club.
			</p>
			<p>
				You do this by clicking the edit or delete button associated with a club on the home page. The page you
				will see when you click on the edit club button is the same as the officers for the club, so refer to the
				Editing Options section above for what each section on the page does.
			</p>

			<h3>Creating a new club</h3>
			<p>
				Creating a new club is really simple: click the create new club button in the header, then enter the club
				name and the president's email address.
			</p>
			<p>
				Then just send an email to the president. The president of the new club will then have access to edit it
				on the site. They just need to fill out the settings page so the website knows more about the club.
			</p>

			<h3>Switching account types</h3>
			<p>
				As admins, you may want to see what different views of the website look like. To do this for our website,
				look at the top of the website, where you will see a banner called Switch account types, with the options
				User, Officer, and Admin to represent the three account types we have.
			</p>
			<p>
				By default, it should have admin selected, but to change the account type, select the user or officer
				button, and to go back to an admin account, just reselect the admin button.
			</p>
		</article>
	</section>
</main>
<Footer/>
</div>

<style>
	:global(html),
	:global(body) {
		margin: 0;
		padding: 0;
		min-height: 100%;
		background: #0b1220;
	}

	:global(body) {
		overflow-x: hidden;
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
		--action: #1f6f8b;
		--action-hover: #195d76;

		flex: 1;
		min-height: 0;
		background: var(--page-bg);
		color: var(--text-main);
		padding: 2rem 1rem 2.25rem;
		font-size: 16px;
	}

	.hero-content {
		width: min(100%, 760px);
		margin: 0 auto;
		padding: 1.25rem 1rem 0.65rem;
		text-align: left;
	}

	h1 {
		margin: 0.35rem 0 0.8rem;
		font-size: clamp(1.9rem, 2.7vw + 1rem, 3rem);
		line-height: 1.12;
		max-width: 22ch;
	}

	.intro {
		max-width: 68ch;
		margin: 0;
		font-size: 1.08rem;
		line-height: 1.65;
		color: var(--text-subtle);
	}

	.grid {
		width: min(100%, 760px);
		margin: 0.45rem auto 0;
		display: block;
		padding: 0 1rem;
	}

	.grid article {
		padding: 0.55rem 0;
	}

	.grid article + article {
		border-top: 1px solid #dbe7f3;
		margin-top: 0.3rem;
		padding-top: 1.2rem;
	}

	.grid h2 {
		margin: 0 0 0.75rem;
		font-size: 1.25rem;
		color: var(--text-main);
	}

	.grid h3 {
		margin: 1.05rem 0 0.35rem;
		font-size: 1rem;
		color: #1b4f73;
	}

	.grid p {
		margin: 0.3rem 0 0;
		line-height: 1.62;
		color: var(--text-subtle);
	}

	ul {
		margin: 0;
		padding-left: 1.2rem;
		line-height: 1.62;
		color: var(--text-subtle);
	}

	li::marker {
		color: var(--action);
	}

	li + li {
		margin-top: 0.35rem;
	}

	strong {
		color: var(--text-main);
	}

	@media (max-width: 820px) {
		.grid {
			padding: 0 0.6rem;
		}

		.grid article + article {
			padding-top: 1rem;
		}

		.manual-page {
			padding: 1rem 0.85rem 1.7rem;
		}

		.hero-content {
			padding: 1rem 0.6rem 0.6rem;
		}

		h1 {
			font-size: clamp(1.65rem, 7vw, 2.2rem);
		}

		.intro {
			font-size: 1rem;
		}
	}

	@media (prefers-color-scheme: dark) {
		.manual-page {
			--page-bg: linear-gradient(180deg, #0b1220 0%, #111827 100%);
			--text-main: #e5edf8;
			--text-subtle: #b6c7df;
			--action: #2b8fb5;
			--action-hover: #3aa3cb;
		}

		.grid article + article {
			border-top-color: #273449;
		}

		.grid h3 {
			color: #7dc4e0;
		}
	}
</style>