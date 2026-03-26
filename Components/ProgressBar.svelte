<script>
	let barWidth = $state(0); // start from 0 so there is always some amount of progress made/some appearance (1-99 range)
	let progress = $state(0); // range of 0-100
	let {maxLimit = 100} = $props();

	let elmWidth = $state(0);
	let elmHeight = $state(0);

	export function advanceProgress(incrementAmt) {
		progress += incrementAmt;
		changeWidth();
	}

	function changeWidth() {
		if (progress)
		console.log(elmWidth);

		barWidth = (progress/maxLimit)*100;
	}

	
</script>

<div bind:clientWidth={elmWidth} bind:clientHeight={elmHeight}>
	<!-- <h1>
		width height is : ({elmWidth}, {elmHeight})
	</h1>
	<h2>
		Progress is {progress}  and progress width is {barWidth}
	</h2> -->
	<div id="progBar" style="width: {barWidth}%" > </div>

<button onclick={() => advanceProgress(1)}>Run JavaScript</button>
</div>

<style>

#progBar {
		height: 30px;
		background-color: rgb(87, 247, 17);	
		border: 1px solid rgb(87, 247, 17);
}

div {
	      border: 5px solid;
				background: #1C6EA4;
}
	
</style>