<script>

// TODO make it look like its somthing that can be swipable
// TODO make it look pretty
// TODO centered on page
// TODO make general progress bar (new file that can be added for other parts)
// TODO let arrowkeys work for swiping
  import { blur, fade, fly, scale, slide } from 'svelte/transition'; //import transition animations

  let { show = true } = $props();
	let counter = $state(0); // count what term we are on to stay on the list
	let directionInt = $state('none'); // gets direction information from swiping app
	
	// list of items
	//https://svelte.dev/playground/805300f5895f4ea89b73ba75de393db8?version=5.53.6
	let items = $state([ // terms that will be included in the cards shown to the user

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

	
	let term = $derived(items[counter].term); // create a local variable using the derived rune that dynamically updates as we move cards back and forth
	let def = $derived(items[counter].def);	// same as the term but for the definition
    let tag = $derived(items[counter].tag); // same as the term but for the type (personality or activity)
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


	export function getTerm() { // because the list building is done in another file (SwipeingApp.svelte) we need a way to get the current term before changing it
		return term;
	}

	export function getTag() { // because the list building is done in another file (SwipeingApp.svelte) we need a way to get the current tag before changing it
		return tag;
	}
</script>



<div>
<!-- <button onclick={() => show = !show}>
  Toggle Elements
</button> -->

{#key term} <!-- This key term allows for the cards to file in one after another and change propperly without weird graphical issues -->
<div class="Card"
        in:fade = {{delay: 1000, duration: 1000}}
		out:fly={{ x: (directionInt)*(100), duration: 1000}} > <!-- -(negative) is to the left, + (positive) is to the right -->
    
    <h1 in:fade = {{delay: 1000, duration: 1000}}> {term} </h1>
		<p in:fade={{delay:1000, speed: 1000 }}> {def}</p>

</div>
{/key}


</div>
<style>
    .Card {

        background: #1C6EA4;
        background: -moz-linear-gradient(-45deg, #1C6EA4 0%, #144E75 100%);
        background: -webkit-linear-gradient(-45deg, #1C6EA4 0%, #144E75 100%);
        background: linear-gradient(135deg, #1C6EA4 0%, #144E75 100%);

        -webkit-box-shadow: 0px 5px 15px 0px #000000; 
        box-shadow: 0px 5px 15px 0px #000000;

        border: 1px solid;


        color: #91ff00;
        width: 100vw;
        height: 30vw;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;
        padding: 1.5rem;
        box-sizing: border-box;
    }

</style>
<!-- 
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
</style> -->