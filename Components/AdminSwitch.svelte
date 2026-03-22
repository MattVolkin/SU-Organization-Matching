<!-- @component Creates a component for toggling between different admin preview modes. that the admin can use to preview the site as different user types. -->
<script>
/**
 * @type {props} enabled - boolean to control whether the admin switch is shown or not, defaults to false
 * @type {props} value - the current preview mode, can be 'admin', 'officer' or 'user', defaults to 'admin'
 * @type {props} onChange - callback function that is called when the preview mode is changed, receives the next view as an argument
 * @function setView - helper function to change the preview mode when a button is clicked, calls the onChange callback with the next view if it is different from the current value
 */

	let {
		enabled = false,
		value = 'admin',
		onChange = () => {},
	} = $props();

	function setView(nextView) {
		if (nextView === value) {
			return;
		}
		onChange(nextView);
	}
</script>
<!-- only if enabled  it creates the admin banner that when clicked calls the set view function to view the page as the selected user type -->
{#if enabled}
	<section class="admin-banner" aria-label="Admin view selector"> <!-- aria-label="Admin view selector" is for accessibility -->
		<p class="label">Admin preview mode</p>
		<div class="options" role="tablist" aria-label="Preview as">
			<button
				type="button"
				role="tab"
				aria-selected={value === 'user'}
				class:selected={value === 'user'}
				onclick={() => setView('user')}
			>
				Regular User
			</button>
			<button
				type="button"
				role="tab"
				aria-selected={value === 'officer'}
				class:selected={value === 'officer'}
				onclick={() => setView('officer')}
			>
				Officer
			</button>
			<button
				type="button"
				role="tab"
				aria-selected={value === 'admin'}
				class:selected={value === 'admin'}
				onclick={() => setView('admin')}
			>
				Admin
			</button>
		</div>
	</section>
{/if}

<style>
	.admin-banner {
		position: sticky;
		top: 0;
		z-index: 60;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
		padding: 0.55rem 0.85rem;
		background: linear-gradient(90deg, #f2f7ff 0%, #e5f0ff 100%);
		border-bottom: 1px solid #9fbff0;
		color: #1d334d;
		box-shadow: 0 2px 10px rgba(31, 53, 84, 0.12);
	}

	.label {
		margin: 0;
		font-size: 0.9rem;
		font-weight: 600;
		letter-spacing: 0.01em;
	}

	.options {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		flex-wrap: wrap;
	}

	.options button {
		border: 1px solid #b8d1f5;
		border-radius: 999px;
		background: #ffffff;
		color: #23415f;
		font-size: 0.8rem;
		font-weight: 600;
		padding: 0.35rem 0.65rem;
		cursor: pointer;
	}

	.options button.selected {
		background: #1f6f8b;
		border-color: #1f6f8b;
		color: #ffffff;
	}

	.options button:focus-visible {
		outline: 2px solid #60a5fa;
		outline-offset: 2px;
	}

	@media (max-width: 700px) {
		.admin-banner {
			flex-direction: column;
			align-items: stretch;
			gap: 0.5rem;
		}

		.options {
			width: 100%;
		}

		.options button {
			flex: 1 1 auto;
			text-align: center;
		}
	}

	@media (prefers-color-scheme: dark) {
		.admin-banner {
			background: #0f1d33;
			border-bottom-color: #2a3f62;
			color: #dbe9ff;
		}

		.options button {
			background: #162642;
			border-color: #35507b;
			color: #d6e5ff;
		}

		.options button.selected {
			background: #2b8fb5;
			border-color: #2b8fb5;
			color: #ffffff;
		}
	}
</style>
