<script>

// TODO make it look like its somthing that can be swipable
	// put items on a stack of cards rather than just box on a page
	
  	import { blur, fade, fly, scale, slide } from 'svelte/transition'; //import transition animations
  	import { APICreater } from './APIHandler.svelte'; // import the API handler to make calls to the backend and get the data for the cards
	import { redirect } from '@sveltejs/kit';



	let counter = $state(0); // count what term we are on to stay on the list
	let directionInt = $state('none'); // gets direction information from swiping app
	let items = $state([]); // active card list used by the component at runtime
	let showServerErrorPopup = $state(false);
	let isLoadingQuestions = $state(true);

	function normalizeApiCardItems(payload) {
		if (!Array.isArray(payload)) {
			return [];
		}

		return payload
			.map((item) => {
				if (!item || typeof item !== 'object') {
					return null;
				}

				const questionType = item.question_type === 'personality_traits'
					? 'personality'
					: (item.question_type || 'activities');
				const englishTranslation = item.translations?.en;
				const englishTerm = Array.isArray(englishTranslation)
					? (englishTranslation[0] || '')
					: (typeof item.en === 'string' ? item.en : (item.en?.term || ''));
				const englishDef = Array.isArray(englishTranslation)
					? (englishTranslation[1] || '')
					: (typeof item.en?.def === 'string' ? item.en.def : '');

				if (!englishTerm) {
					return null;
				}

				return {
					id: item.id,
					question_type: questionType,
					en: {
						term: englishTerm,
						def: englishDef,
					},
				};
			})
			.filter(Boolean);
	}



	async function loadAdjectivesAndPersonalityTraits() { // load the adjectives and personality traits from the backend to be used in the cards
		isLoadingQuestions = true;
		showServerErrorPopup = false;

		try {
			const adjectivesList = await APICreater('GET', '/api/adjectives', null);
			const normalizedItems = normalizeApiCardItems(adjectivesList);

			if (normalizedItems.length === 0) {
				throw new Error('No questions were returned by the server');
			}

			items = normalizedItems;
			counter = 0;
		} catch (error) {
			console.error('Unable to load adjectives and personality traits', error);
			items = [];
			showServerErrorPopup = true;
		} finally {
			isLoadingQuestions = false;
		}
	}


	export async function sendSwipeInformation(responses) { // send all swipe responses at once using the array payload expected by /response

		try {
	      	await APICreater('POST', '/response', responses);
	} catch (error) {
      console.error('Unable to send swipe information', error);
    }
	}


	const range = (start, end, targetArray) => { // range function to get a sub-array of the full question array (without modifying the original)
    let arr = [];
    for (let i = start; i < end; i += 1) arr.push(targetArray[i]);
    return arr;
  };

  let lang = $state('en'); // language variable to be used in the future when the app is translated into multiple languages, currently just set to english
	
	let term = $derived(items[counter]?.[lang]?.term || (isLoadingQuestions ? 'Loading questions...' : 'No questions available')); // create a local variable using the derived rune that dynamically updates as we move cards back and forth
	let def = $derived(items[counter]?.[lang]?.def || '');	// same as the term but for the definition
    let question_type = $derived(items[counter]?.question_type || ''); // same as the term but for the type (personality or activity)


	export function advanceCard( index = 0) { // update card information
		if (items.length === 0) {
			return;
		}

		if (index < 0 || index >= items.length) {
			return;
		}
		
		counter = index;
		term = items[counter][lang].term;
		def = items[counter][lang].def;
		console.log("tried to advance card to" + counter + ":" + term + ":" + def);
		
	}

	export function setDirection(dir = 0) { // change direction. This has to be an exported function so it can be accessed by other files, like the swipingApp.svelte
		directionInt = dir;
	}


	export function getTerm() { // because the list building is done in another file (SwipeingApp.svelte) we need a way to get the current term before changing it
		return term;
	}

	export function getTag() { // because the list building is done in another file (SwipeingApp.svelte) we need a way to get the current question_type before changing it
		return question_type;
	}

	export function getID() { // because the list building is done in another file (SwipeingApp.svelte) we need a way to get the current id before changing it
		return items[counter]?.id;
	}

	export function getListLength() {
		return items.length;
	}
	
	loadAdjectivesAndPersonalityTraits(); // load the adjectives and personality traits from the backend to be used in the cards or, if unable to load, use the hardcoded list of items

</script>




