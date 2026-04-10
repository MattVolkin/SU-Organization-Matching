<script lang='ts'>
	let name = $state('world');
	let count = $state(0);
	// https://svelte.dev/playground/f696ca27e6374f2cab1691727409a31d?version=5.53.2
	// https://svelte.dev/docs/svelte/animate

	import { useSwipe , type SwipeCustomEvent } from 'svelte-gestures'; // import methods and events from Svelte-Gesture library
	import Card from './Card.svelte'; // import the Card Component and it's relevant method (advanceCard) to display the details of a word and its description
	import Bar from './ProgressBar.svelte'
	import API from './APIHandler.svelte' // import the API handler to make calls to the backend and get the data for the cards



	let cardObject = $state({}); //create this object as a blank version of the Card class so that we can refer to one object that can internally track its data
	let barObject = $state({}); // create object so we can call methods from the bar file
	
  import { fade, fly, scale } from 'svelte/transition'; // import transitions

  let left = $state(true); // if we swipe left or right, augment given booleans to make sure the card is able to come back into frame/know what direction to move card
	let right = $state(true);
 
	let activities = $state(['']) // track the activities that the user likes doing to be fed into the fitness function
	let personality = $state(['']) // track personality traits that the user likes doing (to feed into fitness function)

	let direction = $state("none yet"); // track direction swiping goes to determine which direction elements move
	let directionInt = $state(-1); // represent direction as an integer for calculations of speed and position
	let pointerType: string; // debug information of what type (mouse or touch control)

	let mx: number; // x position of cursor
	let my: number; // y position of cursor

	let target: HTMLElement | null; // what kind of element is being interacted with as well as 

	function swipeLeft() { // method to swipe left, this is broken out into a function so that it can be called by both the swipe gesture and arrow keys
			
		//animation code
		left = !left;
			directionInt = -1;
			setTimeout(() => {
				left = !left;
												} , 3000);
			console.log(typeof cardObject);
			count += 1;
			cardObject.setDirection(directionInt)
			cardObject.advanceCard(count);
				barObject.advanceProgress(count)

	}

	function swipeRight() {// method to swipe right, this is broken out into a function so that it can be called by both the swipe gesture and arrow keys
		
			if(cardObject.getTag() === 'activities') { // if the card that was swiped right is an activity, add it to the activities list, otherwise add it to the personality list
				activities = [...activities, cardObject.getTerm()];
			} else {
				personality = [...personality, cardObject.getTerm()];
			}

			// animation code
			right = !right;
			directionInt = 1;
			setTimeout(() => {
				right = !right;
												} , 3000);
			console.log(typeof cardObject);
			count += 1;
			cardObject.setDirection(directionInt)
			cardObject.advanceCard(count);
			barObject.advanceProgress(count)		

	}

	function rewind() { // method to go back to the previous card, this is called when the user clicks the undo button
		if(count > 0) { // prevent user from setting count to invalid values/trying to rewind past the stack of terms

			// same animation code as swipeLeft
		left = !left;
			directionInt = -1;
			setTimeout(() => {
				left = !left;
												} , 3000);
			console.log(typeof cardObject);

		count -= 1;
		cardObject.setDirection(directionInt);
		cardObject.advanceCard(count);
		barObject.advanceProgress(count)

			// because we are trying to go to the previous term, we should remove it from the array (if it exists), we should 
		deleteGivenCardFromUserArray()
		console.log("Made it out of remove function with: " + cardObject.getTerm())

			
		}
	}

	function deleteGivenCardFromUserArray(string: GivenItem) { // remove last term in array if it is the card we rewinded to
	let GivenItem = cardObject.getTerm();
		
		console.log("Made it into remove function with: " + GivenItem);

		console.log("does " + cardObject.getTerm() + " match " + GivenItem)
			if(cardObject.getTag() === 'activities') { // if the card that was added to the array is an activity, then check the activities list to check if we need to remove it
					if(activities[activities.length-1] === GivenItem) {
						activities.pop();
						console.log("removed " + GivenItem);
					}
			
			
			} else {
					if(personality[personality.length-1] === GivenItem) {
						personality.pop();
						console.log("removed " + GivenItem);
					}
			}
		
	}


	function swipeHandler(event: SwipeCustomEvent) { // method that captures the swipe motion and records the given data, as well as performs some variable assignments to help with the motion of the animations
		direction = event.detail.direction; 
		pointerType = event.detail.pointerType;
		target = event.detail.target as HTMLElement;
		if (direction === 'left') { // we want to make sure that when the user swipes left, the animation of the card moves in the corresponding direction
			swipeLeft();
		}
		else if (direction === 'right') { // move card right when mouse is swiped right

			swipeRight();
		}

		else if (direction === 'top') { // DEBUG; when we swipe up print the activities list to the console
			console.log(activities);
		}
		else if (direction === 'bottom') { // DEBUG; when we swipe down print the personality list to the console
			console.log(personality);
		}




		
	}

	function keyboardHandler(event: KeyboardEvent) { // method to capture the left and right arrow keys to perform the same function as swiping left and right, this is for accessibility reasons
		console.log(event.key);

		if (event.key === 'ArrowLeft') {
			swipeLeft();
		}
		else if (event.key === 'ArrowRight') {
			swipeRight();
		}
	}

	function moveHandler(event: GestureCustomEvent) { // get the x y value of mouse
		mx = event.detail.x;
		my = event.detail.y;
	}
	
</script>


<!-- 
<h1>Hello {name}!</h1>

<input bind:value={name} />
<button onclick={() => count += 1}>
	clicks: {count}
</button>
 -->

<svelte:window onkeydown={keyboardHandler} />


<div class="swipedContainer">
	<h1>Swipe right if you like the activity,</h1>
	<h1>Swipe left if you don't.</h1>
	<!-- Debug information -->
	<!-- <h2>Activities you like: {activities} </h2>
	<h2>Personality traits you like: {personality} </h2>
	<h2> {count} </h2> -->
	
	<Bar maxLimit = 37 bind:this={barObject} class="bar"/> <!-- create the progress bar and bind it to the barObject so that we can call methods from the progress bar file -->



<section
	{...useSwipe(swipeHandler, () => ({ timeframe: 300, minSwipeDistance: 25, touchAction: 'none' }), { // bound swipe function to work given a specific bounds of a section element
		onswipemove: moveHandler
	})}
	class="box"
>
		<button onclick={() => rewind()}>Undo</button> <!-- Button to go back to the previous card -->


	<div class="content">
					<Card bind:this={cardObject} />						
	</div>


	
</section>	
	

</div>


<style>

	h1 {
		margin: 10px;
		margin-bottom: 1px;
		text-align: center;
		font-size: 4vw;
		color: white;
	}

	button {
		align-self: flex-start;
		height: 8vw;
		width: 10vw;
		font-size: 3vw;
		margin: 10px;
		background-color:  #1a1a1a;
		color:white;
	}
	
	.content {
		user-select: none;
		align-items: center;
		justify-content: center;
		
		display: grid;
  place-items: center; /* Centers both horizontally and vertically */
  height: 100vh;       /* Example height to fill the viewport */
	}


 @media (prefers-color-scheme: dark) {
	.swipedContainer {
		
	background: #ad0000;
	background: linear-gradient(90deg,rgba(173, 0, 0, 1) 0%, rgba(0, 0, 0, 255) 33%, rgba(0, 0, 0, 255) 66%, rgba(0, 168, 67, 1) 100%);		
		flex-direction: column;
		align-items: center;
		justify-content: center;


		overflow: hidden;

	}
}

</style>