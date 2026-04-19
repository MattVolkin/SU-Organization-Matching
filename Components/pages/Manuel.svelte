<main class="manual-page">
	<section class="hero">
		<p class="eyebrow">User Manual</p>
		<h1>How to use the SU Organization Matching tool</h1>
		<p class="intro">
			This tool was made to give suggestions based on the information provided. Please do not take it as the
			end-all be-all, and do not be afraid to explore clubs that it does not provide.
		</p>
	</section>

	<section class="grid">
		<article class="card">
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

		<article class="card">
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

		<article class="card">
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

<style>
	:global(body) {
		margin: 0;
		background:
			radial-gradient(circle at top left, rgba(20, 184, 166, 0.14), transparent 28%),
			radial-gradient(circle at top right, rgba(244, 177, 66, 0.14), transparent 24%),
			linear-gradient(180deg, #f8fbfd 0%, #eef5f7 100%);
		color: #173042;
		font-family: Georgia, 'Times New Roman', serif;
	}

	.manual-page {
		width: min(100%, 1100px);
		margin: 0 auto;
		padding: 2rem 1rem 3rem;
	}

	.hero {
		padding: 1.75rem 1.5rem 1.25rem;
		border-radius: 1.4rem;
		background: rgba(255, 255, 255, 0.8);
		border: 1px solid rgba(23, 48, 66, 0.1);
		box-shadow: 0 18px 40px rgba(21, 44, 58, 0.08);
		backdrop-filter: blur(8px);
		margin-bottom: 1.25rem;
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
		margin: 0;
		font-size: clamp(2rem, 3vw + 1rem, 3.8rem);
		line-height: 1.05;
		max-width: 12ch;
	}

	.intro {
		max-width: 70ch;
		margin: 0.9rem 0 0;
		font-size: 1.06rem;
		line-height: 1.65;
		color: #355063;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 1rem;
	}

	.card {
		background: rgba(255, 255, 255, 0.9);
		border: 1px solid rgba(23, 48, 66, 0.1);
		border-radius: 1.15rem;
		padding: 1.15rem 1.2rem;
		box-shadow: 0 12px 30px rgba(21, 44, 58, 0.06);
	}

	.card h2 {
		margin: 0 0 0.6rem;
		font-size: 1.25rem;
		color: #123145;
	}

	.card h3 {
		margin: 1rem 0 0.35rem;
		font-size: 1rem;
		color: #173042;
	}

	.card p {
		margin: 0.3rem 0 0;
		line-height: 1.7;
		color: #2e4657;
	}

	ul {
		margin: 0;
		padding-left: 1.2rem;
		line-height: 1.7;
		color: #2e4657;
	}

	li + li {
		margin-top: 0.35rem;
	}

	strong {
		color: #123145;
	}

	@media (max-width: 820px) {
		.grid {
			grid-template-columns: 1fr;
		}

		.manual-page {
			padding-top: 1rem;
		}
	}
</style>