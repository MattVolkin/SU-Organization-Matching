<script lang="ts">
// using https://svelte.dev/playground/9116699f10ac42668e7b58d120c4bc8c?version=5.53.2
//using https://github.com/react2svelte/swipeable

  import { swipable } from '@react2svelte/swipeable'
  import type { SwipeEventData }  from '@react2svelte/swipeable'

  function handler(e: CustomEvent<SwipeEventData>) {
    console.log('User swiped!', e.detail);
  }

</script>



<div 
    use:swipeable 
    on:swiped={handler} 
    
  
  
/>