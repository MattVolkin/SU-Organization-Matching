<script lang='ts'>
	let name = $state('world');
	let count = $state(0);
	// https://svelte.dev/playground/f696ca27e6374f2cab1691727409a31d?version=5.53.2
	// https://svelte.dev/docs/svelte/animate

	import { useSwipe , type SwipeCustomEvent, typeGestureCustomEvent } from 'svelte-gestures'; // import methods and events from Svelte-Gesture library
	import Card from './Card.svelte'; // import the Card Component and it's relevant method (advanceCard) to display the details of a word and its description
	import { advanceCard } from './Card.svelte'

	let cardObject = $state({}); //create this object as a blank version of the Card class so that we can refer to one object that can internally track its data
	
  import { fade, fly, scale } from 'svelte/transition'; // import transitions

  let left = $state(true); // if we swipe left or right, augment given booleans to make sure the card is able to come back into frame/know what direction to move card
	let right = $state(true);

	let activities = ['testActivities'] // track the activities that the user likes doing to be fed into the fitness function
	let personality = ['testPersonalities'] // track personality traits that the user likes doing (to feed into fitness function)

	let direction = $state("none yet"); // track direction swiping goes to determine which direction elements move
	let directionInt = $state(-1); // represent direction as an integer for calculations of speed and position
	let pointerType: string; // debug information of what type (mouse or touch control)

	let mx: number; // x position of cursor
	let my: number; // y position of cursor

	let target: HTMLElement | null; // what kind of element is being interacted with as well as 

	function swipeHandler(event: SwipeCustomEvent) { // method that captures the swipe motion and records the given data, as well as performs some variable assignments to help with the motion of the animations
		direction = event.detail.direction; 
		pointerType = event.detail.pointerType;
		target = event.detail.target as HTMLElement;
		if (direction === 'left') { // we want to make sure that when the user swipes left, the animation of the card moves in the corresponding direction
			left = !left;
			directionInt = -1;
			setTimeout(() => {
				left = !left;
												} , 3000);
			console.log(typeof cardObject);
			count += 1;
			cardObject.setDirection(directionInt)
			cardObject.advanceCard(count);
		}
		else if (direction === 'right') { // move card right when mouse is swiped right

			activities = [...activities, cardObject.getTerm()]; // because right is considered an accept/agree, we add the term to our list
			
			
			right = !right;
			directionInt = 1;
			setTimeout(() => {
				right = !right;
												} , 3000);
			console.log(typeof cardObject);
			count += 1;
			cardObject.setDirection(directionInt)
			cardObject.advanceCard(count);
		}

		else if (direction === 'top') { // DEBUG; when we swipe up print the activities list to the console
			console.log(activities);
		}
		
	}

	

	function moveHandler(event: GestureCustomEvent) { // get the x y value of mouse
		mx = event.detail.x;
		my = event.detail.y;
	}
	
</script>



<h1>Hello {name}!</h1>

<input bind:value={name} />
<button onclick={() => count += 1}>
	clicks: {count}
</button>


<section
	{...useSwipe(swipeHandler, () => ({ timeframe: 300, minSwipeDistance: 25, touchAction: 'none' }), {
		onswipemove: moveHandler
	})}
	class="box"
>
	<div class="content">
		<h2>Swipe Handler</h2>
		<div>swipe direction: {direction}</div>
		<div>pointerType {pointerType}</div>
		<div>target: {target?.tagName}</div>
		<h2>move</h2>
		<div>x: {mx}</div>
		<div>y: {my}</div>


		<!-- https://www.youtube.com/watch?v=5oEo98BrRqs -->


					<Card bind:this={cardObject} />				

		
	</div>



	
</section>	

<!--
{#if right}
	<div transition:fly={{ x: 100, duration: 700 }}>
		whoopdedoo
	</div>
{/if}
-->

<h2> lets see how the animationBooleanWorks: {left} : {right} : {directionInt} </h2>
<h1>Transition Demo</h1>
<Card left/>


<style>
	.content {
		user-select: none;
	}

	.box {
		border: 1px solid;
		padding:0.5rem;
		height: 75%;
	  width: 75%;


	}

	
</style>