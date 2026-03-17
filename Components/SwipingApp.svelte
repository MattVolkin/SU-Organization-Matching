<script lang='ts'>
	let name = $state('world');
	let count = $state(0);
	// https://svelte.dev/playground/f696ca27e6374f2cab1691727409a31d?version=5.53.2
	// https://svelte.dev/docs/svelte/animate

	import { useSwipe , type SwipeCustomEvent, tyoeGestureCustomEvent } from 'svelte-gestures';
	import Card from './Card.svelte';
	import { advanceCard } from './Card.svelte'

	let cardObject = $state({});
	
  import { fade, fly, scale } from 'svelte/transition';
  let left = $state(true);
	let right = $state(true);

	let activities = ['testActivities']
	let personality = ['testPersonalities']

	let direction = $state("none yet");
	let directionInt = $state(-1);
	let pointerType: string;

	let mx: number;
	let my: number;

	let target: HTMLElement | null;


	function swipeHandler(event: SwipeCustomEvent) {
		direction = event.detail.direction;
		pointerType = event.detail.pointerType;
		target = event.detail.target as HTMLElement;
		if (direction === 'left') {
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
		else if (direction === 'right') {

			activities = [...activities, cardObject.getTerm()];
			
			
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

		else if (direction === 'top') {
			console.log(activities);
		}
		
	}

	

	function moveHandler(event: GestureCustomEvent) {
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