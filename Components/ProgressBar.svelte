<!-- @component Displays a visual progress bar that fills based on quiz completion percentage. -->
<script>
/**
 * @type {state} barWidth - Calculated percentage width of the progress bar (0-100), updated when progress changes
 * @type {state} progress - Current progress value tracked from quiz completion (0-100 range)
 * @type {state} elmWidth - Current width of the progress bar container element in pixels
 * @type {state} elmHeight - Current height of the progress bar container element in pixels
 * @type {props} maxLimit - Maximum progress value for calculating percentage, defaults to 100
 * @function advanceProgress - Updates the progress value and recalculates the bar width for animation
 * @function changeWidth - Calculates the new bar width percentage based on current progress and max limit
 */
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