<div class="parent">
	{#if showServerErrorPopup}
		<div class="errorPopup" role="alert" aria-live="assertive">
			<p>Server could not be reached. Please try again later.</p>
			<button type="button" onclick={() => (showServerErrorPopup = false)}>Close</button>
		</div>
	{/if}


	<!-- first card which has to be animated differently from the rest of the stack -->
	{#key term} <!-- This key term allows for the cards to file in one after another and change properly without weird graphical issues -->
<div class="TopCard"
    in:fly={{ y: (100), duration: 1000}}
		out:fly={{ x: (directionInt)*(100), duration: 1000}} > <!-- -(negative) is to the left, + (positive) is to the right -->
    
    <h1 > {term} </h1> 
		<p > {def}</p> 

</div>
{/key}


	<!-- rest of the cards in a stack -->
    <div 
      class="card" 
    >
    </div>

	


</div>
<style>

 @media (prefers-color-scheme: dark) {
	:root {
		/* shared variables */
	  	--card-width: min(46vw, 42rem);
	--card-height: min(42vh, 24rem);
	--card--color: #1a1a1a;
	--text--color: white;
	--border-color: #ccc
	

	}
}

 @media (prefers-color-scheme: light) {
	:root {
		/* shared variables */
	  	--card-width: min(46vw, 42rem);
	--card-height: min(42vh, 24rem);
	--card--color: #ffffff;
	--text--color:#000000;
	--border-color: #000000;
	

	}
}

	h1 {
		
		font-size: clamp(1.5rem, 3vw, 3rem);
		margin: 10px;
		border-bottom: 3px solid var(--border-color);
		font-weight: 1000; /* make the text a larger weight so it stands out more*/
		color: var(--text--color);
		
	}

	p {
		font-size: clamp(1.2rem, 2.3vw, 2.2rem);
				margin: 10px;
						color: var(--text--color);

	}

	.parent {
		position: relative;
		display: grid;
		place-items: center;
		width: var(--card-width);
		height: var(--card-height);
	}

	.errorPopup {
		position: absolute;
		top: 1rem;
		left: 50%;
		transform: translateX(-50%);
		z-index: 2;
		width: min(90vw, 28rem);
		padding: 0.9rem 1rem;
		border: 2px solid var(--border-color);
		border-radius: 10px;
		background: var(--card--color);
		box-shadow: 0 0.4rem 1rem rgba(0, 0, 0, 0.2);
		text-align: center;
	}

	.errorPopup p {
		margin: 0 0 0.7rem 0;
		font-size: clamp(0.95rem, 2.2vw, 1.1rem);
		color: var(--text--color);
	}

	.errorPopup button {
		padding: 0.35rem 0.8rem;
		border: 2px solid var(--border-color);
		border-radius: 6px;
		background: transparent;
		color: var(--text--color);
		cursor: pointer;
	}
	
	.card:not(:first-child) {
		box-shadow: -3rem 0 3rem -2rem #000;
		color: transparent
	}
	
  .card {
		
		position: absolute; /* render following cards under the top card*/
		inset: 0;
		z-index: 0;

		
		transition: none; /* clear the transitions so that cards look normal until hovered over */

		
    grid-area: 1 / 1;  /* place the background cards in the same position as the top card so they look like they are under the first one */

		/* use shared variables so that all the cards are the same size*/
    width: var(--card-width); /* make the cards slightly smaller than the top card so that they look like they are under the top card */
    height: var(--card-height); /* make the cards slightly smaller than the top card so that they look like they are under the top card */
    background: var(--card--color);
		
    border: 4px solid var(--border-color); 
    border-radius: 8px;
    transition: transform 0.3s ease;
		transform: translate(0.5rem, 0.5rem);


  }

	.TopCard {
			
	  position: absolute;
		inset: 0;
		z-index: 1; /* setting this index to 1 allows it to render above the background stack*/

		transition: none;

	grid-area: 1 / 1;  /* place the background cards in the same position as the top card so they look like they are under the first one */

			
    width: var(--card-width);
    height: var(--card-height);
    background: var(--card--color);
			
    border: 4px solid var(--border-color);
    border-radius: 8px;
    transition: transform 0.3s ease;
			

		text-align: center;

			
    }

	.TopCard:hover {
		/* if the card is hovered over, it does a small rotation animation */

		transform: translateY(-0.45rem) rotate(2deg);
        

    }

	@media (max-width: 900px) {
		:root {
			--card-width: min(88vw, 30rem);
			--card-height: min(45vh, 22rem);
		}

		h1 {
			font-size: clamp(1.25rem, 6vw, 2rem);
		}

		p {
			font-size: clamp(1rem, 4.8vw, 1.5rem);
		}
	}

</style>