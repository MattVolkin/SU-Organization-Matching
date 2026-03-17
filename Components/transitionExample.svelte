<script>
  import { fade, fly, scale } from 'svelte/transition';

  let show = true;
</script>

<button on:click={() => show = !show}>
  Toggle Elements
</button>

{#if show}
  <div transition:fade={{ duration: 2000 }}>
    < Slow Fade (2s)
  </div>

  <div transition:fly={{ y: 100, duration: 700 }}>
    = Fly in from below
  </div>

  <div transition:scale={{ start: 0.5, duration: 600 }}>
    = Scaling in
  </div>
{/if}