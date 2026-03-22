<script>
  import { fade, fly, scale } from 'svelte/transition'; //import transition animations

  let { show = true } = $props();
	let counter = $state(0); // count what term we are on to stay on the list
	let directionInt = $state('none'); // gets direction information from swiping app
	
	// list of items
	//https://svelte.dev/playground/805300f5895f4ea89b73ba75de393db8?version=5.53.6
	let items = $state([ // terms that will be included in the cards shown to the user
		{id: 0, term: "Likable", def: "someone who is an enjoyable person to be around"},
		{id: 1, term: "Tired", def: "someone who has a lack of sleep"},
		{id: 2, term: "Happy", def: "State of euphoria"}

	]);

	
	let term = $derived(items[counter].term); // create a local variable using the derived rune that dynamically updates as we move cards back and forth
	let def = $derived(items[counter].def);	// same as the term but for the definition

	export function advanceCard( index = 0) { // update card information
		
		counter = index;
		counter %= items.length
		term = items[counter].term;
		def = items[counter].def;
		console.log("tried to advance card" + counter + ":" + term + ":" + def);
		
	}

	export function setDirection(dir = 0) { // change direction. This has to be an exported function so it can be accessed by other files, like the swipingApp.svelte
		directionInt = dir;
	}


	export function getTerm() { // because the list buildijg is done in another file (SwipeingApp.svelte) we need a way to get the current term before changing it
		return term;
	}

</script>




<button onclick={() => show = !show}>
  Toggle Elements
</button>

{#key term} <!-- This key term allows for the cards to file in one after another and change propperly without weird graphical issues -->
<div 
		out:fly={{ x: (directionInt)*(100), duration: 1000}} 
		onoutroend={() => console.log('tried to update contents of text')}> <!-- -(negative) is to the left, + (positive) is to the right -->
    <h1> {term} </h1>
		<p> {def}</p>

</div>
{/key}

<div> adding this for spacing reasons: {show}</div>

<style>
	.content {
		user-select: none;
	}

	div {
		border: 1px solid;
		font-size: 12px;
		color: red;
	}

	.box {
		border: 1px solid;
		padding:0.5rem;
		height:50%;

	}
</style>