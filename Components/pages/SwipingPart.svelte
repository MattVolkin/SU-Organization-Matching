<script>
    import Header from "../header.svelte";
    import Footer from "../footer.svelte";
    import SwipingApp from "../SwipingApp.svelte";
</script>

<div class="swipe-page">
    <Header />
    <main class="swipe-main">
        <SwipingApp />
    </main>
    <Footer />
</div>

<style>
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