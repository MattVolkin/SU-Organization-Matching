<script>
	let barWidth = $state(0); // start from 0 so there is always some amount of progress made/some appearance (1-99 range)
	let progress = $state(0); // range of 0-100
	let {maxLimit = 100} = $props();

	let elmWidth = $state(0);
	let elmHeight = $state(0);

	export function advanceProgress(incrementAmt) {
		progress = incrementAmt;
		changeWidth();
	}

	function changeWidth() {
		if (progress)
		console.log(elmWidth);

		barWidth = (progress/(maxLimit-1))*100;
	}

	
</script>

<div class="progress-shell" bind:clientWidth={elmWidth} bind:clientHeight={elmHeight}>
	<!-- <h1>
		width height is : ({elmWidth}, {elmHeight})
	</h1>
	<h2>
		Progress is {progress}  and progress width is {barWidth}
	</h2> -->
	<div id="progBar" style="width: {barWidth}%" > </div>

</div>

<style>

#progBar {
		height: 100%;
		background-color: rgb(87, 247, 17);	
		border: none;
		transition: width 0.2s ease;
}

.progress-shell {
	width: 100%;
	max-width: min(56rem, 90vw);
	height: clamp(0.55rem, 1.25vh, 0.85rem);
	border: 2px solid #1C6EA4;
	border-radius: 999px;
	background: rgba(28, 110, 164, 0.2);
	overflow: hidden;
}
	
</style>