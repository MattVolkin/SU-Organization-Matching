<script>
  import { fade, fly, scale } from 'svelte/transition';

  let { show = true } = $props();
	let counter = $state(0);
	let directionInt = $state('none');
	
	// list of items
	//https://svelte.dev/playground/805300f5895f4ea89b73ba75de393db8?version=5.53.6
	let items = $state([
		{id: 0, term: "Likable", def: "someone who is an enjoyable person to be around"},
		{id: 1, term: "Tired", def: "someone who has a lack of sleep"},
		{id: 2, term: "Happy", def: "State of euphoria"}

	]);

	
	let term = $derived(items[counter].term);
	let def = $derived(items[counter].def);	

	export function advanceCard( index = 0) {
		
		counter = index;
		counter %= items.length
		term = items[counter].term;
		def = items[counter].def;
		console.log("tried to advance card" + counter + ":" + term + ":" + def);
		
	}

	export function setDirection(dir = 0) {
		directionInt = dir;
	}


	export function getTerm() {
		return term;
	}

</script>




<button onclick={() => show = !show}>
  Toggle Elements
</button>

{#key term}
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