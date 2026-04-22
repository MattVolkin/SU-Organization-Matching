<main class="manual-page">
	<section class="hero">
		<p class="eyebrow">User Manual</p>
		<h1>How to use the SU Organization Matching Tool</h1>
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

		<article class="card">
			<h2>Officers</h2>
			<h3>Editing a club's information</h3>
			<p>
				In addition to the regular quiz and matching capabilities of a regular user, a student who is an officer
				of one or more clubs will be able to edit information associated with the club(s).
			</p>
			<p>
				You do this by clicking on the Manage Club tab on the header. If you are an officer of multiple clubs,
				then you can hover over the button to view all of them, then click on the one you want to manage.
			</p>

			<h3>Editing options</h3>
			<p>Below are all the settings our website allows officers to change, along with descriptions of what they are and how they affect the sorting algorithm.</p>
			<ul>
				<li><strong>Description</strong> - information about the club that is displayed on the results page</li>
				<li><strong>Activities Description (comma-separated)</strong> - a list of activities displayed the results page</li>
				<li><strong>Meeting Information</strong> - the organization's meeting time, day, and location displayed on the results page</li>
				<li><strong>Social media/website</strong> - a link to your organization's website or social media displayed on the results page</li>
				<li><strong>Image</strong> - you can upload a file to be displayed on the results page</li>
				<li><strong>Personality Trait & Activity Select</strong> - a list of personality traits and activities that you can select to change how the system sorts users into your club</li>
				<li><strong>Genders (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Ethnicities (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Religions (comma-separated)</strong> - trends that officers see in their club slightly affect sorting</li>
				<li><strong>Dedicated Majors (comma-separated)</strong> - a list of majors that are inherently linked to a departmental organization </li>
				<li><strong>Associated Majors (comma-separated)<strong> - a list of majors that are common in your current members</li>
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