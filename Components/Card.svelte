<script>

// TODO make it look like its somthing that can be swipable
	// put items on a stack of cards rather than just box on a page
	
  	import { blur, fade, fly, scale, slide } from 'svelte/transition'; //import transition animations
  	import { APICreater } from './APIHandler.svelte'; // import the API handler to make calls to the backend and get the data for the cards

	let counter = $state(0); // count what term we are on to stay on the list
	let directionInt = $state('none'); // gets direction information from swiping app

	let Hardcodeditems = $state([ // terms that will be included in the cards shown to the user

{id:0, question_type: "activities", en: {term: "Community Service/Fundraising", def: "Voluntary work/raising funds intended to help people in a particular area."}},
{id:1, question_type: "activities", en: {term: "Social Justice", def: "Advocating for the fair treatment and equitable status of all individuals and social groups within a society"}},
{id:2, question_type: "activities", en: {term: "Retreats", def: "Planned off-campus getaway or sleepover for organization members designed to bond, foster unity, and align on organizational goals"}},
{id:3, question_type: "activities", en: {term: "Dance", def: "Moving rhythmically to music, typically following a set sequence of steps"}},
{id:4, question_type: "activities", en: {term: "Board Games", def: "Playing games that involves the movement of counters or other pieces on a marked board, often with the use of other components such as dice or cards"}},
{id:5, question_type: "activities", en: {term: "Movies", def: "Watching a story or event recorded by a camera as a set of moving images and shown in a theater or on television; a motion picture"}},
{id:6, question_type: "activities", en: {term: "Video Games", def: "Playing games by electronically manipulating images produced by a computer program on a television screen or other display screen."}},
{id:7, question_type: "activities", en: {term: "Arts & Crafts", def: "Making objects, such as decorations, toys, furniture, and pottery by hand"}},
{id:8, question_type: "activities", en: {term: "Music", def: ""}},
{id:9, question_type: "activities", en: {term: "Exercise", def: ""}},
{id:10, question_type: "activities", en: {term: "Writing", def: ""}},
{id:11, question_type: "activities", en: {term: "Professional Development", def: "Help for finding future jobs/educational opperunities"}},
{id:12, question_type: "activities", en: {term: "Caring for Animals", def: ""}},
{id:13, question_type: "activities", en: {term: "Giving Presentations", def: ""}},
{id:14, question_type: "activities", en: {term: "Trivia", def: "A game involving questions about various subjects"}},
{id:15, question_type: "activities", en: {term: "Literary Analysis", def: ""}},
{id:16, question_type: "activities", en: {term: "Study Groups", def: ""}},
{id:17, question_type: "activities", en: {term: "Guest Speakers", def: ""}},
{id:18, question_type: "activities", en: {term: "Group Lunch/Dinner", def: ""}},
{id:19, question_type: "activities", en: {term: "Discussion", def: ""}},
{id:20, question_type: "personality", en: {term: "Welcoming", def: "to greet hospitably and with courtesy or cordiality"}},
{id:21, question_type: "personality", en: {term: "Hard Working", def: "constantly, regularly, or habitually engaged in earnest and energetic work"}},
{id:22, question_type: "personality", en: {term: "Caring", def: "feeling or showing concern for or kindness to others"}},
{id:23, question_type: "personality", en: {term: "Creative", def: "the ability or power to make somthing"}},
{id:24, question_type: "personality", en: {term: "Outgoing", def: "openly friendly and responsive"}},
{id:25, question_type: "personality", en: {term: "Open Minded", def: "receptive to arguments or ideas"}},
{id:26, question_type: "personality", en: {term: "Eager to Learn", def: "someone with a strong, enthusiastic desire to acquire new knowledge, skills, or experiences"}},
{id:27, question_type: "personality", en: {term: "Confident", def: "having or showing assurance and self-reliance"}},
{id:28, question_type: "personality", en: {term: "Nerdy", def: "a person devoted to intellectual, academic, or technical pursuits or interests"}},
{id:29, question_type: "personality", en: {term: "Leader", def: "a person who guides "}},
{id:30, question_type: "personality", en: {term: "Enthusiastic", def: "filled with or marked by strong excitement of feeling"}},
{id:31, question_type: "personality", en: {term: "Collaborative", def: "involving or done by two or more people or groups working together"}},
{id:32, question_type: "personality", en: {term: "Curious", def: "desire to investigate and learn"}},
{id:33, question_type: "personality", en: {term: "Organized", def: "arranged in a systematic way, especially on a large scale."}},
{id:34, question_type: "personality", en: {term: "Social", def: "pleasant companionship with friends or associates"}},
{id:35, question_type: "personality", en: {term: "Fun", def: "providing entertainment, amusement, or enjoyment"}}

	]);

	let items = $state(Hardcodeditems); // active card list used by the component at runtime

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
				const englishValue = typeof item.en === 'string'
					? item.en
					: (item.en?.term || item.en?.def || '');

				if (!englishValue) {
					return null;
				}

				return {
					id: item.id,
					question_type: questionType,
					en: {
						term: englishValue,
						def: typeof item.en?.def === 'string' ? item.en.def : '',
					},
				};
			})
			.filter(Boolean);
	}



	async function loadAdjectivesAndPersonalityTraits() { // load the adjectives and personality traits from the backend to be used in the cards or, if unable to load, use the hardcoded list of items
    try {
      const adjectivesList = await APICreater('GET', '/api/adjectives', null);
	  const normalizedItems = normalizeApiCardItems(adjectivesList);
	  items = normalizedItems.length > 0 ? normalizedItems : Hardcodeditems;

	} catch (error) {
      console.error('Unable to load adjectives and personality traits', error);
	  items = Hardcodeditems;
    }
}

	async function sendSwipeInformation() { // after finishing all of the cards, send the information about the activities and personality traits that the user likes to the backend to be used in the fitness function and matching algorithm
    try {
      const adjectivesList = await APICreater('GET', '/api/adjectives', null);

	} catch (error) {
      console.error('Unable to load officer clubs', error);
    }

}

	const range = (start, end, targetArray) => { // range function to get a sub-array of the full question array (without modifying the original)
    let arr = [];
    for (let i = start; i < end; i += 1) arr.push(targetArray[i]);
    return arr;
  };

  let lang = $state('en'); // language variable to be used in the future when the app is translated into multiple languages, currently just set to english
	
	let term = $derived(items[counter][lang].term); // create a local variable using the derived rune that dynamically updates as we move cards back and forth
	let def = $derived(items[counter][lang].def);	// same as the term but for the definition
    let question_type = $derived(items[counter].question_type); // same as the term but for the type (personality or activity)


	export function advanceCard( index = 0) { // update card information
		
		counter = index;
		term = items[counter][lang].term;
		def = items[counter][lang].def;
		console.log("tried to advance card" + counter + ":" + term + ":" + def);

		//TODO: if at end of list and finished with swiping, do something (undefined at this time)
		
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

	export function getListLength() {
		return items.length;
	}
	
	loadAdjectivesAndPersonalityTraits(); // load the adjectives and personality traits from the backend to be used in the cards or, if unable to load, use the hardcoded list of items

</script>




<div class="parent">


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
	{#each range(counter+1, items.length, items) as card, index} 
    <div 
      class="card" 
    >
    </div>
	{/each}

	


</div>
<style>

	:root {
		/* shared variables */
  	--card-width: 50vw;
	--card-height: 60vh;
	--card--color: #1a1a1a;

	}

	h1 {
		
		font-size: 4vw;
		margin: 10px;
		border-bottom: 3px solid #ccc;
		font-weight: 1000; /* make the text a larger weight so it stands out more*/
		
	}

	p {
		font-size: 3vw;
				margin: 10px;
	}

	.parent {
		/* althougth there is nothing in the styling for the parent element, this is here just incase somthing needs to be added*/
	}
	
	.card:not(:first-child) {
		box-shadow: -3rem 0 3rem -2rem #000;
		color: transparent
	}
	
  .card {
		
		position: relative; /* render following cards under the top card*/
		z-index: 0;
;

		
		transition: none; /* clear the transitions so that cards look normal until hovered over */

		
    grid-area: 1 / 1;  /* place the background cards in the same position as the top card so they look like they are under the first one */

		/* use shared variables so that all the cards are the same size*/
    width: calc(var(--card-width) - 5vw); /* make the cards slightly smaller than the top card so that they look like they are under the top card */
    height: calc(var(--card-height) - 5vh); /* make the cards slightly smaller than the top card so that they look like they are under the top card */
    background: var(--card--color);
		
    border: 2px solid #ccc; 
    border-radius: 8px;
    transition: transform 0.3s ease;
		margin-bottom: 160px;


  }

    .TopCard {
			
	  position: absolute;
		z-index: 1; /* setting this index to 1 allows it to render above the background stack*/

		transition: none;

	grid-area: 1 / 1;  /* place the background cards in the same position as the top card so they look like they are under the first one */

			
    width: var(--card-width);
    height: var(--card-height);
    background: var(--card--color);
			
    border: 2px solid #ccc;
    border-radius: 8px;
    transition: transform 0.3s ease;
			

		text-align: center;

			
    }

	.TopCard:hover {
		/* if the card is hovered over, it does a small rotation animation */

		transform: translateY(-1rem) rotate(3deg);
        

    }

</style>