<script>

// TODO make it look like its somthing that can be swipable
	// put items on a stack of cards rather than just box on a page
	
  	import { blur, fade, fly, scale, slide } from 'svelte/transition'; //import transition animations
  	import { APICreater } from './APIHandler.svelte'; // import the API handler to make calls to the backend and get the data for the cards

	let counter = $state(0); // count what term we are on to stay on the list
	let directionInt = $state('none'); // gets direction information from swiping app

	let Hardcodeditems = $state([ // terms that will be included in the cards shown to the user

{id:0, tag: "activities", term: "Community Service/Fundraising", def: "Voluntary work/raising funds intended to help people in a particular area."},
{id:1, tag: "activities", term: "Social Justice", def: "Advocating for the fair treatment and equitable status of all individuals and social groups within a society"},
{id:2, tag: "activities", term: "Retreats", def: "Planned off-campus getaway or sleepover for organization members designed to bond, foster unity, and align on organizational goals"},
{id:3, tag: "activities", term: "Dance", def: "Moving rhythmically to music, typically following a set sequence of steps"},
{id:4, tag: "activities", term: "Board Games", def: "Playing games that involves the movement of counters or other pieces on a marked board, often with the use of other components such as dice or cards"},
{id:5, tag: "activities", term: "Movies", def: "Watching a story or event recorded by a camera as a set of moving images and shown in a theater or on television; a motion picture"},
{id:6, tag: "activities", term: "Video Games", def: "Playing games by electronically manipulating images produced by a computer program on a television screen or other display screen."},
{id:7, tag: "activities", term: "Arts & Crafts", def: "Making objects, such as decorations, toys, furniture, and pottery by hand"},
{id:8, tag: "activities", term: "Music", def: ""},
{id:9, tag: "activities", term: "Exercise", def: ""},
{id:10, tag: "activities", term: "Writing", def: ""},
{id:11, tag: "activities", term: "Professional Development", def: "Help for finding future jobs/educational opperunities"},
{id:12, tag: "activities", term: "Caring for Animals", def: ""},
{id:13, tag: "activities", term: "Giving Presentations", def: ""},
{id:14, tag: "activities", term: "Trivia", def: ""},
{id:15, tag: "activities", term: "Literary Analysis", def: ""},
{id:16, tag: "activities", term: "Study Groups", def: ""},
{id:17, tag: "activities", term: "Guest Speakers", def: ""},
{id:18, tag: "activities", term: "Group Lunch/Dinner", def: ""},
{id:19, tag: "activities", term: "Discussion", def: ""},
{id:20, tag: "personality", term: "Welcoming", def: "to greet hospitably and with courtesy or cordiality"},
{id:21, tag: "personality", term: "Hard Working", def: "constantly, regularly, or habitually engaged in earnest and energetic work"},
{id:22, tag: "personality", term: "Caring", def: "feeling or showing concern for or kindness to others"},
{id:23, tag: "personality", term: "Creative", def: "the ability or power to make somthing"},
{id:24, tag: "personality", term: "Outgoing", def: "openly friendly and responsive"},
{id:25, tag: "personality", term: "Open Minded", def: "receptive to arguments or ideas"},
{id:26, tag: "personality", term: "Eager to Learn", def: "someone with a strong, enthusiastic desire to acquire new knowledge, skills, or experiences"},
{id:27, tag: "personality", term: "Confident", def: "having or showing assurance and self-reliance"},
{id:28, tag: "personality", term: "Nerdy", def: "a person devoted to intellectual, academic, or technical pursuits or interests"},
{id:29, tag: "personality", term: "Leader", def: "a person who guides "},
{id:30, tag: "personality", term: "Enthusiastic", def: "filled with or marked by strong excitement of feeling"},
{id:31, tag: "personality", term: "Collaborative", def: "involving or done by two or more people or groups working together"},
{id:32, tag: "personality", term: "Curious", def: "desire to investigate and learn"},
{id:33, tag: "personality", term: "Organized", def: "arranged in a systematic way, especially on a large scale."},
{id:34, tag: "personality", term: "Social", def: "pleasant companionship with friends or associates"},
{id:35, tag: "personality", term: "Fun", def: "providing entertainment, amusement, or enjoyment"}

	]);



	async function loadAdjectivesAndPersonalityTraits() { // load the adjectives and personality traits from the backend to be used in the cards or, if unable to load, use the hardcoded list of items
    try {
      const adjectivesList = await APICreater('GET', '/api/adjectives', null, null);
	  const personalityTraitsList = await APICreater('GET', '/api/personality_traits', null, null);
	  items = Hardcodeditems | [...adjectivesList, ...personalityTraitsList]; // combine the two lists into one list of items for the cards

	} catch (error) {
      console.error('Unable to load officer clubs', error);
    }

}

	const range = (start, end, targetArray) => { // range function to get a sub-array of the full question array (without modifying the original)
    let arr = [];
    for (let i = start; i < end; i += 1) arr.push(targetArray[i]);
    return arr;
  };

	
	let term = $derived(items[counter].term); // create a local variable using the derived rune that dynamically updates as we move cards back and forth
	let def = $derived(items[counter].def);	// same as the term but for the definition
    let tag = $derived(items[counter].tag); // same as the term but for the type (personality or activity)

	
	export function advanceCard( index = 0) { // update card information
		
		counter = index;
		term = items[counter].term;
		def = items[counter].def;
		console.log("tried to advance card" + counter + ":" + term + ":" + def);

		//TODO: if at end of list and finished with swiping, do something (undefined at this time)
		
	}

	export function setDirection(dir = 0) { // change direction. This has to be an exported function so it can be accessed by other files, like the swipingApp.svelte
		directionInt = dir;
	}


	export function getTerm() { // because the list building is done in another file (SwipeingApp.svelte) we need a way to get the current term before changing it
		return term;
	}

	export function getTag() { // because the list building is done in another file (SwipeingApp.svelte) we need a way to get the current tag before changing it
		return tag;
	}

	export function getListLength() {
		return items.length;
	}
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
	--card-height: 60vw;
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

		
		transition: none; /* clear the transitions so that cards look normal until hovered over */

		
    grid-area: 1 / 1;  /* place the background cards in the same position as the top card so they look like they are under the first one */

		/* use shared variables so that all the cards are the same size*/
    width: var(--card-width);
    height: var(--card-height);
    background: var(--card--color);
		
    border: 5px solid #ccc; /* offset bottom cards so the cards look like they are in a stack*/
    border-radius: 8px;
    transition: transform 0.3s ease;
		margin-bottom: 160px;


  }

    .TopCard {
			
	  position: absolute;
		z-index: 1; /* setting this index to 1 allows it to render above the background stack*/

		transition: none;
			
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