<!-- @component Swiping page wrapper combining header, swiping interface, and footer for the personality traits quiz section. -->
<script>
/**
 * No internal state - this is a layout wrapper component.
 * Delegates all swiping functionality to the child SwipingApp component.
 */
    // Shared page chrome + swipe experience component imports.
    import Header from "../header.svelte";
    import Footer from "../footer.svelte";
    import SwipingApp from "../SwipingApp.svelte";
</script>

<div class="swipe-page">
    <!-- Top navigation -->
    <Header />

    <!-- Main swiping experience -->
    <main class="swipe-main">
        <SwipingApp />
    </main>

    <!-- Global footer -->
    <Footer />
</div>

<style>
    /* Base viewport reset for full-height swipe layout */
    :global(html),
    :global(body) {
        margin: 0;
        padding: 0;
        height: 100%;
    }

    :global(body) {
        overflow: hidden;
    }

    .swipe-page {
        height: 100vh;
        display: flex;
        flex-direction: column;
    }

    .swipe-main {
        flex: 1;
        min-height: 0;
        display: flex;
    }

    .swipe-main :global(.swipedContainer) {
        width: 100%;
        height: 100%;
    }
</style>