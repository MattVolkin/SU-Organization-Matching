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




<!-- <button onclick={() => show = !show}>
  Toggle Elements
</button> -->

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
    .admin-home {
        --bg: linear-gradient(180deg, #edf4fb 0%, #f6f9fd 100%);
        --card: #ffffff;
        --text: #132c45;
        --muted: #4f6781;
        --border: #d4e0ec;
        --accent: #0f6d8c;
        --accent-hover: #0b5972;
        --danger: #a5294a;
        --danger-hover: #8a1f3d;

        background: var(--bg);
        color: var(--text);
        min-height: 100vh;
        padding-bottom: 0.75rem;
    }

    .admin-home h1 {
        width: min(100%, 1040px);
        margin: 1rem auto 0.6rem auto;
        padding: 0 1rem;
        font-size: clamp(1.35rem, 1.5vw + 0.9rem, 2rem);
        line-height: 1.2;
    }

    .club-management {
        width: min(100%, 1040px);
        margin: 0 auto 1rem auto;
        padding: 1rem;
        border: 1px solid var(--border);
        border-radius: 1rem;
        background: var(--card);
        box-shadow: 0 10px 24px rgba(13, 37, 62, 0.1);
    }

    .club-management > p {
        margin: 0 0 0.8rem 0;
        color: var(--muted);
        line-height: 1.45;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        border: 1px solid var(--border);
        border-radius: 0.8rem;
        overflow: hidden;
        background: #fff;
    }

    th,
    td {
        text-align: left;
        padding: 0.85rem 1.1rem;
        border-bottom: 1px solid #e8eef4;
    }

    .col-name {
        width: 68%;
    }

    .col-actions {
        width: 32%;
    }

    th {
        background: #f5f9ff;
        text-transform: uppercase;
        letter-spacing: 0.03em;
        font-size: 0.8rem;
        color: #2f4a66;
    }

    th:nth-child(2) {
        text-align: right;
    }

    tbody tr:hover {
        background: #f9fcff;
    }

    td:nth-child(2) {
        display: flex;
        gap: 0.65rem;
        flex-wrap: wrap;
        justify-content: flex-end;
    }

    button {
        border: none;
        border-radius: 0.55rem;
        padding: 0.42rem 0.78rem;
        font-size: 0.84rem;
        font-weight: 700;
        color: #fff;
        background: var(--accent);
        cursor: pointer;
        transition: background-color 0.2s ease;
    }

    td button:last-child {
        background: var(--danger);
    }

    button:hover {
        background: var(--accent-hover);
    }

    td button:last-child:hover {
        background: var(--danger-hover);
    }

    button:disabled {
        opacity: 0.55;
        cursor: not-allowed;
    }

    .pager {
        margin-top: 0.9rem;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.6rem;
    }

    .pager span {
        color: var(--muted);
        font-weight: 600;
    }

    @media (max-width: 760px) {
        .admin-home h1,
        .club-management {
            width: calc(100% - 1rem);
        }

        .club-management {
            padding: 0.8rem;
            border-radius: 0.85rem;
        }

        table,
        thead,
        tbody,
        th,
        td,
        tr {
            display: block;
            width: 100%;
        }

        thead {
            display: none;
        }

        tr {
            border-bottom: 1px solid #e6edf5;
            padding: 0.45rem 0;
        }

        td {
            border: none;
            padding: 0.35rem 0.2rem;
        }

        td:nth-child(2) {
            justify-content: flex-start;
        }

        .pager {
            flex-direction: column;
            align-items: stretch;
        }

        .pager button {
            width: 100%;
        }
    }

    @media (prefers-color-scheme: dark) {
        .admin-home {
            --bg: linear-gradient(180deg, #0c1725 0%, #0a111b 100%);
            --card: #0f1c2d;
            --text: #deebfb;
            --muted: #b2c6df;
            --border: #2a3b53;
            --accent: #2b8fb5;
            --accent-hover: #3ca3cb;
            --danger: #b84267;
            --danger-hover: #d05079;
        }

        .club-management {
            box-shadow: 0 12px 28px rgba(0, 0, 0, 0.45);
        }

        table {
            background: #102133;
            border-color: #2a3b53;
        }

        th {
            background: #13273e;
            color: #c9d9ef;
        }

        td,
        tr {
            border-bottom-color: #253952;
        }

        tbody tr:hover {
            background: #14263d;
        }
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